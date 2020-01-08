/**
 * Copyright 2020 Panther Labs Inc
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import React from 'react';
import { Flex, Text } from 'pouncejs';
import urls from 'Source/urls';
import { Link } from 'react-router-dom';
import AuthPageContainer from 'Components/auth-page-container';
import Banner from 'Assets/sign-in-banner.jpg';
import SetPasswordForm from 'Components/forms/set-password-form';
import MfaForm from 'Components/forms/mfa-form';
import TotpForm from 'Components/forms/totp-form';
import SignInForm from 'Components/forms/signin-form';
import useAuth from 'Hooks/useAuth';
import { CHALLENGE_NAMES } from 'Components/utils/auth-context';

/**
 * This page is only visitable by non-auth users (since it's sitting behind a guarded route). Thus,
 * no authenticated user will ever be in this page so we don't have to handle the redirect here
 * (meaning the redirect to the referrer page when the user became authenticated). This happens in
 * the `guarded-route` component which handles auth-related accesses and redirects.
 */
const SignInPage: React.FC = () => {
  const { currentAuthChallengeName, userInfo } = useAuth();

  // if there is an user object, then we check to see if he should go through any challenges. If he
  // should then we render the appropriate screen
  switch (currentAuthChallengeName) {
    case CHALLENGE_NAMES.SOFTWARE_TOKEN_MFA:
      return (
        <AuthPageContainer banner={Banner}>
          <AuthPageContainer.Caption
            title="One last thing..."
            subtitle="Enter your MFA code to complete the sign-in"
          />
          <MfaForm />
          <AuthPageContainer.AltOptions>
            <Text size="medium" color="grey200">
              Can{"'"}t seem to get it right?{' '}
              <a
                href={`mailto:support@runpanther.io?subject=MFA issues for ${userInfo?.email}`}
                rel="noopener noreferrer"
              >
                Email us
              </a>
            </Text>
          </AuthPageContainer.AltOptions>
        </AuthPageContainer>
      );
    case CHALLENGE_NAMES.MFA_SETUP:
      return (
        <AuthPageContainer banner={Banner}>
          <AuthPageContainer.Caption
            title="Great!"
            subtitle="Now let's set up two-factor authentication for your account."
          />
          <TotpForm />
        </AuthPageContainer>
      );
    case CHALLENGE_NAMES.NEW_PASSWORD_REQUIRED:
      return (
        <AuthPageContainer banner={Banner}>
          <AuthPageContainer.Caption
            title="First things first"
            subtitle="We need to set you up with a new password."
          />
          <SetPasswordForm />
        </AuthPageContainer>
      );
    default:
      return (
        <AuthPageContainer banner={Banner}>
          <AuthPageContainer.Caption title="Sign in" subtitle="to continue to Panther" />
          <SignInForm />
          <Text mt={10} size="large" color="grey200" textAlign="center">
            <Link to={urls.account.auth.forgotPassword()}>Forgot your password?</Link>
          </Text>
          <AuthPageContainer.AltOptions>
            <Flex alignItems="center">
              <Text size="medium" color="grey200" is="span" mr={3}>
                Don{"'"}t have an account? Talk to your admin
              </Text>
            </Flex>
          </AuthPageContainer.AltOptions>
        </AuthPageContainer>
      );
  }
};

export default SignInPage;
