import React from 'react';
import { Field } from 'formik';
import * as Yup from 'yup';
import FormikTextInput from 'Components/fields/text-input';
import { DestinationConfigInput } from 'Generated/schema';
import BaseDestinationForm, {
  BaseDestinationFormValues,
  defaultValidationSchema,
} from 'Components/forms/common/base-destination-form';

type SlackFieldValues = Pick<DestinationConfigInput, 'slack'>;

interface SlackDestinationFormProps {
  initialValues: BaseDestinationFormValues<SlackFieldValues>;
  onSubmit: (values: BaseDestinationFormValues<SlackFieldValues>) => void;
}

const slackFieldsValidationSchema = Yup.object().shape({
  outputConfig: Yup.object().shape({
    slack: Yup.object().shape({
      webhookURL: Yup.string()
        .url('Must be a valid webhook URL')
        .required(),
    }),
  }),
});

// @ts-ignore
// We merge the two schemas together: the one deriving from the common fields, plus the custom
// ones that change for each destination.
// https://github.com/jquense/yup/issues/522
const mergedValidationSchema = defaultValidationSchema.concat(slackFieldsValidationSchema);

const SlackDestinationForm: React.FC<SlackDestinationFormProps> = ({ onSubmit, initialValues }) => {
  return (
    <BaseDestinationForm<SlackFieldValues>
      initialValues={initialValues}
      validationSchema={mergedValidationSchema}
      onSubmit={onSubmit}
    >
      <Field
        as={FormikTextInput}
        name="outputConfig.slack.webhookURL"
        label="Slack Webhook URL"
        placeholder="Where should we send a push notification to?"
        mb={6}
        aria-required
      />
    </BaseDestinationForm>
  );
};

export default SlackDestinationForm;
