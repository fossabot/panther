/**
 * Panther is a scalable, powerful, cloud-native SIEM written in Golang/React.
 * Copyright (C) 2020 Panther Labs Inc
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

/* eslint-disable react/display-name */
import React from 'react';
import { Card, Flex, Text, Heading, Alert, Box } from 'pouncejs';
import { INTEGRATION_TYPES } from 'Source/constants';
import Wizard from 'Components/wizard';
import urls from 'Source/urls';
import { extractErrorMessage } from 'Helpers/utils';
import { useMutation, gql } from '@apollo/client';
import { LIST_LOG_SOURCES } from 'Pages/list-sources/subcomponents/log-source-table';
import SubmitButton from 'Components/utils/SubmitButton';
import ErrorBoundary from 'Components/error-boundary';
import useRouter from 'Hooks/useRouter';
import LogSourceFormWrapper, {
  LogSourceFormWrapperValues,
} from 'Components/forms/log-source-form-wrapper';
import { Field } from 'formik';
import FormikTextInput from 'Components/fields/text-input';
import SnsLogConnectionPanel from './sns-log-connection-panel';
import PanelWrapper from './panel-wrapper';

const ADD_LOG_SOURCE = gql`
  mutation AddSource($input: AddIntegrationInput!) {
    addIntegration(input: $input)
  }
`;

const initialLogSourceFormValues = {
  integrationLabel: '',
  sourceSnsTopicArn: '',
  logProcessingRoleArn: '',
};

const CreateLogSource: React.FC = () => {
  const { history } = useRouter();
  const [addLogSource, { data, loading, error }] = useMutation(ADD_LOG_SOURCE);

  const submitSourceToServer = React.useCallback(
    (values: LogSourceFormWrapperValues) =>
      addLogSource({
        awaitRefetchQueries: true,
        variables: {
          input: {
            integrations: [
              {
                ...values,
                integrationType: INTEGRATION_TYPES.AWS_LOGS,
              },
            ],
          },
        },
        refetchQueries: [{ query: LIST_LOG_SOURCES }],
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
            extractErrorMessage(error) || "We couldn't store your source due to an internal error"
          }
          mb={6}
        />
      )}
      <Card p={9}>
        <LogSourceFormWrapper
          initialValues={initialLogSourceFormValues}
          onSubmit={submitSourceToServer}
        >
          {({ isValid, errors, dirty }) => (
            <Flex justifyContent="center" alignItems="center" width={1}>
              <Wizard<LogSourceFormWrapperValues>
                autoCompleteLastStep
                steps={[
                  {
                    title: 'Source Details',
                    icon: 'add' as const,
                    renderStep: ({ goToNextStep }) => {
                      const shouldEnableNextButton =
                        dirty && !errors.integrationLabel && !errors.sourceSnsTopicArn;

                      return (
                        <PanelWrapper>
                          <PanelWrapper.Content>
                            <Box width={460} m="auto">
                              <Heading size="medium" m="auto" mb={5} color="grey400">
                                Let{"'"}s start with the basics
                              </Heading>
                              <Text size="large" color="grey200" mb={10}>
                                We need to know where to get your logs from
                              </Text>
                              <ErrorBoundary>
                                <Field
                                  name="integrationLabel"
                                  as={FormikTextInput}
                                  label="Label"
                                  placeholder="A nickname for this log source"
                                  aria-required
                                  mb={6}
                                />
                                <Field
                                  name="sourceSnsTopicArn"
                                  as={FormikTextInput}
                                  label="SNS Topic ARN"
                                  placeholder="The SNS Topic receiving log delivery notifications"
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
                    title: 'AWS Permissions Grant',
                    icon: 'search',
                    renderStep: ({ goToPrevStep, goToNextStep }) => {
                      const shouldEnableNextButton = dirty && isValid;

                      return (
                        <PanelWrapper>
                          <PanelWrapper.Content>
                            <SnsLogConnectionPanel />
                          </PanelWrapper.Content>
                          <PanelWrapper.WizardActions
                            goToPrevStep={goToPrevStep}
                            goToNextStep={goToNextStep}
                            isNextStepDisabled={!shouldEnableNextButton}
                          />
                        </PanelWrapper>
                      );
                    },
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
                          Add New Log Source
                        </SubmitButton>
                      </Flex>
                    ),
                  },
                ]}
              />
            </Flex>
          )}
        </LogSourceFormWrapper>
      </Card>
    </Box>
  );
};

export default CreateLogSource;
