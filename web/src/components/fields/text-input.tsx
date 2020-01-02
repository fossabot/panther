import React from 'react';
import { TextInput, TextInputProps } from 'pouncejs';
import { FieldConfig, useField } from 'formik';

const FormikTextInput: React.FC<TextInputProps & Required<Pick<FieldConfig, 'name'>>> = props => {
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  const [field, meta] = useField(props.name);
  return <TextInput {...props} error={meta.touched && meta.error} />;
};

export default FormikTextInput;
