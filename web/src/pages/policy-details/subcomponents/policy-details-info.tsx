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
import { Link } from 'react-router-dom';
import { Badge, Box, Button, Grid, Icon, Label, Text } from 'pouncejs';
import { capitalize, formatDatetime } from 'Helpers/utils';
import Panel from 'Components/panel';
import { ComplianceStatusEnum, PolicyDetails } from 'Generated/schema';
import { READONLY_ROLES_ARRAY, SEVERITY_COLOR_MAP } from 'Source/constants';
import { LinkifyProps } from 'linkifyjs/react';
import urls from 'Source/urls';
import JsonViewer from 'Components/json-viewer';
import useModal from 'Hooks/useModal';
import { MODALS } from 'Components/utils/modal-context';
import RoleRestrictedAccess from 'Components/role-restricted-access';

const Linkify = React.lazy(() =>
  import(/* webpackChunkName: "linkify" */ 'linkifyjs/react.js')
) as React.FC<LinkifyProps>;

interface ResourceDetailsInfoProps {
  policy?: PolicyDetails;
}

const PolicyDetailsInfo: React.FC<ResourceDetailsInfoProps> = ({ policy }) => {
  const { showModal } = useModal();

  return (
    <Panel
      size="large"
      title="Policy Details"
      actions={
        <RoleRestrictedAccess deniedRoles={READONLY_ROLES_ARRAY}>
          <Box>
            <Button
              size="large"
              variant="default"
              mr={4}
              is={Link}
              to={urls.policies.edit(policy.id)}
            >
              Edit
            </Button>
            <Button
              size="large"
              variant="default"
              color="red300"
              onClick={() =>
                showModal({
                  modal: MODALS.DELETE_POLICY,
                  props: { policy },
                })
              }
            >
              Delete
            </Button>
          </Box>
        </RoleRestrictedAccess>
      }
    >
      <Grid gridTemplateColumns="repeat(3, 1fr)" gridGap={6}>
        <Box my={1}>
          <Label mb={1} is="div" size="small" color="grey300">
            ID
          </Label>
          <Text size="medium" color="black">
            {policy.id}
          </Text>
        </Box>
        <Box my={1}>
          <Label mb={1} is="div" size="small" color="grey300">
            DISPLAY NAME
          </Label>
          <Text size="medium" color={policy.displayName ? 'black' : 'grey200'}>
            {policy.displayName || 'No display name found'}
          </Text>
        </Box>
        <Box my={1}>
          <Label mb={1} is="div" size="small" color="grey300">
            ENABLED
          </Label>
          {policy.enabled ? (
            <Icon type="check" color="green300" size="large" />
          ) : (
            <Icon type="close" color="red300" size="large" />
          )}
        </Box>
        <Box my={1}>
          <Label mb={1} is="div" size="small" color="grey300">
            STATUS
          </Label>
          <Text
            size="medium"
            color={policy.complianceStatus === ComplianceStatusEnum.Pass ? 'green300' : 'red300'}
          >
            {capitalize(policy.complianceStatus.toLowerCase())}
          </Text>
        </Box>
        <Box my={1}>
          <Label mb={1} is="div" size="small" color="grey300">
            IGNORE PATTERNS
          </Label>
          {policy.suppressions.length ? (
            policy.suppressions.map(suppression => (
              <Text size="medium" color="black" key={suppression}>
                {suppression}
              </Text>
            ))
          ) : (
            <Text size="medium" color="grey200">
              No resource is being ignored
            </Text>
          )}
        </Box>
        <Box my={1}>
          <Label mb={1} is="div" size="small" color="grey300">
            REFERENCE
          </Label>
          <Text size="medium" color={policy.reference ? 'blue300' : 'grey200'}>
            {policy.reference ? (
              <a href={policy.reference} target="_blank" rel="noopener noreferrer">
                {policy.reference}
              </a>
            ) : (
              'No reference found'
            )}
          </Text>
        </Box>
        <Box my={1}>
          <Label mb={1} is="div" size="small" color="grey300">
            RESOURCE TYPES
          </Label>
          {policy.resourceTypes.length ? (
            policy.resourceTypes.map(resourceType => (
              <Text size="medium" color="black" key={resourceType}>
                {resourceType}
              </Text>
            ))
          ) : (
            <Text size="medium" color="black">
              All resources
            </Text>
          )}
        </Box>
        <Box my={1}>
          <Label mb={1} is="div" size="small" color="grey300">
            DESCRIPTION
          </Label>
          <Text size="medium" color={policy.description ? 'black' : 'grey200'}>
            <React.Suspense fallback={<span>{policy.description}</span>}>
              <Linkify>{policy.description || 'No description available'}</Linkify>
            </React.Suspense>
          </Text>
        </Box>
        <Box my={1}>
          <Label mb={1} is="div" size="small" color="grey300">
            RUNBOOK
          </Label>
          <Text size="medium" color={policy.runbook ? 'black' : 'grey200'}>
            <React.Suspense fallback={<span>{policy.runbook}</span>}>
              <Linkify>{policy.runbook || 'No runbook available'}</Linkify>
            </React.Suspense>
          </Text>
        </Box>
        <Box my={1}>
          <Label mb={1} is="div" size="small" color="grey300">
            SEVERITY
          </Label>
          <Badge color={SEVERITY_COLOR_MAP[policy.severity]}>{policy.severity}</Badge>
        </Box>
        <Box my={1}>
          <Label mb={1} is="div" size="small" color="grey300">
            TAGS
          </Label>
          {policy.tags.length ? (
            policy.tags.map((tag, index) => (
              <Text size="medium" color="black" key={tag} is="span">
                {tag}
                {index !== policy.tags.length - 1 ? ', ' : null}
              </Text>
            ))
          ) : (
            <Text size="medium" color="grey200">
              No tags assigned
            </Text>
          )}
        </Box>
        <Box my={1}>
          <Label mb={1} is="div" size="small" color="grey300">
            AUTO-REMEDIATION
          </Label>
          <Text size="medium" color={policy.autoRemediationId ? 'black' : 'grey200'}>
            {policy.autoRemediationId || 'Not configured'}
          </Text>
        </Box>
        <Box my={1}>
          <Label mb={1} is="div" size="small" color="grey300">
            CREATED
          </Label>
          <Text size="medium" color="black">
            {formatDatetime(policy.createdAt)}
          </Text>
        </Box>
        <Box my={1}>
          <Label mb={1} is="div" size="small" color="grey300">
            LAST MODIFIED
          </Label>
          <Text size="medium" color="black">
            {formatDatetime(policy.lastModified)}
          </Text>
        </Box>
        {policy.autoRemediationId && (
          <Box my={1}>
            <Label mb={1} is="div" size="small" color="grey300">
              REMEDIATION PARAMETERS
            </Label>
            <JsonViewer data={JSON.parse(policy.autoRemediationParameters)} />
          </Box>
        )}
      </Grid>
    </Panel>
  );
};

export default React.memo(PolicyDetailsInfo);
