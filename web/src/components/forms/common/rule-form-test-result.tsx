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
import { ComplianceStatusEnum } from 'Generated/schema';
import { Card, defaultTheme, Flex, Label } from 'pouncejs';

// A mapping from status to background color for our test results (background color of where it says
// 'pass', 'fail' or 'error'
export const mapTestStatusToColor: {
  [key in ComplianceStatusEnum]: keyof typeof defaultTheme['colors'];
} = {
  [ComplianceStatusEnum.Pass]: 'green200',
  [ComplianceStatusEnum.Fail]: 'red300',
  [ComplianceStatusEnum.Error]: 'orange300',
};

interface TestResultProps {
  /** The name of the test */
  testName: string;

  /** The value that is going to displayed to the user as a result for this test */
  status: ComplianceStatusEnum;
}

const TestResult: React.FC<TestResultProps> = ({ testName, status }) => (
  <Flex alignItems="center">
    <Card bg={mapTestStatusToColor[status]} mr={2} width={50} py={1}>
      <Label size="small" color="white" mx="auto" is="div" textAlign="center">
        {status}
      </Label>
    </Card>
    <Label size="medium" color="grey400">
      {testName}
    </Label>
  </Flex>
);

export default TestResult;
