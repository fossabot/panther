import React from 'react';
import { MultiCombobox, MultiComboboxProps } from 'pouncejs';
import { useFormikContext, FieldConfig } from 'formik';

function FormikMultiCombobox<T>(
  props: MultiComboboxProps<T> & Required<Pick<FieldConfig, 'name'>>
): React.ReactNode {
  const { setFieldValue } = useFormikContext<any>();
  return <MultiCombobox<T> {...props} onChange={value => setFieldValue(props.name, value)} />;
}

export default FormikMultiCombobox;
