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
import { Field, useFormikContext } from 'formik';
import FormikTextInput from 'Components/fields/text-input';
import { InputElementLabel, Grid, Flex, Box, InputElementErrorLabel, Text } from 'pouncejs';
import { SeverityEnum } from 'Generated/schema';
import { capitalize } from 'Helpers/utils';
import FormikTextArea from 'Components/fields/textarea';
import FormikSwitch from 'Components/fields/switch';
import FormikCombobox from 'Components/fields/combobox';
import FormikMultiCombobox from 'Components/fields/multi-combobox';
import FormikEditor from 'Components/fields/editor';
import { LOG_TYPES, RESOURCE_TYPES } from 'Source/constants';
import { RuleFormValues } from 'Components/forms/rule-form';
import { PolicyFormValues } from 'Components/forms/policy-form';

export const ruleCoreEditableFields = [
  'body',
  'description',
  'displayName',
  'enabled',
  'id',
  'reference',
  'runbook',
  'severity',
  'tags',
] as const;

interface BaseRuleCoreFieldsProps {
  type: 'rule' | 'policy';
}

type FormValues = Required<Pick<RuleFormValues, typeof ruleCoreEditableFields[number]>> &
  Pick<RuleFormValues, 'logTypes'> &
  Pick<PolicyFormValues, 'resourceTypes' | 'suppressions'>;

const severityOptions = Object.values(SeverityEnum);

const BaseRuleCoreFields: React.FC<BaseRuleCoreFieldsProps> = ({ type }) => {
  // Read the values from the "parent" form. We expect a formik to be declared in the upper scope
  // since this is a "partial" form. If no Formik context is found this will error out intentionally
  const { values, errors, touched, initialValues } = useFormikContext<FormValues>();

  return (
    <section>
      <Grid gridTemplateColumns="1fr 1fr" gridRowGap={2} gridColumnGap={9}>
        <Box>
          <Flex justifyContent="space-between">
            <Flex alignItems="center">
              <InputElementLabel htmlFor="enabled" mr={6}>
                Enabled
              </InputElementLabel>
              <Field as={FormikSwitch} name="enabled" />
            </Flex>
            <Flex alignItems="center">
              <InputElementLabel htmlFor="severity" mr={6}>
                * Severity
              </InputElementLabel>
              <Field
                as={FormikCombobox}
                name="severity"
                items={severityOptions}
                itemToString={severity => capitalize(severity.toLowerCase())}
              />
            </Flex>
          </Flex>
        </Box>
        <div />
        <Field
          as={FormikTextInput}
          label="* ID"
          placeholder={`The unique ID of this ${type}`}
          name="id"
          disabled={initialValues.id}
          aria-required
        />
        <Field
          as={FormikTextInput}
          label="Display Name"
          placeholder={`A human-friendly name for this ${type}`}
          name="displayName"
        />
        <Field
          as={FormikTextInput}
          label="Runbook"
          placeholder={`Procedures and operations related to this ${type}`}
          name="runbook"
        />
        <Field
          as={FormikTextInput}
          label="Reference"
          placeholder={`An external link to why this ${type} exists`}
          name="reference"
        />
        <Field
          as={FormikTextArea}
          label="Description"
          placeholder={`Additional context about this ${type}`}
          name="description"
        />
        {type === 'policy' && (
          <React.Fragment>
            <Field
              as={FormikMultiCombobox}
              searchable
              name="suppressions"
              label="Resource Ignore Patterns"
              items={values.suppressions}
              allowAdditions
              inputProps={{
                placeholder: 'i.e. aws::s3::* (separate with <Enter>)',
              }}
            />
            <Box>
              <Field
                as={FormikMultiCombobox}
                searchable
                label="Resource Types"
                name="resourceTypes"
                items={RESOURCE_TYPES}
                inputProps={{ placeholder: 'Filter affected resource types' }}
              />
              <Text size="small" color="grey300" mt={2}>
                Leave empty to apply to all resources
              </Text>
            </Box>
          </React.Fragment>
        )}
        <Field
          as={FormikMultiCombobox}
          searchable
          name="tags"
          label="Custom Tags"
          items={values.tags}
          allowAdditions
          validateAddition={tag => !values.tags.includes(tag)}
          inputProps={{
            placeholder: 'i.e. Bucket Security (separate with <Enter>)',
          }}
        />
        {type === 'rule' && (
          <Box>
            <Field
              as={FormikMultiCombobox}
              searchable
              label="Log Types"
              name="logTypes"
              items={LOG_TYPES}
              inputProps={{ placeholder: 'Filter affected log types' }}
            />
            <Text size="small" color="grey300" mt={2}>
              Leave empty to apply to all logs
            </Text>
          </Box>
        )}
      </Grid>
      <Box my={6}>
        <InputElementLabel htmlFor="enabled">{`* ${capitalize(type)} Function`}</InputElementLabel>
        <Field
          as={FormikEditor}
          placeholder={`# Enter the body of the ${type} here...`}
          name="body"
          width="100%"
          minLines={16}
          mode="python"
          aria-required
        />
        {errors.body && touched.body && (
          <InputElementErrorLabel mt={6}>{errors.body}</InputElementErrorLabel>
        )}
      </Box>
    </section>
  );
};

export default BaseRuleCoreFields;
