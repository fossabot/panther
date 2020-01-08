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
import LogSourceFormWrapper, {
  LogSourceFormWrapperProps,
} from 'Components/forms/log-source-form-wrapper';

export type UpdateLogSourceFormProps = Omit<LogSourceFormWrapperProps, 'children'>;

const UpdateLogSourceForm: React.FC<UpdateLogSourceFormProps> = ({ onSubmit, initialValues }) => {
  return (
    <LogSourceFormWrapper onSubmit={onSubmit} initialValues={initialValues}>
      {({ handleSubmit, isValid, dirty, isSubmitting }) => (
        <form onSubmit={handleSubmit}>
          <Field
            name="integrationLabel"
            as={FormikTextInput}
            label="Label"
            placeholder="A nickname for this log source"
            aria-required
            mb={6}
          />
          <Field
            disabled
            readonly
            name="sourceSnsTopicArn"
            as={FormikTextInput}
            label="SNS Topic ARN"
            placeholder="The SNS Topic receiving log delivery notifications"
            aria-required
            mb={6}
          />
          <Field
            disabled
            readonly
            name="logProcessingRoleArn"
            as={FormikTextInput}
            label="Assumable Role ARN"
            placeholder="The ARN of the role to read logs from S3"
            aria-required
            mb={6}
          />
          <SubmitButton
            width={1}
            disabled={isSubmitting || !isValid || !dirty}
            submitting={isSubmitting}
          >
            Update
          </SubmitButton>
        </form>
      )}
    </LogSourceFormWrapper>
  );
};

export default UpdateLogSourceForm;
