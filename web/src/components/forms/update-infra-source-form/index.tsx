import React from 'react';
import { Field } from 'formik';
import FormikTextInput from 'Components/fields/text-input';
import SubmitButton from 'Components/utils/SubmitButton';
import InfraSourceFormWrapper, {
  InfraSourceFormWrapperProps,
} from 'Components/forms/infra-source-form-wrapper';

export type UpdateInfraSourceFormProps = Omit<InfraSourceFormWrapperProps, 'children'>;

const UpdateInfraSourceForm: React.FC<UpdateInfraSourceFormProps> = ({
  onSubmit,
  initialValues,
}) => {
  return (
    <InfraSourceFormWrapper initialValues={initialValues} onSubmit={onSubmit}>
      {({ isSubmitting, isValid, dirty }) => (
        <React.Fragment>
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
          <SubmitButton
            width={1}
            disabled={isSubmitting || !isValid || !dirty}
            submitting={isSubmitting}
          >
            {initialValues.awsAccountId ? 'Update' : 'Add'}
          </SubmitButton>
        </React.Fragment>
      )}
    </InfraSourceFormWrapper>
  );
};

export default UpdateInfraSourceForm;
