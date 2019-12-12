package aws

import "github.com/aws/aws-sdk-go/service/waf"

const (
	WafWebAclSchema         = "AWS.WAF.WebACL"
	WafRegionalWebAclSchema = "AWS.WAF.Regional.WebACL"
)

// WafWebAcl contains all information about a web acl
type WafWebAcl struct {
	// Generic resource fields
	GenericAWSResource
	GenericResource

	// Fields embedded from waf.WebAcl
	DefaultAction *waf.WafAction
	MetricName    *string

	// Additional fields
	Rules []*WafRule
}

type WafRule struct {
	*waf.ActivatedRule
	*waf.Rule
	RuleId *string
}
