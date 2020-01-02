import React from 'react';
import { TextArea, TextAreaProps } from 'pouncejs';
import { FieldConfig, useField } from 'formik';

const FormikTextArea: React.FC<TextAreaProps & Required<Pick<FieldConfig, 'name'>>> = props => {
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  const [__, meta] = useField(props.name);
  return <TextArea {...props} error={meta.touched && meta.error} />;
};

export default FormikTextArea;
