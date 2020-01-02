import React from 'react';
import TablePlaceholder from 'Components/table-placeholder';
import { Box, Card } from 'pouncejs';
import Panel from 'Components/panel';

const RuleDetailsPageSkeleton: React.FC = () => {
  return (
    <React.Fragment>
      <Card p={6}>
        <TablePlaceholder rowCount={2} rowHeight={10} />
      </Card>
      <Box mt={2} mb={6}>
        <Panel size="large" title="Alerts">
          <Box mt={6}>
            <TablePlaceholder />
          </Box>
        </Panel>
      </Box>
    </React.Fragment>
  );
};

export default RuleDetailsPageSkeleton;
