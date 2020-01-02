import React from 'react';
import Editor, { EditorProps } from 'Components/editor';
import { useFormikContext, FieldConfig } from 'formik';
import debounce from 'lodash-es/debounce';

const FormikEditor: React.FC<EditorProps & Required<Pick<FieldConfig, 'name'>>> = ({
  // we destruct `onBlur` since we shouldn't pass it as a prop to `Editor`. This is becase we are
  // manually syncing the changes of the editor to the formik instance through the
  // `syncValueFromEditor`. Thus, we don't need an `onBlur`
  // eslint-disable-next-line @typescript-eslint/no-unused-vars
  onBlur,
  ...rest
}) => {
  const { setFieldValue } = useFormikContext<any>();

  // For performance enhancing reasons, we are debouncing the syncing of the editor value to
  // the formik controller. The editor is *not* a controlled component by nature, so we are
  // only syncing its internal state to formik with some delays.
  // It's worth noting that this is contradictory to all the other components in the `fields`
  // folder, since they are controlled
  const syncValueFromEditor = React.useCallback(
    debounce((value: string) => {
      setFieldValue(rest.name, value);
    }, 200),
    [setFieldValue, rest.name]
  );

  return <Editor {...rest} onChange={syncValueFromEditor} />;
};

export default FormikEditor;
