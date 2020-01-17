/**
 * Panther is a scalable, powerful, cloud-native SIEM written in Golang/React.
 * Copyright (C) 2020 Panther Labs Inc
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
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
