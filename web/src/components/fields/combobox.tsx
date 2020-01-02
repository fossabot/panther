import React from 'react';
import { Combobox, ComboboxProps } from 'pouncejs';
import { useFormikContext, FieldConfig } from 'formik';

function FormikCombobox<T>(
  props: ComboboxProps<T> & Required<Pick<FieldConfig, 'name'>>
): React.ReactNode {
  const { setFieldValue } = useFormikContext<any>();
  return <Combobox<T> {...props} onChange={value => setFieldValue(props.name, value)} />;
}

export default FormikCombobox;
