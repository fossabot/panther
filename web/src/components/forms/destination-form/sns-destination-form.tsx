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

const SNS_TOPIC_POLICY = {
  Version: '2012-10-17',
  Statement: [
    {
      Sid: 'AllowPantherAlarming',
      Effect: 'Allow',
      Action: 'sns:Publish',
      Principal: {
        AWS: process.env.AWS_ACCOUNT_ID,
      },
      Resource: '<The ARN of the SNS Topic you are adding as the destination>',
    },
  ],
};

type SNSFieldValues = Pick<DestinationConfigInput, 'sns'>;

interface SNSDestinationFormProps {
  initialValues: BaseDestinationFormValues<SNSFieldValues>;
  onSubmit: (values: BaseDestinationFormValues<SNSFieldValues>) => void;
}

const snsFieldsValidationSchema = Yup.object().shape({
  outputConfig: Yup.object().shape({
    sns: Yup.object().shape({
      topicArn: Yup.string().required(),
    }),
  }),
});

// @ts-ignore
// We merge the two schemas together: the one deriving from the common fields, plus the custom
// ones that change for each destination.
// https://github.com/jquense/yup/issues/522
const mergedValidationSchema = defaultValidationSchema.concat(snsFieldsValidationSchema);

const SNSDestinationForm: React.FC<SNSDestinationFormProps> = ({ onSubmit, initialValues }) => {
  return (
    <BaseDestinationForm<SNSFieldValues>
      initialValues={initialValues}
      validationSchema={mergedValidationSchema}
      onSubmit={onSubmit}
    >
      <Field
        as={FormikTextInput}
        name="outputConfig.sns.topicArn"
        label="Topic ARN"
        placeholder="Where should we publish a notification to?"
        mb={6}
        aria-required
      />
      <Text size="medium" mb={2}>
        <b>Note</b>: You would need to allow Panther <b>sns:Publish</b> access to send alert
        messages to your SNS topic
      </Text>
      <JsonViewer data={SNS_TOPIC_POLICY} collapsed={false} />
    </BaseDestinationForm>
  );
};

export default SNSDestinationForm;
