import { Heading, SideSheet, useSnackbar } from 'pouncejs';
import React from 'react';
import { useMutation, gql } from '@apollo/client';
import { LIST_INFRA_SOURCES } from 'Pages/list-sources/subcomponents/infra-source-table';
import { LIST_LOG_SOURCES } from 'Pages/list-sources/subcomponents/log-source-table';
import useSidesheet from 'Hooks/useSidesheet';
import { Integration, UpdateIntegrationInput } from 'Generated/schema';
import UpdateInfraSourceForm from 'Components/forms/update-infra-source-form';
import pick from 'lodash-es/pick';
import { extractErrorMessage } from 'Helpers/utils';
import { INTEGRATION_TYPES } from 'Source/constants';
import UpdateLogSourceForm from 'Components/forms/update-log-source-form';
import { LogSourceFormWrapperValues } from 'Components/forms/log-source-form-wrapper';
import { InfraSourceFormWrapperValues } from 'Components/forms/infra-source-form-wrapper';

const UPDATE_SOURCE = gql`
  mutation UpdateSource($input: UpdateIntegrationInput!) {
    updateIntegration(input: $input)
  }
`;

export interface UpdateAWSSourcesSidesheetProps {
  source: Integration;
}

interface ApolloMutationInput {
  input: UpdateIntegrationInput;
}

export const UpdateAwsSourcesSidesheet: React.FC<UpdateAWSSourcesSidesheetProps> = ({ source }) => {
  const isInfraSource = source.integrationType === INTEGRATION_TYPES.AWS_INFRA;
  const [updateSource, { data, error }] = useMutation<boolean, ApolloMutationInput>(UPDATE_SOURCE);
  const { pushSnackbar } = useSnackbar();
  const { hideSidesheet } = useSidesheet();

  React.useEffect(() => {
    if (error) {
      pushSnackbar({
        variant: 'error',
        title: extractErrorMessage(error) || 'Failed to update your source due to an unknown error',
      });
    }
  }, [error]);

  React.useEffect(() => {
    if (data) {
      pushSnackbar({ variant: 'success', title: `Successfully updated sources` });
      hideSidesheet();
    }
  }, [data]);

  const handleSubmit = (values: InfraSourceFormWrapperValues | LogSourceFormWrapperValues) =>
    updateSource({
      awaitRefetchQueries: true,
      variables: {
        input: {
          ...values,
          integrationId: source.integrationId,
        },
      },
      refetchQueries: [{ query: isInfraSource ? LIST_INFRA_SOURCES : LIST_LOG_SOURCES }],
    });

  return (
    <SideSheet open onClose={hideSidesheet}>
      <Heading size="medium" mb={8}>
        Update Account
      </Heading>
      {isInfraSource ? (
        <UpdateInfraSourceForm
          initialValues={
            pick(source, ['awsAccountId', 'integrationLabel']) as InfraSourceFormWrapperValues
          }
          onSubmit={handleSubmit}
        />
      ) : (
        <UpdateLogSourceForm
          initialValues={
            pick(source, [
              'integrationId',
              'integrationLabel',
              'sourceSnsTopicArn',
              'logProcessingRoleArn',
            ]) as LogSourceFormWrapperValues
          }
          onSubmit={handleSubmit}
        />
      )}
    </SideSheet>
  );
};

export default UpdateAwsSourcesSidesheet;
