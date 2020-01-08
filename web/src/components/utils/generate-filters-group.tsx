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
import { Formik, Field } from 'formik';
import { Grid, Button } from 'pouncejs';
import mapValues from 'lodash-es/mapValues';
import map from 'lodash-es/map';
import FormikTextInput from 'Components/fields/text-input';
import FormikTextArea from 'Components/fields/textarea';
import FormikCombobox from 'Components/fields/combobox';
import FormikMultiCombobox from 'Components/fields/multi-combobox';

interface FiltersGroupData<T> {
  /** The component to render for this particular form entry */
  component:
    | typeof FormikCombobox
    | typeof FormikMultiCombobox
    | typeof FormikTextInput
    | typeof FormikTextArea;

  /** The props that should be given to the form entry */
  props: { [key: string]: any };
}

interface GenerateFiltersGroupProps<T> {
  /** The initial values of the filters group */
  initialValues: T;

  /** A callback for when the `cancel` button is clicked */
  onCancel: () => void;

  /** A callback for when the `apply` button is clicked or - in general - the filters are applied */
  onSubmit: (values: T) => void;

  /** A filter configuration */
  filters: {
    [K in Extract<keyof T, string>]: FiltersGroupData<T>;
  };
}

const getFilterGroupDefaultValue = (
  component: FiltersGroupData<any>['component']
): [] | '' | null => {
  switch (component) {
    case FormikCombobox:
      return null;
    case FormikMultiCombobox:
      return [];
    default:
      return '';
  }
};

function GenerateFiltersGroup<T extends { [key: string]: any }>({
  filters,
  onCancel,
  onSubmit,
  initialValues,
}: GenerateFiltersGroupProps<T>): React.ReactElement<GenerateFiltersGroupProps<T>> {
  // These are the default values that each field should have. This is related to the type of the
  // field (a.k.a. `component`) and is the "fallback initial value of the field" (since its actual
  // initial value comes from the URL)
  const defaultValues = React.useMemo(() => {
    return mapValues(filters, filterData => getFilterGroupDefaultValue(filterData.component)) as T;
  }, []);

  // We initialize the values of the form based on the current URL. This only happens during mount
  // time. The value of `initialValues` doesn't get updated as the component updates, since we only
  // need it during form initialization (a.k.a. component mount-time)
  const initialValuesWithDefaults = React.useMemo(() => {
    return mapValues(filters, (value, name) => initialValues[name] || defaultValues[name]) as T;
  }, []);

  // On a successful submit, the URL params are updated and the page query gets re-fetched, since
  // the page query depends on what the URL is. Essentially, we are using the URL params as a
  // "store" that we observe on the index
  return (
    <Formik<T> initialValues={initialValuesWithDefaults} onSubmit={onSubmit}>
      {({ handleSubmit, setValues, submitForm, resetForm }) => (
        <form onSubmit={handleSubmit}>
          <Grid gridTemplateColumns="repeat(3, 1fr)" gridGap={6} mb={8}>
            {map(filters, (filterData, filterName) => (
              <Field
                key={filterName}
                as={filterData.component}
                name={filterName}
                {...filterData.props}
              />
            ))}
          </Grid>
          <Button type="submit" size="large" variant="primary" mr={4}>
            Apply
          </Button>
          <Button
            type="button"
            size="large"
            variant="default"
            mr={4}
            onClick={() => {
              resetForm();
              onCancel();
            }}
          >
            Cancel
          </Button>
          <Button
            type="button"
            size="large"
            variant="default"
            color="red300"
            onClick={() => {
              setValues(defaultValues);
              submitForm();
            }}
          >
            Clear all
          </Button>
        </form>
      )}
    </Formik>
  );
}

export default GenerateFiltersGroup;
