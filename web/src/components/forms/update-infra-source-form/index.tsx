/**
 * Panther is a scalable, powerful, cloud-native SIEM written in Golang/React.
 * Copyright (C) 2020 Panther Labs Inc
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

import React from 'react';
import { Field } from 'formik';
import FormikTextInput from 'Components/fields/text-input';
import SubmitButton from 'Components/utils/SubmitButton';
import InfraSourceFormWrapper, {
  InfraSourceFormWrapperProps,
} from 'Components/forms/infra-source-form-wrapper';

export type UpdateInfraSourceFormProps = Omit<InfraSourceFormWrapperProps, 'children'>;

const UpdateInfraSourceForm: React.FC<UpdateInfraSourceFormProps> = ({
  onSubmit,
  initialValues,
}) => {
  return (
    <InfraSourceFormWrapper initialValues={initialValues} onSubmit={onSubmit}>
      {({ isSubmitting, isValid, dirty }) => (
        <React.Fragment>
          <Field
            name="awsAccountId"
            as={FormikTextInput}
            label="AWS Account ID"
            placeholder="Your 12-digit AWS Account ID"
            aria-required
            mb={6}
          />
          <Field
            name="integrationLabel"
            as={FormikTextInput}
            label="Label"
            placeholder="A nickname for your account"
            aria-required
            mb={6}
          />
          <SubmitButton
            width={1}
            disabled={isSubmitting || !isValid || !dirty}
            submitting={isSubmitting}
          >
            {initialValues.awsAccountId ? 'Update' : 'Add'}
          </SubmitButton>
        </React.Fragment>
      )}
    </InfraSourceFormWrapper>
  );
};

export default UpdateInfraSourceForm;
