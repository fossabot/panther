/**
 * Copyright 2020 Panther Labs Inc
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *    http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

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
