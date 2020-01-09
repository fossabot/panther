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
import { Formik } from 'formik';
import SubmitButton from 'Components/utils/SubmitButton';
import * as Yup from 'yup';
import { Flex, Button } from 'pouncejs';
import useRouter from 'Hooks/useRouter';

interface IdValue {
  id: string;
}

export interface BaseRuleFormProps<BaseRuleFormValues extends IdValue> {
  /** The initial values of the form */
  initialValues: BaseRuleFormValues;

  /** callback for the submission of the form */
  onSubmit: (values: BaseRuleFormValues) => void;

  /** The validation schema that the form will have */
  validationSchema: Yup.ObjectSchema<Yup.Shape<object, Partial<BaseRuleFormValues> & IdValue>>;
}

function BaseRuleForm<BaseRuleFormValues extends IdValue>({
  initialValues,
  onSubmit,
  validationSchema,
  children,
}: React.PropsWithChildren<BaseRuleFormProps<BaseRuleFormValues>>): React.ReactElement<
  BaseRuleFormProps<BaseRuleFormValues>
> {
  const { history } = useRouter();

  return (
    <Formik<BaseRuleFormValues>
      initialValues={initialValues}
      onSubmit={onSubmit}
      enableReinitialize
      validationSchema={validationSchema}
    >
      {({ handleSubmit, isSubmitting, isValid, dirty }) => {
        return (
          <form onSubmit={handleSubmit}>
            {children}
            <Flex
              borderTop="1px solid"
              borderColor="grey100"
              pt={6}
              mt={10}
              justifyContent="flex-end"
            >
              <Flex>
                <Button variant="default" size="large" onClick={history.goBack} mr={4}>
                  Cancel
                </Button>
                <SubmitButton
                  submitting={isSubmitting}
                  disabled={!dirty || !isValid || isSubmitting}
                >
                  {initialValues.id ? 'Update' : 'Create'}
                </SubmitButton>
              </Flex>
            </Flex>
          </form>
        );
      }}
    </Formik>
  );
}

export default BaseRuleForm;