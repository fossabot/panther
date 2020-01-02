import React from 'react';
import { Formik, FormikProps } from 'formik';
import * as Yup from 'yup';

export interface LogSourceFormWrapperValues {
  integrationLabel: string;
  sourceSnsTopicArn: string;
  logProcessingRoleArn: string;
}

export interface LogSourceFormWrapperProps {
  initialValues: LogSourceFormWrapperValues;
  onSubmit: (values: LogSourceFormWrapperValues) => Promise<any> | void;
  children: (renderProps: Partial<FormikProps<LogSourceFormWrapperValues>>) => React.ReactElement;
}

const validationSchema = Yup.object().shape({
  integrationLabel: Yup.string().required(),
  sourceSnsTopicArn: Yup.string().required(),
  logProcessingRoleArn: Yup.string().required(),
});

export const LogSourceFormWrapper = ({
  onSubmit,
  initialValues,
  children,
}: LogSourceFormWrapperProps): React.ReactElement => {
  return (
    <Formik<LogSourceFormWrapperValues>
      initialValues={initialValues}
      validationSchema={validationSchema}
      onSubmit={onSubmit}
    >
      {({ handleSubmit, ...rest }) => <form onSubmit={handleSubmit}>{children(rest)}</form>}
    </Formik>
  );
};

export default LogSourceFormWrapper;
