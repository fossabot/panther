import React from 'react';
import * as Yup from 'yup';
import pick from 'lodash-es/pick';
import { GET_ORGANIZATION } from 'Pages/general-settings';
import { Alert, Box, Text, useSnackbar } from 'pouncejs';
import { Field, Formik } from 'formik';
import FormikTextInput from 'Components/fields/text-input';
import SubmitButton from 'Components/utils/SubmitButton';
import { useMutation, useQuery } from '@apollo/client';
import { UPDATE_ORGANIZATION } from 'Pages/general-settings/subcomponent/company-information-form';
import { GetOrganizationResponse, UpdateOrganizationInput } from 'Generated/schema';

interface AddRemediationLambdaFormValues {
  awsRemediationLambdaArn: string;
}

interface ApolloQueryData {
  organization: GetOrganizationResponse;
}

interface ApolloMutationInput {
  input: UpdateOrganizationInput;
}

const validationSchema = Yup.object().shape({
  awsRemediationLambdaArn: Yup.string().required(),
});

const AddRemediationLambdaForm: React.FC = () => {
  const { pushSnackbar } = useSnackbar();

  // Get Organization for the existing remediation lambda arn
  const { error: fetchOrganizationError, data: organizationData } = useQuery<ApolloQueryData>(
    GET_ORGANIZATION
  );

  const [
    updateOrganization,
    {
      data: updateOrganizationData,
      error: updateOrganizationError,
      loading: updateOrganizationLoading,
    },
  ] = useMutation<boolean, ApolloMutationInput>(UPDATE_ORGANIZATION);

  React.useEffect(() => {
    if (updateOrganizationData) {
      pushSnackbar({ variant: 'success', title: 'Remediation lambda Arn has been updated' });
    }
  }, [updateOrganizationData]);

  if (fetchOrganizationError) {
    return (
      <Alert
        variant="error"
        title="Failed to query company information"
        description="Sorry, something went wrong and we couldn't fetch the details of your remediation AWS lambda"
      />
    );
  }

  return (
    <Formik<AddRemediationLambdaFormValues>
      enableReinitialize
      initialValues={{
        awsRemediationLambdaArn:
          organizationData?.organization.organization.remediationConfig.awsRemediationLambdaArn,
      }}
      validationSchema={validationSchema}
      onSubmit={async values => {
        await updateOrganization({
          variables: {
            input: {
              ...pick(organizationData.organization.organization, [
                'displayName',
                'email',
                'alertReportFrequency',
              ]),
              remediationConfig: {
                awsRemediationLambdaArn: values.awsRemediationLambdaArn,
              },
            },
          },
          refetchQueries: [{ query: GET_ORGANIZATION }],
        });
      }}
    >
      {({ handleSubmit, isValid, dirty }) => (
        <form onSubmit={handleSubmit}>
          <Box width={0.5} mb={5}>
            <Field
              name="awsRemediationLambdaArn"
              as={FormikTextInput}
              label="Remediation Lambda Arn"
            />
            <Text size="small" color="grey300" mt={2}>
              When the stack is done creating, copy the Role ARN from the &quot;Outputs&quot; tab
              and paste it in the box above
            </Text>
          </Box>
          <SubmitButton
            disabled={updateOrganizationLoading || !isValid || !dirty}
            submitting={updateOrganizationLoading}
          >
            Save ARN
          </SubmitButton>
          {updateOrganizationError && (
            <Alert variant="error" title="Failed to update remediation lambda Arn" mt={6} />
          )}
        </form>
      )}
    </Formik>
  );
};

export default AddRemediationLambdaForm;
