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

import React, { useState } from 'react';
import Panel from 'Components/panel';
import { Flex, Text } from 'pouncejs';
import CompanyInformationForm from 'Pages/general-settings/subcomponent/company-information-form';

interface CompanyInformationProps {
  displayName: string;
  email: string;
}

const CompanyInformation: React.FC<CompanyInformationProps> = ({ displayName, email }) => {
  const [isEditing, setEditingState] = useState<boolean>(false);

  return (
    <Panel size="large" title={'Company Information'}>
      {isEditing ? (
        <CompanyInformationForm
          onSuccess={() => setEditingState(false)}
          displayName={displayName}
          email={email}
        />
      ) : (
        <CompanyInformationReadOnly displayName={displayName} email={email} />
      )}
    </Panel>
  );
};

const CompanyInformationReadOnly: React.FC<CompanyInformationProps> = ({ displayName, email }) => (
  <React.Fragment>
    <Flex mb={6}>
      <Text size="medium" minWidth={150} color="grey400" fontWeight="bold">
        NAME
      </Text>
      <Text size="medium" color="grey400">
        {displayName || '-'}
      </Text>
    </Flex>
    <Flex>
      <Text size="medium" minWidth={150} color="grey400" fontWeight="bold">
        EMAIL
      </Text>
      <Text size="medium" color="grey400">
        {email || '-'}
      </Text>
    </Flex>
  </React.Fragment>
);

export default CompanyInformation;
