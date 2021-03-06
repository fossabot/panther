/**
 * Panther is a scalable, powerful, cloud-native SIEM written in Golang/React.
 * Copyright (C) 2020 Panther Labs Inc
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

import { Alert, Heading, SideSheet, useSnackbar, Box } from 'pouncejs';
import React from 'react';
import { useMutation, gql } from '@apollo/client';

import pick from 'lodash-es/pick';
import useSidesheet from 'Hooks/useSidesheet';
import {
  Destination,
  DestinationConfigInput,
  DestinationInput,
  DestinationTypeEnum,
} from 'Generated/schema';

import { BaseDestinationFormValues } from 'Components/forms/common/base-destination-form';
import {
  EmailDestinationForm,
  GithubDestinationForm,
  JiraDestinationForm,
  MicrosoftTeamsDestinationForm,
  OpsgenieDestinationForm,
  PagerDutyDestinationForm,
  SlackDestinationForm,
  SNSDestinationForm,
  SQSDestinationForm,
} from 'Components/forms/destination-form';
import { extractErrorMessage } from 'Helpers/utils';

const UPDATE_DESTINATION = gql`
  mutation UpdateSlackDestination($input: DestinationInput!) {
    updateDestination(input: $input) {
      createdBy
      creationTime
      displayName
      lastModifiedBy
      lastModifiedTime
      outputId
      outputType
      outputConfig {
        slack {
          webhookURL
        }
        sns {
          topicArn
        }
        email {
          destinationAddress
        }
        pagerDuty {
          integrationKey
        }
        github {
          repoName
          token
        }
        jira {
          orgDomain
          projectKey
          userName
          apiKey
          assigneeID
        }
        opsgenie {
          apiKey
        }
        msTeams {
          webhookURL
        }
      }
      verificationStatus
      defaultForSeverity
    }
  }
`;

interface DestinationMutationData {
  updateDestination: Destination;
}

interface DestinationMutationInput {
  input: DestinationInput;
}

// Normally the `destination` doesn't contain the severities, but because we receive it as a prop
// from the destinations table, we are able to access a `defaultForSeverities` key that the table
// has assigned for us. Thus the `destination` that we actually received in enhanced with this
// property.
export interface UpdateDestinationSidesheetProps {
  destination: Destination;
}

export const UpdateDestinationSidesheet: React.FC<UpdateDestinationSidesheetProps> = ({
  destination,
}) => {
  const { pushSnackbar } = useSnackbar();
  const { hideSidesheet } = useSidesheet();

  // If destination object exist, handleSubmit should call updateDestination and use attributes from the destination object for form initial values
  const [
    updateDestination,
    { data: updateDestinationData, error: updateDestinationError },
  ] = useMutation<DestinationMutationData, DestinationMutationInput>(UPDATE_DESTINATION);

  React.useEffect(() => {
    if (updateDestinationData) {
      pushSnackbar({
        variant: 'success',
        title: `Successfully updated ${updateDestinationData.updateDestination.displayName}`,
      });
      hideSidesheet();
    }
  }, [updateDestinationData]);

  const handleSubmit = React.useCallback(
    async (values: BaseDestinationFormValues<Partial<DestinationConfigInput>>) => {
      const { displayName, defaultForSeverity, outputConfig } = values;

      await updateDestination({
        variables: {
          input: {
            // static form values that are present on all destinations
            displayName,
            defaultForSeverity,

            // needed fields from the server in order to update the selected destination
            outputId: destination.outputId,
            outputType: destination.outputType,

            // dynamic form values that depend on the selected destination
            outputConfig,
          },
        },
      });
    },
    []
  );

  // Normally we would want to perform a single `pick` operation per switch-case and not extend the
  // commonInitialValues that are defined here. Unfortunately, if you do deep picking (i.e. x.w.y)
  // on  lodash's pick, it messes typings and TS fails cause it thinks it doesn't have all the
  // fields it needs. We use `commonInitialValues` to satisfy this exact constraint that was set by
  // the `initialValues` prop of each form.
  const commonInitialValues = pick(destination, ['displayName', 'defaultForSeverity']);

  const renderFullDestinationForm = () => {
    switch (destination.outputType) {
      case DestinationTypeEnum.Email:
        return (
          <EmailDestinationForm
            initialValues={{
              ...commonInitialValues,
              outputConfig: pick(destination.outputConfig, 'email.destinationAddress'),
            }}
            onSubmit={handleSubmit}
          />
        );
      case DestinationTypeEnum.Pagerduty:
        return (
          <PagerDutyDestinationForm
            initialValues={{
              ...commonInitialValues,
              outputConfig: pick(destination.outputConfig, 'pagerDuty.integrationKey'),
            }}
            onSubmit={handleSubmit}
          />
        );
      case DestinationTypeEnum.Github:
        return (
          <GithubDestinationForm
            initialValues={{
              ...commonInitialValues,
              outputConfig: pick(destination.outputConfig, ['github.repoName', 'github.apiKey']),
            }}
            onSubmit={handleSubmit}
          />
        );
      case DestinationTypeEnum.Jira:
        return (
          <JiraDestinationForm
            initialValues={{
              ...commonInitialValues,
              outputConfig: pick(destination.outputConfig, [
                'jira.orgDomain',
                'jira.projectKey',
                'jira.userName',
                'jira.apiKey',
                'jira.assigneeID',
              ]),
            }}
            onSubmit={handleSubmit}
          />
        );
      case DestinationTypeEnum.Opsgenie:
        return (
          <OpsgenieDestinationForm
            initialValues={{
              ...commonInitialValues,
              outputConfig: pick(destination.outputConfig, 'opsgenie.apiKey'),
            }}
            onSubmit={handleSubmit}
          />
        );
      case DestinationTypeEnum.Msteams:
        return (
          <MicrosoftTeamsDestinationForm
            initialValues={{
              ...commonInitialValues,
              outputConfig: pick(destination.outputConfig, 'msTeams.webhookURL'),
            }}
            onSubmit={handleSubmit}
          />
        );
      case DestinationTypeEnum.Sns:
        return (
          <SNSDestinationForm
            initialValues={{
              ...commonInitialValues,
              outputConfig: pick(destination.outputConfig, 'sns.topicArn'),
            }}
            onSubmit={handleSubmit}
          />
        );
      case DestinationTypeEnum.Sqs:
        return (
          <SQSDestinationForm
            initialValues={{
              ...commonInitialValues,
              outputConfig: pick(destination.outputConfig, 'sqs.queueUrl'),
            }}
            onSubmit={handleSubmit}
          />
        );
      case DestinationTypeEnum.Slack:
        return (
          <SlackDestinationForm
            initialValues={{
              ...commonInitialValues,
              outputConfig: pick(destination.outputConfig, 'slack.webhookURL'),
            }}
            onSubmit={handleSubmit}
          />
        );
      default:
        return null;
    }
  };

  return (
    <SideSheet open onClose={hideSidesheet}>
      <Box width={465}>
        <Heading size="medium" mb={8}>
          Update {destination.outputType}
        </Heading>
        {updateDestinationError && (
          <Alert
            mt={2}
            mb={6}
            variant="error"
            title="Destination not updated"
            description={
              extractErrorMessage(updateDestinationError) ||
              'An unknown error has occured while trying to update your destination'
            }
          />
        )}
        {renderFullDestinationForm()}
      </Box>
    </SideSheet>
  );
};

export default UpdateDestinationSidesheet;
