import * as Yup from 'yup';
import { Field, Formik } from 'formik';
import React from 'react';
import { Box } from 'pouncejs';
import FormikTextInput from 'Components/fields/text-input';
import SubmitButton from 'Components/utils/SubmitButton';
import useAuth from 'Hooks/useAuth';

interface SignInFormValues {
  username: string;
  password: string;
}

const initialValues = {
  username: '',
  password: '',
};

const validationSchema = Yup.object().shape({
  username: Yup.string()
    .email('Needs to be a valid email')
    .required(),
  password: Yup.string().required(),
});

const SignInForm: React.FC = () => {
  const { signIn } = useAuth();

  return (
    <Formik<SignInFormValues>
      initialValues={initialValues}
      validationSchema={validationSchema}
      onSubmit={async ({ username, password }, { setErrors }) =>
        signIn({
          username,
          password,
          onError: ({ message }) =>
            setErrors({
              password: message,
            }),
        })
      }
    >
      {({ handleSubmit, isSubmitting, isValid, dirty }) => (
        <Box width={1} is="form" onSubmit={handleSubmit}>
          <Field
            as={FormikTextInput}
            label="Email"
            placeholder="Enter your company email..."
            type="email"
            name="username"
            aria-required
            mb={6}
          />
          <Field
            as={FormikTextInput}
            label="Password"
            placeholder="The name of your cat"
            name="password"
            type="password"
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

export default SignInForm;
