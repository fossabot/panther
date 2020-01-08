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
import { Box, Label } from 'pouncejs';
import { ComplianceStatusEnum, TestPolicyResponse } from 'Generated/schema';
import PolicyFormTestResult, { mapTestStatusToColor } from './rule-form-test-result';

interface PolicyFormTestResultsProps {
  results: TestPolicyResponse;
  running: boolean;
}

const RuleFormTestResultList: React.FC<PolicyFormTestResultsProps> = ({ running, results }) => {
  return (
    <Box bg="#FEF5ED" p={5}>
      {running && (
        <Label size="medium" is="p">
          Running your tests...
        </Label>
      )}
      {!running && results && (
        <React.Fragment>
          {results.testsPassed.map(testName => (
            <Box mb={1} key={testName}>
              <PolicyFormTestResult testName={testName} status={ComplianceStatusEnum.Pass} />
            </Box>
          ))}
          {results.testsFailed.map(testName => (
            <Box mb={1} key={testName}>
              <PolicyFormTestResult testName={testName} status={ComplianceStatusEnum.Fail} />
            </Box>
          ))}
          {results.testsErrored.map(({ name: testName, errorMessage }) => (
            <Box key={testName} mb={1}>
              <PolicyFormTestResult testName={testName} status={ComplianceStatusEnum.Error} />
              <Label size="small" is="pre" color={mapTestStatusToColor[ComplianceStatusEnum.Error]}>
                {errorMessage}
              </Label>
            </Box>
          ))}
        </React.Fragment>
      )}
    </Box>
  );
};

export default React.memo(RuleFormTestResultList);
