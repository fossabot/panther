import React from 'react';
import { Switch, SwitchProps } from 'pouncejs';
import { useFormikContext, FieldConfig, useField } from 'formik';

const FormikSwitch: React.FC<SwitchProps & Required<Pick<FieldConfig, 'name'>>> = props => {
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  const [field, meta] = useField<boolean>(props.name);
  const { setFieldValue } = useFormikContext<any>();
  return (
    <Switch {...props} checked={field.value} onChange={value => setFieldValue(field.name, value)} />
  );
};

export default FormikSwitch;
