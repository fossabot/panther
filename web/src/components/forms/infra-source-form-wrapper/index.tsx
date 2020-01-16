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
import { Formik, FormikProps } from 'formik';
import * as Yup from 'yup';

export interface InfraSourceFormWrapperValues {
  awsAccountId: string;
  integrationLabel: string;
}

export interface InfraSourceFormWrapperProps {
  initialValues: InfraSourceFormWrapperValues;
  onSubmit: (values: InfraSourceFormWrapperValues) => Promise<any> | void;
  children: (renderProps: Partial<FormikProps<InfraSourceFormWrapperValues>>) => React.ReactElement;
}

const validationSchema = Yup.object().shape({
  awsAccountId: Yup.string()
    .matches(/[0-9]+/, 'Must only contain numbers')
    .length(12)
    .required(),
  integrationLabel: Yup.string().required(),
});

const InfraSourceFormWrapper = ({
  onSubmit,
  initialValues,
  children,
}: InfraSourceFormWrapperProps): React.ReactElement => {
  return (
    <Formik<InfraSourceFormWrapperValues>
      initialValues={initialValues}
      validationSchema={validationSchema}
      onSubmit={onSubmit}
    >
      {({ handleSubmit, ...rest }) => <form onSubmit={handleSubmit}>{children(rest)}</form>}
    </Formik>
  );
};

export default InfraSourceFormWrapper;
