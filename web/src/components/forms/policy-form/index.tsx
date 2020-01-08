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
import { PolicyDetails, PolicyUnitTest } from 'Generated/schema';
import * as Yup from 'yup';
import { Box, Heading } from 'pouncejs';
import BaseRuleForm, { BaseRuleFormProps } from 'Components/forms/common/base-rule-form';
import ErrorBoundary from 'Components/error-boundary';
import PolicyFormAutoRemediationFields from './policy-form-auto-remediation-fields';
import RuleFormCoreFields, { ruleCoreEditableFields } from '../common/rule-form-core-fields';
import PolicyFormTestFields from '../common/rule-form-test-fields';

export const policyEditableFields = [
  ...ruleCoreEditableFields,
  'autoRemediationId',
  'autoRemediationParameters',
  'suppressions',
  'resourceTypes',
  'tests',
] as const;

// The validation checks that Formik will run
const validationSchema = Yup.object().shape({
  id: Yup.string().required(),
  body: Yup.string().required(),
  severity: Yup.string().required(),
  tests: Yup.array<PolicyUnitTest>()
    .of(
      Yup.object().shape({
        name: Yup.string().required(),
      })
    )
    .unique('Test names must be unique', 'name'),
});

export type PolicyFormValues = Pick<PolicyDetails, typeof policyEditableFields[number]>;
export type PolicyFormProps = Pick<
  BaseRuleFormProps<PolicyFormValues>,
  'initialValues' | 'onSubmit'
>;

const PolicyForm: React.FC<PolicyFormProps> = ({ initialValues, onSubmit }) => {
  return (
    <BaseRuleForm<PolicyFormValues>
      initialValues={initialValues}
      onSubmit={onSubmit}
      validationSchema={validationSchema}
    >
      <Box is="article">
        <ErrorBoundary>
          <RuleFormCoreFields type="policy" />
        </ErrorBoundary>
        <ErrorBoundary>
          <PolicyFormTestFields />
        </ErrorBoundary>
      </Box>
      <Box is="article" mt={10}>
        <Heading size="medium" pb={8} borderBottom="1px solid" borderColor="grey100">
          Auto Remediation Settings
        </Heading>
        <Box mt={8}>
          <ErrorBoundary>
            <PolicyFormAutoRemediationFields />
          </ErrorBoundary>
        </Box>
      </Box>
    </BaseRuleForm>
  );
};

export default PolicyForm;
