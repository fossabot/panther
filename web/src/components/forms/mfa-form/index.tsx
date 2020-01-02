import { Field, Formik } from 'formik';
import React from 'react';
import * as Yup from 'yup';
import { Box } from 'pouncejs';
import SubmitButton from 'Components/utils/SubmitButton';
import FormikTextInput from 'Components/fields/text-input';
import useAuth from 'Hooks/useAuth';

interface MfaFormValues {
  mfaCode: string;
}

const initialValues = {
  mfaCode: '',
};

const validationSchema = Yup.object().shape({
  mfaCode: Yup.string()
    .matches(/\b\d{6}\b/, 'Code should contain exactly six digits.')
    .required(),
});

const MfaForm: React.FC = () => {
  const { confirmSignIn } = useAuth();

  return (
    <Formik<MfaFormValues>
      initialValues={initialValues}
      validationSchema={validationSchema}
      onSubmit={async ({ mfaCode }, { setErrors }) =>
        confirmSignIn({
          mfaCode,
          onError: ({ message }) =>
            setErrors({
              mfaCode: message,
            }),
        })
      }
    >
      {({ handleSubmit, isValid, isSubmitting, dirty }) => (
        <Box is="form" width={1} onSubmit={handleSubmit}>
          <Field
            autoFocus
            as={FormikTextInput}
            placeholder="The 6-digit MFA code"
            name="mfaCode"
            autoComplete="off"
            aria-required
            mb={6}
          />
          <SubmitButton
            width={1}
            submitting={isSubmitting}
            disabled={isSubmitting || !isValid || !dirty}
          >
            Sign in
          </SubmitButton>
        </Box>
      )}
    </Formik>
  );
};

export default MfaForm;
