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
import { Alert, Box, Combobox, Grid, InputElementLabel, Spinner } from 'pouncejs';
import { Field, useFormikContext } from 'formik';
import FormikTextInput from 'Components/fields/text-input';
import { formatJSON, extractErrorMessage } from 'Helpers/utils';
import { useQuery, gql } from '@apollo/client';
import FormikEditor from 'Components/fields/editor';
import { PolicyFormValues } from './index';

export const LIST_REMEDIATIONS = gql`
  query ListRemediations {
    remediations
  }
`;

interface ApolloQueryData {
  remediations: string;
}

const PolicyFormAutoRemediationFields: React.FC = () => {
  // Read the values from the "parent" form. We expect a formik to be declared in the upper scope
  // since this is a "partial" form. If no Formik context is found this will error out intentionally
  const { values, setFieldValue } = useFormikContext<PolicyFormValues>();

  // This state is used to track/store the value of the auto-remediation combobox. This combobox
  // doesn't belong to the form and we wouldn't wanna pollute our form with undesired information.
  // Instead what this checkbox does, is to control the value of the actual fields in the form which
  // are the ID and Params of the auto remediation.
  // Here we are parsing & reformatting for display purposes only (since the JSON that arrives as a
  // string doesn't have any formatting)
  const [autoRemediationSelection, setAutoRemediationSelection] = React.useState<[string, string]>([
    values.autoRemediationId,
    values.autoRemediationParameters,
  ]);

  // Currently there is a bug in apollo. On requests where the cache is read, if there
  // is an error, then the second time you read from the cache the "error" key is undefined. Thus,
  // you don't know whether there was actually an error before. For this reason, we don't add cache
  // for requests where we would want to still store the error in the cache (since the error would
  // mean that no remediation lambda is present). That's why `no-cache` is added here.
  // https://github.com/apollographql/apollo-client/issues/4138
  // TODO: convert fetchPolicy to `cache-first` if the above issue is resolved
  const { data, loading, error } = useQuery<ApolloQueryData>(LIST_REMEDIATIONS, {
    fetchPolicy: 'no-cache',
  });

  if (loading) {
    return <Spinner size="medium" />;
  }

  if (error) {
    return (
      <Alert
        variant="warning"
        title="Couldn't load your available remediations"
        description={[
          extractErrorMessage(error),
          '. For more info, please consult the ',
          <a
            key="docs"
            href="https://docs.runpanther.io/amazon-web-services/aws-setup/automatic-remediation"
            target="_blank"
            rel="noopener noreferrer"
          >
            related docs
          </a>,
        ]}
      />
    );
  }

  const remediationTuples = Object.entries(
    JSON.parse(data.remediations)
  ).map(([id, params]: [string, { [key: string]: string }]) => [id, formatJSON(params)]) as [
    string,
    string
  ][];

  return (
    <section>
      <Grid gridTemplateColumns="1fr 1fr" gridRowGap={2} gridColumnGap={9}>
        <Combobox<[string, string]>
          searchable
          label="Remediation"
          items={[['', '{}'], ...remediationTuples]}
          itemToString={remediationTuple => remediationTuple[0] || '(No remediation)'}
          value={autoRemediationSelection}
          onChange={remediationTuple => {
            setFieldValue('autoRemediationId', remediationTuple[0]);
            setFieldValue('autoRemediationParameters', remediationTuple[1]);
            setAutoRemediationSelection(remediationTuple);
          }}
        />
      </Grid>
      <Box hidden>
        <Field as={FormikTextInput} name="autoRemediationId" />
      </Box>
      <Box mt={10} hidden={!values.autoRemediationId}>
        <InputElementLabel htmlFor="enabled">Remediation Parameters</InputElementLabel>
        <Field
          as={FormikEditor}
          placeholder="# Enter a JSON object describing the parameters of the remediation"
          name="autoRemediationParameters"
          width="100%"
          minLines={9}
          mode="json"
        />
      </Box>
    </section>
  );
};

export default React.memo(PolicyFormAutoRemediationFields);
