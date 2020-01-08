import { Alert, Flex, Heading, Icon, IconButton, SideSheet, useSnackbar, Box } from 'pouncejs';
import React from 'react';
import { useMutation, gql } from '@apollo/client';

import useSidesheet from 'Hooks/useSidesheet';
import { SIDESHEETS } from 'Components/utils/sidesheet-context';
import {
  Destination,
  DestinationConfigInput,
  DestinationInput,
  DestinationTypeEnum,
} from 'Generated/schema';
import { BaseDestinationFormValues } from 'Components/forms/common/base-destination-form';
import {
  PagerDutyDestinationForm,
  EmailDestinationForm,
  JiraDestinationForm,
  OpsgenieDestinationForm,
  MicrosoftTeamsDestinationForm,
  SlackDestinationForm,
  GithubDestinationForm,
  SNSDestinationForm,
  SQSDestinationForm,
} from 'Components/forms/destination-form';
import { LIST_DESTINATIONS, ListDestinationsQueryData } from 'Pages/destinations';
import { capitalize, extractErrorMessage } from 'Helpers/utils';

const ADD_DESTINATION = gql`
  mutation AddSlackDestination($input: DestinationInput!) {
    addDestination(input: $input) {
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
  addDestination: Destination;
}

interface DestinationMutationInput {
  input: DestinationInput;
}

export interface AddDestinationSidesheetProps {
  destinationType: DestinationTypeEnum;
}

const AddDestinationSidesheet: React.FC<AddDestinationSidesheetProps> = ({ destinationType }) => {
  const { pushSnackbar } = useSnackbar();
  const { hideSidesheet, showSidesheet } = useSidesheet();

  // If destination object doesn't exist, handleSubmit should call addDestination to create a new destination and use default initial values
  const [addDestination, { data: addDestinationData, error: addDestinationError }] = useMutation<
    DestinationMutationData,
    DestinationMutationInput
  >(ADD_DESTINATION);

  React.useEffect(() => {
    if (addDestinationData) {
      pushSnackbar({
        variant: 'success',
        title: `Successfully added ${addDestinationData.addDestination.displayName}`,
      });
      hideSidesheet();
    }
  }, [addDestinationData]);

  // The typescript on `values` simply says that we expect to have DestinationFormValues with an
  // `outputType` that partially implements the DestinationConfigInput (we say partially since each
  // integration will add each own config). Ideally we would want to say "exactly 1". We can't
  // specify the exact one since `const` are not allowed to have generics and `useCallback` can only
  // be assigned to a const
  const handleSubmit = React.useCallback(
    async (values: BaseDestinationFormValues<Partial<DestinationConfigInput>>) => {
      const { displayName, defaultForSeverity, outputConfig } = values;
      await addDestination({
        variables: {
          input: {
            // form values that are present in all destinations
            displayName,
            defaultForSeverity,

            // dynamic form values that depend on the selected destination
            outputType: destinationType,
            outputConfig,
          },
        },
        update: (proxy, { data }: { data: DestinationMutationData }) => {
          // Read the data from our cache for this query.
          const existingData: ListDestinationsQueryData = proxy.readQuery({
            query: LIST_DESTINATIONS,
          });

          // Write our data back to the cache with the new comment in it
          proxy.writeQuery({
            query: LIST_DESTINATIONS,
            data: {
              ...existingData,
              destinations: [data.addDestination, ...existingData.destinations],
            },
          });
        },
      });
    },
    []
  );

  const commonInitialValues = {
    displayName: '',
    defaultForSeverity: [],
  };

  const renderFullDestinationForm = () => {
    switch (destinationType) {
      case DestinationTypeEnum.Email:
        return (
          <EmailDestinationForm
            initialValues={{
              ...commonInitialValues,
              outputConfig: { email: { destinationAddress: '' } },
            }}
            onSubmit={handleSubmit}
          />
        );
      case DestinationTypeEnum.Pagerduty:
        return (
          <PagerDutyDestinationForm
            initialValues={{
              ...commonInitialValues,
              outputConfig: { pagerDuty: { integrationKey: '' } },
            }}
            onSubmit={handleSubmit}
          />
        );
      case DestinationTypeEnum.Github:
        return (
          <GithubDestinationForm
            initialValues={{
              ...commonInitialValues,
              outputConfig: { github: { repoName: '', token: '' } },
            }}
            onSubmit={handleSubmit}
          />
        );
      case DestinationTypeEnum.Jira:
        return (
          <JiraDestinationForm
            initialValues={{
              ...commonInitialValues,
              outputConfig: {
                jira: {
                  orgDomain: '',
                  projectKey: '',
                  userName: '',
                  apiKey: '',
                  assigneeID: '',
                },
              },
            }}
            onSubmit={handleSubmit}
          />
        );
      case DestinationTypeEnum.Opsgenie:
        return (
          <OpsgenieDestinationForm
            initialValues={{ ...commonInitialValues, outputConfig: { opsgenie: { apiKey: '' } } }}
            onSubmit={handleSubmit}
          />
        );
      case DestinationTypeEnum.Msteams:
        return (
          <MicrosoftTeamsDestinationForm
            initialValues={{ ...commonInitialValues, outputConfig: { msTeams: { webhookURL: '' } } }} // prettier-ignore
            onSubmit={handleSubmit}
          />
        );
      case DestinationTypeEnum.Slack:
        return (
          <SlackDestinationForm
            initialValues={{ ...commonInitialValues, outputConfig: { slack: { webhookURL: '' } } }}
            onSubmit={handleSubmit}
          />
        );
      case DestinationTypeEnum.Sns:
        return (
          <SNSDestinationForm
            initialValues={{
              ...commonInitialValues,
              outputConfig: { sns: { topicArn: '' } },
            }}
            onSubmit={handleSubmit}
          />
        );
      case DestinationTypeEnum.Sqs:
        return (
          <SQSDestinationForm
            initialValues={{
              ...commonInitialValues,
              outputConfig: { sqs: { queueUrl: '' } },
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
        <Flex mb={8} alignItems="center" mt={-2}>
          <IconButton
            mr={4}
            variant="default"
            onClick={() =>
              showSidesheet({
                sidesheet: SIDESHEETS.SELECT_DESTINATION,
              })
            }
          >
            <Icon size="large" type="arrow-back" />
          </IconButton>
          <Heading size="medium">{capitalize(destinationType)} Configuration</Heading>
        </Flex>
        {addDestinationError && (
          <Alert
            mt={2}
            mb={6}
            variant="error"
            title="Destination not added"
            description={
              extractErrorMessage(addDestinationError) ||
              "An unknown error occured and we couldn't add your new destination"
            }
          />
        )}
        {renderFullDestinationForm()}
      </Box>
    </SideSheet>
  );
};

export default AddDestinationSidesheet;
