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

type MicrosoftTeamsFieldValues = Pick<DestinationConfigInput, 'msTeams'>;

interface MicrosoftTeamsDestinationFormProps {
  initialValues: BaseDestinationFormValues<MicrosoftTeamsFieldValues>;
  onSubmit: (values: BaseDestinationFormValues<MicrosoftTeamsFieldValues>) => void;
}

const msTeamsFieldsValidationSchema = Yup.object().shape({
  outputConfig: Yup.object().shape({
    msTeams: Yup.object().shape({
      webhookURL: Yup.string()
        .url('Must be a valid webhook URL')
        .required(),
    }),
  }),
});

// @ts-ignore
// We merge the two schemas together: the one deriving from the common fields, plus the custom
// ones that change for each destination.
// https://github.com/jquense/yup/issues/522
const mergedValidationSchema = defaultValidationSchema.concat(msTeamsFieldsValidationSchema);

const MicrosoftTeamsDestinationForm: React.FC<MicrosoftTeamsDestinationFormProps> = ({
  onSubmit,
  initialValues,
}) => {
  return (
    <BaseDestinationForm<MicrosoftTeamsFieldValues>
      initialValues={initialValues}
      validationSchema={mergedValidationSchema}
      onSubmit={onSubmit}
    >
      <Field
        as={FormikTextInput}
        name="outputConfig.msTeams.webhookURL"
        label="Microsoft Teams Webhook URL"
        placeholder="Where should we send a push notification to?"
        mb={6}
        aria-required
      />
    </BaseDestinationForm>
  );
};

export default MicrosoftTeamsDestinationForm;
