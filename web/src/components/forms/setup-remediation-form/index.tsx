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
import * as Yup from 'yup';
import { FastField as Field, Formik } from 'formik';
import { Box, Button, Flex, InputElementLabel } from 'pouncejs';
import { AWS_ACCOUNT_ID_REGEX } from 'Source/constants';
import FormikCheckbox from 'Components/fields/checkbox';
import FormikTextInput from 'Components/fields/text-input';

interface SetupRemediationFormValues {
  isSatellite: boolean;
  adminAWSAccountId: string;
}

interface SetupRemediationFormProps {
  getStackUrl: (values: SetupRemediationFormValues) => string;
  onStackLaunch?: () => void;
}

const initialValues = {
  isSatellite: false,
  adminAWSAccountId: '',
};

const validationSchema = Yup.object().shape({
  isSatellite: Yup.boolean(),
  username: Yup.string()
    .matches(AWS_ACCOUNT_ID_REGEX)
    .required(),
});

const SetupRemediationForm: React.FC<SetupRemediationFormProps> = ({
  getStackUrl,
  onStackLaunch,
}) => {
  return (
    <Formik<SetupRemediationFormValues>
      initialValues={initialValues}
      onSubmit={() => {}}
      validationSchema={validationSchema}
    >
      {({ handleSubmit, values }) => (
        <Box is="form" onSubmit={handleSubmit}>
          <Flex mb={6} alignItems="center">
            <InputElementLabel htmlFor="isSatellite" mr={3}>
              I want to setup automatic remediation in a satellite account
            </InputElementLabel>
            <Field as={FormikCheckbox} id="isSatellite" name="isSatellite" />
          </Flex>
          <Box hidden={!values.isSatellite} mb={10} width={0.3}>
            <Field
              as={FormikTextInput}
              label="Your Auto-Remediation Master AWS Account ID"
              name="adminAWSAccountId"
              placeholder="i.e. 548784460855"
              aria-required
            />
          </Box>
          <Button
            size="large"
            variant="default"
            target="_blank"
            is="a"
            rel="noopener noreferrer"
            href={getStackUrl(values)}
            onClick={onStackLaunch}
          >
            Launch Stack
          </Button>
        </Box>
      )}
    </Formik>
  );
};

export default React.memo(SetupRemediationForm);
