import Banner from 'Assets/sign-up-banner.jpg';
import AuthPageContainer from 'Components/auth-page-container';
import ForgotPasswordForm from 'Components/forms/forgot-password-form';
import { Button, Flex, Text } from 'pouncejs';
import urls from 'Source/urls';
import { Link } from 'react-router-dom';
import React from 'react';

interface EmailStatusState {
  state: 'SENT' | 'FAILED' | 'PENDING';
  message?: string;
}

const ForgotPasswordPage: React.FC = () => {
  return (
    <AuthPageContainer banner={Banner}>
      <AuthPageContainer.Caption
        title="Forgot your password?"
        subtitle="We'll help you reset your password and get back on track."
      />
      <ForgotPasswordForm />
      <Text size="small" color="grey200" mt={8} is="p" textAlign="center">
        <i>
          By clicking the button above you will receive an email with instructions on how to reset
          your password
        </i>
      </Text>
      <AuthPageContainer.AltOptions>
        <Flex alignItems="center">
          <Text size="medium" color="grey200" is="span" mr={3}>
            Remembered it all of a sudden?
          </Text>
          <Button
            size="small"
            variant="default"
            is={Link}
            to={urls.account.auth.signIn()}
            style={{ textDecoration: 'none' }}
          >
            Sign in
          </Button>
        </Flex>
      </AuthPageContainer.AltOptions>
    </AuthPageContainer>
  );
};

export default ForgotPasswordPage;
