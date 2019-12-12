package processor

import (
	"errors"
	"path"
	"strconv"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/sqs"
	jsoniter "github.com/json-iterator/go"
	"go.uber.org/zap"

	alertmodels "github.com/panther-labs/panther/api/alert_processor"
	enginemodels "github.com/panther-labs/panther/api/gateway/analysis"
	analysismodels "github.com/panther-labs/panther/api/gateway/analysis/models"
	complianceops "github.com/panther-labs/panther/api/gateway/compliance/client/operations"
	compliancemodels "github.com/panther-labs/panther/api/gateway/compliance/models"
	resourcemodels "github.com/panther-labs/panther/api/gateway/resources/models"
	"github.com/panther-labs/panther/pkg/awsbatch/sqsbatch"
)

const defaultDelaySeconds = 30

// Map policy/resource ID to the instance of the object
type policyMap map[string]*analysismodels.EnabledPolicy
type resourceMap map[string]*resourcemodels.Resource

// Every batch of sqs messages results in compliance updates and alert/remediation deliveries
type batchResults struct {
	StatusEntries []*compliancemodels.SetStatus
	Alerts        []*sqs.SendMessageBatchRequestEntry
}

// Handle is the entry point for the resource analysis.
func Handle(batch *events.SQSEvent) error {
	resources := make(resourceMap)
	var results batchResults

	for _, record := range batch.Records {
		resource, policy := parseQueueMsg(record.Body)

		if policy != nil {
			// Policy updated - analyze applicable resources now
			if err := results.analyzeUpdatedPolicy(policy); err != nil {
				return err
			}
		} else if resource != nil {
			// Resource updated - analyze with applicable policies (after grouping + deduping)
			resources[string(resource.ID)] = resource
		}
	}

	// Analyze updated resources with applicable policies
	if err := results.analyze(resources, nil); err != nil {
		return err
	}

	return results.deliver()
}

func parseQueueMsg(body string) (*resourcemodels.Resource, *analysismodels.Policy) {
	// There are 2 kinds of possible messages:
	//    a) An updated resource which needs to be evaluated with all applicable policies
	//    b) An updated policy which needs to be evaluated against applicable resources
	var resource resourcemodels.Resource
	err := jsoniter.UnmarshalFromString(body, &resource)
	if err == nil && resource.Attributes != nil && resource.Type != "" {
		zap.L().Info("found new/updated resource",
			zap.String("resourceId", string(resource.ID)))
		return &resource, nil
	}

	// Not a resource - it must be a policy
	var policy analysismodels.Policy
	if err := jsoniter.UnmarshalFromString(body, &policy); err != nil || policy.Body == "" {
		zap.L().Error("failed to parse msg as resource or as policy", zap.String("body", body))
		return nil, nil
	}

	return nil, &policy
}

// Analyze all resources related to a single policy (may require several policy-engine invocations).
func (r *batchResults) analyzeUpdatedPolicy(policy *analysismodels.Policy) error {
	// convert policy to policy map
	policies := policyMap{
		string(policy.ID): &analysismodels.EnabledPolicy{
			Body:          policy.Body,
			ID:            policy.ID,
			ResourceTypes: policy.ResourceTypes,
			Severity:      policy.Severity,
			Suppressions:  policy.Suppressions,
			VersionID:     policy.VersionID,
		},
	}

	// Analyze each page of resources
	// TODO - check for duplicates here
	var totalPages int64 = 1
	for pageno := int64(1); pageno <= totalPages; pageno++ {
		resources, pageCount, err := getResources(policy.ResourceTypes, pageno)
		if err != nil {
			return err
		}

		if err := r.analyze(resources, policies); err != nil {
			return err
		}
		totalPages = pageCount
	}

	return nil
}

// Analyze each org in turn and report status entries and alert notifications across the entire batch
//
// Policies can either be provided by the caller or else they will be fetched from policy-api.
func (r *batchResults) analyze(resources resourceMap, policies policyMap) error {
	// Fetch policies and evaluate them against the resources
	var err error
	if policies == nil {
		if policies, err = getPolicies(); err != nil {
			return err
		}
	}

	if len(policies) == 0 {
		return nil
	}

	var analysis *enginemodels.PolicyEngineOutput
	analysis, err = evaluatePolicies(policies, resources)
	if err != nil {
		return err
	}

	// Add a status entry for every policy/resource pair
	for _, result := range analysis.Resources {
		for _, policyError := range result.Errored {
			entry := buildStatus(policies[policyError.ID], resources[result.ID], compliancemodels.StatusERROR)
			entry.ErrorMessage = compliancemodels.ErrorMessage(policyError.Message)
			r.StatusEntries = append(r.StatusEntries, entry)
		}

		for _, policyID := range result.Failed {
			policy, resource := policies[policyID], resources[result.ID]
			entry := buildStatus(policy, resource, compliancemodels.StatusFAIL)
			r.StatusEntries = append(r.StatusEntries, entry)

			if entry.Suppressed {
				// Suppressed resources are recorded in compliance-api, but do not trigger
				// alerts nor remediations.
				continue
			}

			// Check the current pass/fail status for this policy/resource pair
			var response *complianceops.GetStatusOK
			response, err = complianceClient.Operations.GetStatus(&complianceops.GetStatusParams{
				PolicyID:   policyID,
				ResourceID: result.ID,
				HTTPClient: httpClient,
			})

			if _, ok := err.(*complianceops.GetStatusNotFound); err != nil && !ok {
				// An error other than NotFound
				zap.L().Error("failed to fetch compliance status", zap.Error(err))
				return err
			}

			status := compliancemodels.StatusPASS
			if response != nil {
				status = response.Payload.Status
			}

			zap.L().Info("loaded previous compliance status",
				zap.String("policyId", policyID),
				zap.String("resourceId", result.ID),
				zap.String("complianceStatus", string(status)),
			)

			// Every failed policy, if not suppressed, will trigger the remediation flow
			complianceNotification := &alertmodels.ComplianceNotification{
				ResourceID:      aws.String(string(resource.ID)),
				PolicyID:        aws.String(string(policy.ID)),
				PolicyVersionID: aws.String(string(policy.VersionID)),
				Timestamp:       aws.Time(time.Now()),

				// We only need to send an alert to the user if the status is newly FAILing
				ShouldAlert: aws.Bool(status != compliancemodels.StatusFAIL),
			}
			var sqsMessageBody string
			if sqsMessageBody, err = jsoniter.MarshalToString(complianceNotification); err != nil {
				zap.L().Error("failed to marshal complianceNotification body", zap.Error(err))
				return err
			}

			r.Alerts = append(r.Alerts, &sqs.SendMessageBatchRequestEntry{
				DelaySeconds: aws.Int64(defaultDelaySeconds),
				Id:           aws.String(strconv.Itoa(len(r.Alerts))),
				MessageBody:  aws.String(sqsMessageBody),
			})
		}

		for _, policyID := range result.Passed {
			entry := buildStatus(policies[policyID], resources[result.ID], compliancemodels.StatusPASS)
			r.StatusEntries = append(r.StatusEntries, entry)
		}
	}

	return nil
}

// Invoke the policy engine.
func evaluatePolicies(policies policyMap, resources resourceMap) (*enginemodels.PolicyEngineOutput, error) {
	input := enginemodels.PolicyEngineInput{
		Policies:  make([]enginemodels.Policy, 0, len(policies)),
		Resources: make([]enginemodels.Resource, 0, len(resources)),
	}
	for _, policy := range policies {
		input.Policies = append(input.Policies, enginemodels.Policy{
			Body:          string(policy.Body),
			ID:            string(policy.ID),
			ResourceTypes: policy.ResourceTypes,
		})
	}
	for _, resource := range resources {
		input.Resources = append(input.Resources, enginemodels.Resource{
			Attributes: resource.Attributes,
			ID:         string(resource.ID),
			Type:       string(resource.Type),
		})
	}

	body, err := jsoniter.Marshal(&input)
	if err != nil {
		zap.L().Error("failed to marshal PolicyEngineInput", zap.Error(err))
		return nil, err
	}

	zap.L().Info("invoking policy engine",
		zap.String("policyEngine", env.AnalysisEngine),
		zap.Int("policyCount", len(input.Policies)),
		zap.Int("resourceCount", len(input.Resources)),
	)
	response, err := lambdaClient.Invoke(&lambda.InvokeInput{FunctionName: &env.AnalysisEngine, Payload: body})
	if err != nil {
		zap.L().Error("failed to invoke policy engine", zap.Error(err))
		return nil, err
	}

	if aws.StringValue(response.FunctionError) != "" {
		errorMessage := string(response.Payload)
		zap.L().Error("policy engine returned an error",
			zap.String("payload", errorMessage))
		return nil, errors.New("policy engine error: " + errorMessage)
	}

	var output enginemodels.PolicyEngineOutput
	if err = jsoniter.Unmarshal(response.Payload, &output); err != nil {
		zap.L().Error("failed to unmarshal PolicyEngineOutput", zap.Error(err))
		return nil, err
	}

	zap.L().Debug("successfully invoked policy-engine",
		zap.Any("policyEngineOutput", output))
	return &output, nil
}

// Convert a policy/resource pair into the compliance status struct
func buildStatus(
	policy *analysismodels.EnabledPolicy,
	resource *resourcemodels.Resource,
	status compliancemodels.Status,
) *compliancemodels.SetStatus {

	// Check if the resource is suppressed by this policy
	suppressed := false
	for _, pattern := range policy.Suppressions {
		match, err := path.Match(pattern, string(resource.ID))
		if err != nil {
			zap.L().Warn("invalid policy suppression pattern",
				zap.String("policyId", string(policy.ID)),
				zap.String("pattern", pattern))
			continue
		}

		if match {
			suppressed = true
			break
		}
	}

	return &compliancemodels.SetStatus{
		PolicyID:       compliancemodels.PolicyID(policy.ID),
		PolicySeverity: compliancemodels.PolicySeverity(policy.Severity),
		ResourceID:     compliancemodels.ResourceID(resource.ID),
		ResourceType:   compliancemodels.ResourceType(resource.Type),
		Suppressed:     compliancemodels.Suppressed(suppressed),
		IntegrationID:  compliancemodels.IntegrationID(resource.IntegrationID),

		Status: status,
	}
}

// Deliver all analysis results to compliance-api and alert-processor
func (r *batchResults) deliver() error {
	if len(r.StatusEntries) == 0 {
		return nil // if there aren't any results, there aren't any alerts either
	}

	zap.L().Info("sending status information to compliance-api",
		zap.Int("statusCount", len(r.StatusEntries)))
	if _, err := complianceClient.Operations.SetStatus(&complianceops.SetStatusParams{
		Body:       &compliancemodels.SetStatusBatch{Entries: r.StatusEntries},
		HTTPClient: httpClient,
	}); err != nil {
		zap.L().Error("failed to update status", zap.Error(err))
		return err
	}

	if len(r.Alerts) == 0 {
		return nil
	}

	// Send all alert notifications to the queue in a batch
	zap.L().Info("sending alert notifications",
		zap.String("alertQueue", env.AlertQueueURL),
		zap.Int("notificationCount", len(r.Alerts)))
	batchInput := &sqs.SendMessageBatchInput{Entries: r.Alerts, QueueUrl: &env.AlertQueueURL}
	if err := sqsbatch.SendMessageBatch(sqsClient, maxBackoff, batchInput); err != nil {
		zap.L().Error("failed to send alert notifications", zap.Error(err))
		return err
	}

	return nil
}
