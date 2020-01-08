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
import { Text } from 'pouncejs';
import { DestinationConfigInput } from 'Generated/schema';
import BaseDestinationForm, {
  BaseDestinationFormValues,
  defaultValidationSchema,
} from 'Components/forms/common/base-destination-form';
import JsonViewer from 'Components/json-viewer';

type SQSFieldValues = Pick<DestinationConfigInput, 'sqs'>;

interface SQSDestinationFormProps {
  initialValues: BaseDestinationFormValues<SQSFieldValues>;
  onSubmit: (values: BaseDestinationFormValues<SQSFieldValues>) => void;
}

const sqsFieldsValidationSchema = Yup.object().shape({
  outputConfig: Yup.object().shape({
    sqs: Yup.object().shape({
      queueUrl: Yup.string()
        .url('Queue URL must be a valid url')
        .required('Queue URL is required'),
    }),
  }),
});

const SQS_QUEUE_POLICY = {
  Version: '2012-10-17',
  Statement: [
    {
      Sid: 'AllowPantherAlarming',
      Effect: 'Allow',
      Action: 'sqs:SendMessage',
      Principal: {
        AWS: process.env.AWS_ACCOUNT_ID,
      },
      Resource: '<The ARN of the SQS Queue they are adding as output>',
    },
  ],
};

// @ts-ignore
// We merge the two schemas together: the one deriving from the common fields, plus the custom
// ones that change for each destination.
// https://github.com/jquense/yup/issues/522
const mergedValidationSchema = defaultValidationSchema.concat(sqsFieldsValidationSchema);

const SQSDestinationForm: React.FC<SQSDestinationFormProps> = ({ onSubmit, initialValues }) => {
  return (
    <BaseDestinationForm<SQSFieldValues>
      initialValues={initialValues}
      validationSchema={mergedValidationSchema}
      onSubmit={onSubmit}
    >
      <Text size="small">Note: Add note here</Text>
      <Field
        as={FormikTextInput}
        name="outputConfig.sqs.queueUrl"
        label="Queue URL"
        placeholder="Where should we send the queue data to?"
        mb={6}
        aria-required
      />
      <Text size="medium" mb={2}>
        <b>Note</b>: You would need to allow Panther <b>sqs:SendMessage</b> access to send alert
        messages to your SQS queue
      </Text>
      <JsonViewer data={SQS_QUEUE_POLICY} collapsed={false} />
    </BaseDestinationForm>
  );
};

export default SQSDestinationForm;
