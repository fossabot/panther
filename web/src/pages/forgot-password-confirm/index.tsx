import React from 'react';
import Banner from 'Assets/sign-up-banner.jpg';
import AuthPageContainer from 'Components/auth-page-container';
import queryString from 'query-string';
import ForgotPasswordConfirmForm from 'Components/forms/forgot-password-confirm-form';
import useRouter from 'Hooks/useRouter';

const ForgotPasswordConfirmPage: React.FC = () => {
  const { location } = useRouter();

  // protect against not having the proper parameters in place
  const { email, token } = queryString.parse(location.search) as { email: string; token: string };
  if (!token || !email) {
    return (
      <AuthPageContainer banner={Banner}>
        <AuthPageContainer.Caption
          title="Something seems off..."
          subtitle="Are you sure that the URL you followed is valid?"
        />
      </AuthPageContainer>
    );
  }

  return (
    <AuthPageContainer banner={Banner}>
      <AuthPageContainer.Caption
        title="Alrighty then.."
        subtitle="Let's set you up with a new password."
      />
      <ForgotPasswordConfirmForm email={email} token={token} />
    </AuthPageContainer>
  );
};

export default ForgotPasswordConfirmPage;
