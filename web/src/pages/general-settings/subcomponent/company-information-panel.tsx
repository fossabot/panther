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
