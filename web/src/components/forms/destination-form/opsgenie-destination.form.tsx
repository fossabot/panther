import React from 'react';
import { Field } from 'formik';
import * as Yup from 'yup';
import FormikTextInput from 'Components/fields/text-input';
import { DestinationConfigInput } from 'Generated/schema';
import BaseDestinationForm, {
  BaseDestinationFormValues,
  defaultValidationSchema,
} from 'Components/forms/common/base-destination-form';

type OpsgenieFieldValues = Pick<DestinationConfigInput, 'opsgenie'>;

interface OpsgenieDestinationFormProps {
  initialValues: BaseDestinationFormValues<OpsgenieFieldValues>;
  onSubmit: (values: BaseDestinationFormValues<OpsgenieFieldValues>) => void;
}

const opsgenieFieldsValidationSchema = Yup.object().shape({
  outputConfig: Yup.object().shape({
    opsgenie: Yup.object().shape({
      apiKey: Yup.string().required(),
    }),
  }),
});

// @ts-ignore
// We merge the two schemas together: the one deriving from the common fields, plus the custom
// ones that change for each destination.
// https://github.com/jquense/yup/issues/522
const mergedValidationSchema = defaultValidationSchema.concat(opsgenieFieldsValidationSchema);

const OpsgenieDestinationForm: React.FC<OpsgenieDestinationFormProps> = ({
  onSubmit,
  initialValues,
}) => {
  return (
    <BaseDestinationForm<OpsgenieFieldValues>
      initialValues={initialValues}
      validationSchema={mergedValidationSchema}
      onSubmit={onSubmit}
    >
      <Field
        as={FormikTextInput}
        name="outputConfig.opsgenie.apiKey"
        label="Opsgenie API key"
        placeholder="What's your organization's Opsgenie API key?"
        mb={6}
        aria-required
        autoComplete="new-password"
      />
    </BaseDestinationForm>
  );
};

export default OpsgenieDestinationForm;
