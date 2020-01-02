import React from 'react';
import { Field, Formik } from 'formik';
import * as Yup from 'yup';
import SubmitButton from 'Components/utils/SubmitButton';
import FormikTextInput from 'Components/fields/text-input';
import useAuth from 'Hooks/useAuth';
import { Card, Text } from 'pouncejs';

interface ForgotPasswordFormValues {
  email: string;
}

const initialValues = {
  email: '',
};

const validationSchema = Yup.object().shape({
  email: Yup.string()
    .email('Needs to be a valid email')
    .required(),
});

const ForgotPasswordForm: React.FC = () => {
  const { forgotPassword } = useAuth();

  return (
    <Formik<ForgotPasswordFormValues>
      initialValues={initialValues}
      validationSchema={validationSchema}
      onSubmit={async ({ email }, { setErrors, setStatus }) =>
        forgotPassword({
          email,
          onSuccess: () => setStatus('SENT'),
          onError: ({ code, message }) => {
            setErrors({
              email:
                code === 'UserNotFoundException'
                  ? "We couldn't find this Panther account"
                  : message,
            });
          },
        })
      }
    >
      {({ handleSubmit, isSubmitting, isValid, dirty, status, values }) => {
        if (status === 'SENT') {
          return (
            <Card bg="#def7e9" p={5} mb={8} boxShadow="none">
              <Text color="green300" size="large">
                We have successfully sent you an email with reset instructions at{' '}
                <b>{values.email}</b>
              </Text>
            </Card>
          );
        }

        return (
          <form onSubmit={handleSubmit}>
            <Field
              as={FormikTextInput}
              label="Email"
              placeholder="Enter your company email..."
              type="email"
              name="email"
              aria-required
              mb={6}
            />
            <SubmitButton
              width={1}
              submitting={isSubmitting}
              disabled={isSubmitting || !isValid || !dirty}
            >
              Reset Password
            </SubmitButton>
          </form>
        );
      }}
    </Formik>
  );
};

export default ForgotPasswordForm;
