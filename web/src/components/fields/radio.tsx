import React from 'react';
import { Radio, RadioProps } from 'pouncejs';
import { useFormikContext, FieldConfig, useField } from 'formik';

const FormikRadio: React.FC<RadioProps & Required<Pick<FieldConfig, 'name' | 'value'>>> = props => {
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  const [field, meta] = useField(props.name);
  const { setFieldValue } = useFormikContext<any>();

  // Here `props.value` is the value that the radio button should have according to the typical HTML
  // and not the value that will be forced into Formik
  return (
    <Radio
      {...props}
      checked={field.value === props.value}
      onChange={() => setFieldValue(field.name, props.value)}
    />
  );
};

export default FormikRadio;
