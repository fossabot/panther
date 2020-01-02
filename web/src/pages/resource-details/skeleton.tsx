import React from 'react';
import TablePlaceholder from 'Components/table-placeholder';
import { Box, Card } from 'pouncejs';
import Panel from 'Components/panel';
import ContentLoader from 'react-content-loader';

const ResourceDetailsPageSkeleton: React.FC = () => {
  return (
    <React.Fragment>
      <Box mb={2}>
        <Card p={6} mb={2}>
          <ContentLoader height={30}>
            <rect x="0" y={0} rx="1" ry="1" width="100%" height="10" />
            <rect x="0" y={15} rx="1" ry="1" width="100%" height="10" />
          </ContentLoader>
        </Card>
        <Card p={6} mb={2}>
          <ContentLoader height={30}>
            <rect x="0" y={0} rx="1" ry="1" width="100%" height="10" />
            <rect x="0" y={15} rx="1" ry="1" width="100%" height="10" />
          </ContentLoader>
        </Card>
      </Box>
      <Panel size="large" title="Policies">
        <Box mt={6}>
          <TablePlaceholder />
        </Box>
      </Panel>
    </React.Fragment>
  );
};

export default ResourceDetailsPageSkeleton;
