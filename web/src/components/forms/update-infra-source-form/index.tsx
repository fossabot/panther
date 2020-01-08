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
