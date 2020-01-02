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
