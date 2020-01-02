import React from 'react';
import { Field } from 'formik';
import * as Yup from 'yup';
import FormikTextInput from 'Components/fields/text-input';
import { DestinationConfigInput } from 'Generated/schema';
import BaseDestinationForm, {
  BaseDestinationFormValues,
  defaultValidationSchema,
} from 'Components/forms/common/base-destination-form';

type GithubFieldValues = Pick<DestinationConfigInput, 'github'>;

interface GithubDestinationFormProps {
  initialValues: BaseDestinationFormValues<GithubFieldValues>;
  onSubmit: (values: BaseDestinationFormValues<GithubFieldValues>) => void;
}

const githubFieldsValidationSchema = Yup.object().shape({
  outputConfig: Yup.object().shape({
    github: Yup.object().shape({
      repoName: Yup.string().required(),
      token: Yup.string().required(),
    }),
  }),
});

// @ts-ignore
// We merge the two schemas together: the one deriving from the common fields, plus the custom
// ones that change for each destination.
// https://github.com/jquense/yup/issues/522
const mergedValidationSchema = defaultValidationSchema.concat(githubFieldsValidationSchema);

const GithubDestinationForm: React.FC<GithubDestinationFormProps> = ({
  onSubmit,
  initialValues,
}) => {
  return (
    <BaseDestinationForm<GithubFieldValues>
      initialValues={initialValues}
      validationSchema={mergedValidationSchema}
      onSubmit={onSubmit}
    >
      <Field
        as={FormikTextInput}
        name="outputConfig.github.repoName"
        label="Repository name"
        placeholder="What's the name of your Github repository?"
        mb={6}
        aria-required
      />
      <Field
        as={FormikTextInput}
        name="outputConfig.github.token"
        label="Token"
        placeholder="What's your Github API token?"
        mb={6}
        aria-required
        autoComplete="new-password"
      />
    </BaseDestinationForm>
  );
};

export default GithubDestinationForm;
