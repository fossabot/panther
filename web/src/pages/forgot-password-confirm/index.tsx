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
