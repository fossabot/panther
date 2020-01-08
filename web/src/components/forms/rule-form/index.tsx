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
import { RuleDetails, PolicyUnitTest } from 'Generated/schema';
import * as Yup from 'yup';
import { Box } from 'pouncejs';
import ErrorBoundary from 'Components/error-boundary';
import BaseRuleForm, { BaseRuleFormProps } from 'Components/forms/common/base-rule-form';
import RuleFormCoreFields, { ruleCoreEditableFields } from '../common/rule-form-core-fields';
import RuleFormTestFields from '../common/rule-form-test-fields';

export const ruleEditableFields = [...ruleCoreEditableFields, 'logTypes', 'tests'] as const;

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

export type RuleFormValues = Pick<RuleDetails, typeof ruleEditableFields[number]>;
export type RuleFormProps = Pick<BaseRuleFormProps<RuleFormValues>, 'initialValues' | 'onSubmit'>;

const RuleForm: React.FC<RuleFormProps> = ({ initialValues, onSubmit }) => {
  return (
    <BaseRuleForm<RuleFormValues>
      initialValues={initialValues}
      onSubmit={onSubmit}
      validationSchema={validationSchema}
    >
      <Box is="article">
        <ErrorBoundary>
          <RuleFormCoreFields type="rule" />
        </ErrorBoundary>
        <ErrorBoundary>
          <RuleFormTestFields />
        </ErrorBoundary>
      </Box>
    </BaseRuleForm>
  );
};

export default RuleForm;
