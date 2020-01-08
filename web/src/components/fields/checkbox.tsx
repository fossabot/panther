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
import { Checkbox, CheckboxProps } from 'pouncejs';
import { useFormikContext, FieldConfig, useField } from 'formik';

const FormikCheckbox: React.FC<CheckboxProps & Required<Pick<FieldConfig, 'name'>>> = props => {
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  const [field, meta] = useField<boolean>(props.name);
  const { setFieldValue } = useFormikContext<any>();
  return (
    <Checkbox
      {...props}
      checked={field.value}
      onChange={value => setFieldValue(field.name, value)}
    />
  );
};

export default FormikCheckbox;
