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
import { Box, Grid, Label, Text } from 'pouncejs';
import Panel from 'Components/panel';
import { capitalize, formatDatetime } from 'Helpers/utils';
import { ComplianceStatusEnum, Integration, ResourceDetails } from 'Generated/schema';

interface ResourceDetailsInfoProps {
  resource?: ResourceDetails & Pick<Integration, 'integrationLabel'>;
}

const ResourceDetailsInfo: React.FC<ResourceDetailsInfoProps> = ({ resource }) => {
  return (
    <Panel size="large" title="Resource Details">
      <Grid gridTemplateColumns="repeat(3, 1fr)" gridGap={6}>
        <Box my={1}>
          <Label mb={1} is="div" size="small" color="grey300">
            ID
          </Label>
          <Text size="medium" color="black">
            {resource.id}
          </Text>
        </Box>
        <Box my={1}>
          <Label mb={1} is="div" size="small" color="grey300">
            TYPE
          </Label>
          <Text size="medium" color="black">
            {resource.type}
          </Text>
        </Box>
        <Box my={1}>
          <Label mb={1} is="div" size="small" color="grey300">
            SOURCE
          </Label>
          <Text size="medium" color="black">
            {resource.integrationLabel}
          </Text>
        </Box>
        <Box my={1}>
          <Label mb={1} is="div" size="small" color="grey300">
            STATUS
          </Label>
          <Text
            size="medium"
            color={resource.complianceStatus === ComplianceStatusEnum.Pass ? 'green300' : 'red300'}
          >
            {capitalize(resource.complianceStatus.toLowerCase())}
          </Text>
        </Box>
        <Box my={1}>
          <Label mb={1} is="div" size="small" color="grey300">
            LAST MODIFIED
          </Label>
          <Text size="medium" color="black">
            {formatDatetime(resource.lastModified)}
          </Text>
        </Box>
      </Grid>
    </Panel>
  );
};

export default React.memo(ResourceDetailsInfo);
