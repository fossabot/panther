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
import * as Yup from 'yup';
import FormikTextInput from 'Components/fields/text-input';
import { DestinationConfigInput } from 'Generated/schema';
import { Text } from 'pouncejs';
import BaseDestinationForm, {
  BaseDestinationFormValues,
  defaultValidationSchema,
} from 'Components/forms/common/base-destination-form';

type EmailFieldValues = Pick<DestinationConfigInput, 'email'>;

interface EmailDestinationFormProps {
  initialValues: BaseDestinationFormValues<EmailFieldValues>;
  onSubmit: (values: BaseDestinationFormValues<EmailFieldValues>) => void;
}

const emailFieldsValidationSchema = Yup.object().shape({
  outputConfig: Yup.object().shape({
    email: Yup.object().shape({
      destinationAddress: Yup.string()
        .email('Must be a valid email address')
        .required(),
    }),
  }),
});

// @ts-ignore
// We merge the two schemas together: the one deriving from the common fields, plus the custom
// ones that change for each destination.
// https://github.com/jquense/yup/issues/522
const mergedValidationSchema = defaultValidationSchema.concat(emailFieldsValidationSchema);

const EmailDestinationForm: React.FC<EmailDestinationFormProps> = ({ onSubmit, initialValues }) => {
  return (
    <BaseDestinationForm<EmailFieldValues>
      initialValues={initialValues}
      validationSchema={mergedValidationSchema}
      onSubmit={onSubmit}
    >
      <Field
        as={FormikTextInput}
        name="outputConfig.email.destinationAddress"
        label="Email Address"
        placeholder="Where should we send an email notification to?"
        mb={3}
        aria-required
      />
      <Text size="small" color="grey300" mb={6}>
        * If the email address is not already verified, we will immediately send a verification
        email to it. Until it gets verified, it will not be eligible to receive any alerts.
      </Text>
    </BaseDestinationForm>
  );
};

export default EmailDestinationForm;
