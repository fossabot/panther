/* eslint-disable react/display-name */
import React from 'react';
import { Card, Flex, Text, Heading, Alert, Box } from 'pouncejs';
import { INTEGRATION_TYPES } from 'Source/constants';
import Wizard from 'Components/wizard';
import InfraSourceFormWrapper, {
  InfraSourceFormWrapperValues,
} from 'Components/forms/infra-source-form-wrapper';
import urls from 'Source/urls';
import { extractErrorMessage } from 'Helpers/utils';
import { useMutation, gql } from '@apollo/client';
import { LIST_INFRA_SOURCES } from 'Pages/list-sources/subcomponents/infra-source-table';
import SubmitButton from 'Components/utils/SubmitButton';
import ErrorBoundary from 'Components/error-boundary';
import { Field } from 'formik';
import FormikTextInput from 'Components/fields/text-input';
import useRouter from 'Hooks/useRouter';
import RemediationPanel from './remediation-panel';
import RealTimeEventPanel from './real-time-event-panel';
import ResourceScanningPanel from './resource-scanning-panel';
import PanelWrapper from './panel-wrapper';

const ADD_INFRA_SOURCE = gql`
  mutation AddInfraSource($input: AddIntegrationInput!) {
    addIntegration(input: $input)
  }
`;

const initialInfraSourceFormValues = {
  awsAccountId: '',
  integrationLabel: '',
};

const CreateInfraSource: React.FC = () => {
  const { history } = useRouter();
  const [addInfraSource, { data, loading, error }] = useMutation(ADD_INFRA_SOURCE);

  const submitSourceToServer = React.useCallback(
    (values: InfraSourceFormWrapperValues) =>
      addInfraSource({
        awaitRefetchQueries: true,
        variables: {
          input: {
            integrations: [
              {
                awsAccountID: values.awsAccountId,
                integrationLabel: values.integrationLabel,
                integrationType: INTEGRATION_TYPES.AWS_INFRA,
              },
            ],
          },
        },
        refetchQueries: [{ query: LIST_INFRA_SOURCES }],
      }),
    []
  );

  React.useEffect(() => {
    if (data) {
      history.push(urls.account.settings.sources.list());
    }
  });

  return (
    <Box>
      {error && (
        <Alert
          variant="error"
          title="An error has occurred"
          description={
            extractErrorMessage(error) || "We couldn't store your account due to an internal error"
          }
          mb={6}
        />
      )}
      <Card p={9}>
        <InfraSourceFormWrapper
          initialValues={initialInfraSourceFormValues}
          onSubmit={submitSourceToServer}
        >
          {({ isValid, dirty }) => (
            <Flex justifyContent="center" alignItems="center" width={1}>
              <Wizard<InfraSourceFormWrapperValues>
                autoCompleteLastStep
                steps={[
                  {
                    title: 'Account Details',
                    icon: 'add' as const,
                    renderStep: ({ goToNextStep }) => {
                      const shouldEnableNextButton = dirty && isValid;

                      return (
                        <PanelWrapper>
                          <PanelWrapper.Content>
                            <Box width={460} m="auto">
                              <Heading size="medium" m="auto" mb={5} color="grey400">
                                Let{"'"}s start with some account information
                              </Heading>
                              <Text size="large" color="grey200" mb={10}>
                                Before we begin, we need to setup a few things in our database
                              </Text>
                              <ErrorBoundary>
                                <Field
                                  name="awsAccountId"
                                  as={FormikTextInput}
                                  label="AWS Account ID"
                                  placeholder="Your 12-digit AWS Account ID"
                                  aria-required
                                  mb={6}
                                />
                                <Field
                                  name="integrationLabel"
                                  as={FormikTextInput}
                                  label="Label"
                                  placeholder="A nickname for your account"
                                  aria-required
                                  mb={6}
                                />
                              </ErrorBoundary>
                            </Box>
                          </PanelWrapper.Content>
                          <PanelWrapper.WizardActions
                            goToNextStep={goToNextStep}
                            isNextStepDisabled={!shouldEnableNextButton}
                          />
                        </PanelWrapper>
                      );
                    },
                  },
                  {
                    title: 'Scanning',
                    icon: 'search',
                    renderStep: ({ goToPrevStep, goToNextStep }) => (
                      <PanelWrapper>
                        <PanelWrapper.Content>
                          <ResourceScanningPanel />
                        </PanelWrapper.Content>
                        <PanelWrapper.WizardActions
                          goToNextStep={goToNextStep}
                          goToPrevStep={goToPrevStep}
                        />
                      </PanelWrapper>
                    ),
                  },
                  {
                    title: 'Real Time',
                    icon: 'sync',
                    renderStep: ({ goToPrevStep, goToNextStep }) => (
                      <PanelWrapper>
                        <PanelWrapper.Content>
                          <RealTimeEventPanel />
                        </PanelWrapper.Content>
                        <PanelWrapper.WizardActions
                          goToNextStep={goToNextStep}
                          goToPrevStep={goToPrevStep}
                        />
                      </PanelWrapper>
                    ),
                  },
                  {
                    title: 'Remediation',
                    icon: 'wrench',
                    renderStep: ({ goToPrevStep, goToNextStep }) => (
                      <PanelWrapper>
                        <PanelWrapper.Content>
                          <RemediationPanel />
                        </PanelWrapper.Content>
                        <PanelWrapper.WizardActions
                          goToNextStep={goToNextStep}
                          goToPrevStep={goToPrevStep}
                        />
                      </PanelWrapper>
                    ),
                  },
                  {
                    title: 'Done!',
                    icon: 'check',
                    renderStep: () => (
                      <Flex
                        justifyContent="center"
                        alignItems="center"
                        flexDirection="column"
                        my={190}
                        mx="auto"
                        width={300}
                      >
                        <Heading size="medium" m="auto" mb={5} color="grey400">
                          Almost done!
                        </Heading>
                        <Text size="large" color="grey300" mb={10}>
                          Click the button below to complete setup!
                        </Text>
                        <SubmitButton width={350} disabled={loading} submitting={loading}>
                          Add New Source
                        </SubmitButton>
                      </Flex>
                    ),
                  },
                ]}
              />
            </Flex>
          )}
        </InfraSourceFormWrapper>
      </Card>
    </Box>
  );
};

export default CreateInfraSource;
