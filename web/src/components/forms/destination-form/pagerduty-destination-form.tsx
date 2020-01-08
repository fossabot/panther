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
import BaseDestinationForm, {
  BaseDestinationFormValues,
  defaultValidationSchema,
} from 'Components/forms/common/base-destination-form';

type PagerDutyFieldValues = Pick<DestinationConfigInput, 'pagerDuty'>;

interface PagerDutyDestinationFormProps {
  initialValues: BaseDestinationFormValues<PagerDutyFieldValues>;
  onSubmit: (values: BaseDestinationFormValues<PagerDutyFieldValues>) => void;
}

const pagerDutyFieldsValidationSchema = Yup.object().shape({
  outputConfig: Yup.object().shape({
    pagerDuty: Yup.object().shape({
      integrationKey: Yup.string()
        .length(32, 'Must be exactly 32 characters')
        .required(),
    }),
  }),
});

// @ts-ignore
// We merge the two schemas together: the one deriving from the common fields, plus the custom
// ones that change for each destination.
// https://github.com/jquense/yup/issues/522
const mergedValidationSchema = defaultValidationSchema.concat(pagerDutyFieldsValidationSchema);

const PagerDutyDestinationForm: React.FC<PagerDutyDestinationFormProps> = ({
  onSubmit,
  initialValues,
}) => {
  return (
    <BaseDestinationForm<PagerDutyFieldValues>
      initialValues={initialValues}
      validationSchema={mergedValidationSchema}
      onSubmit={onSubmit}
    >
      <Field
        as={FormikTextInput}
        name="outputConfig.pagerDuty.integrationKey"
        label="Integration Key"
        placeholder="What's your PagerDuty Integration Key?"
        mb={6}
        aria-required
        autoComplete="new-password"
      />
    </BaseDestinationForm>
  );
};

export default PagerDutyDestinationForm;
