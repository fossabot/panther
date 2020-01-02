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
