import React from 'react';
import { WizardRenderStepParams } from 'Components/wizard';
import { Box, Button, Flex } from 'pouncejs';

interface PanelWrapperWizardActionsProps {
  goToPrevStep?: WizardRenderStepParams<{}>['goToPrevStep'];
  goToNextStep?: WizardRenderStepParams<{}>['goToNextStep'];
  isNextStepDisabled?: boolean;
  isPrevStepDisabled?: boolean;
}

interface PanelWrapperComposition {
  WizardActions: React.FC<PanelWrapperWizardActionsProps>;
  Content: React.FC;
}

const PanelWrapper: React.FC & PanelWrapperComposition = ({ children }) => {
  return (
    <Flex height={550} flexDirection="column">
      {children}
    </Flex>
  );
};

const PanelWrapperContent: React.FC = ({ children }) => {
  return (
    <Box width={600} m="auto">
      {children}
    </Box>
  );
};

const PanelWrapperWizardActions: React.FC<PanelWrapperWizardActionsProps> = ({
  isNextStepDisabled,
  isPrevStepDisabled,
  goToPrevStep,
  goToNextStep,
}) => {
  return (
    <Flex justifyContent="flex-end">
      {goToPrevStep && (
        <Button
          size="large"
          variant="default"
          onClick={goToPrevStep}
          mr={3}
          disabled={isPrevStepDisabled}
        >
          Back
        </Button>
      )}
      {goToNextStep && (
        <Button size="large" variant="primary" onClick={goToNextStep} disabled={isNextStepDisabled}>
          Next
        </Button>
      )}
    </Flex>
  );
};

PanelWrapper.Content = PanelWrapperContent;
PanelWrapper.WizardActions = PanelWrapperWizardActions;

export default PanelWrapper;
