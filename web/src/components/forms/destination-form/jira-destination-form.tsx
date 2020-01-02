import React from 'react';
import { Field } from 'formik';
import * as Yup from 'yup';
import FormikTextInput from 'Components/fields/text-input';
import { DestinationConfigInput } from 'Generated/schema';
import BaseDestinationForm, {
  BaseDestinationFormValues,
  defaultValidationSchema,
} from 'Components/forms/common/base-destination-form';

type JiraFieldValues = Pick<DestinationConfigInput, 'jira'>;

interface JiraDestinationFormProps {
  initialValues: BaseDestinationFormValues<JiraFieldValues>;
  onSubmit: (values: BaseDestinationFormValues<JiraFieldValues>) => void;
}

const jiraFieldsValidationSchema = Yup.object().shape({
  outputConfig: Yup.object().shape({
    jira: Yup.object().shape({
      orgDomain: Yup.string()
        .url('Must be a valid Jira domain')
        .required(),
      userName: Yup.string(),
      projectKey: Yup.string().required(),
      apiKey: Yup.string().required(),
      assigneeId: Yup.string(),
    }),
  }),
});

// @ts-ignore
// We merge the two schemas together: the one deriving from the common fields, plus the custom
// ones that change for each destination.
// https://github.com/jquense/yup/issues/522
const mergedValidationSchema = defaultValidationSchema.concat(jiraFieldsValidationSchema);

const JiraDestinationForm: React.FC<JiraDestinationFormProps> = ({ onSubmit, initialValues }) => {
  return (
    <BaseDestinationForm<JiraFieldValues>
      initialValues={initialValues}
      validationSchema={mergedValidationSchema}
      onSubmit={onSubmit}
    >
      <Field
        as={FormikTextInput}
        name="outputConfig.jira.orgDomain"
        label="Organization Domain"
        placeholder="What's your organization's Jira domain?"
        mb={6}
        aria-required
      />
      <Field
        as={FormikTextInput}
        name="outputConfig.jira.projectKey"
        label="Project Key"
        placeholder="What's your Jira Project key?"
        mb={6}
        aria-required
        autoComplete="new-password"
      />
      <Field
        as={FormikTextInput}
        name="outputConfig.jira.userName"
        label="User Name"
        placeholder="What's the name of the reporting user?"
        mb={6}
      />
      <Field
        as={FormikTextInput}
        name="outputConfig.jira.apiKey"
        label="Jira API Key"
        placeholder="What's the API key of the related Jira account"
        mb={6}
        aria-required
        autoComplete="new-password"
      />

      <Field
        as={FormikTextInput}
        name="outputConfig.jira.assigneeID"
        label="Assignee ID"
        placeholder="Who should we assign this to?"
        mb={6}
      />
    </BaseDestinationForm>
  );
};

export default JiraDestinationForm;
