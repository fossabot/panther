import React from 'react';
import { Field } from 'formik';
import * as Yup from 'yup';
import FormikTextInput from 'Components/fields/text-input';
import { DestinationConfigInput } from 'Generated/schema';
import BaseDestinationForm, {
  BaseDestinationFormValues,
  defaultValidationSchema,
} from 'Components/forms/common/base-destination-form';

type PagerDutyFieldValues = Pick<DestinationConfigInput, 'pagerDuty'>;

interface PagerDutyDestinationFormProps {
  initialValues: BaseDestinationFormValues<PagerDutyFieldValues>;
  onSubmit: (values: BaseDestinationFormValues<PagerDutyFieldValues>) => void;
}

const pagerDutyFieldsValidationSchema = Yup.object().shape({
  outputConfig: Yup.object().shape({
    pagerDuty: Yup.object().shape({
      integrationKey: Yup.string()
        .length(32, 'Must be exactly 32 characters')
        .required(),
    }),
  }),
});

// @ts-ignore
// We merge the two schemas together: the one deriving from the common fields, plus the custom
// ones that change for each destination.
// https://github.com/jquense/yup/issues/522
const mergedValidationSchema = defaultValidationSchema.concat(pagerDutyFieldsValidationSchema);

const PagerDutyDestinationForm: React.FC<PagerDutyDestinationFormProps> = ({
  onSubmit,
  initialValues,
}) => {
  return (
    <BaseDestinationForm<PagerDutyFieldValues>
      initialValues={initialValues}
      validationSchema={mergedValidationSchema}
      onSubmit={onSubmit}
    >
      <Field
        as={FormikTextInput}
        name="outputConfig.pagerDuty.integrationKey"
        label="Integration Key"
        placeholder="What's your PagerDuty Integration Key?"
        mb={6}
        aria-required
        autoComplete="new-password"
      />
    </BaseDestinationForm>
  );
};

export default PagerDutyDestinationForm;
