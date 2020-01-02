import React from 'react';
import ContentLoader from 'react-content-loader';
import { Box, Card } from 'pouncejs';

interface PanelPlaceholderProps {
  /** The number of rows that the placeholder component should render. Defaults to 5 */
  rowCount?: number;

  /** The height of each row. Defaults to 15px */
  rowHeight?: number;
}

const PanelPlaceholder: React.FC<PanelPlaceholderProps> = ({ rowCount = 4, rowHeight = 15 }) => (
  <Card
    width={1}
    borderBottom="1px solid"
    borderColor="grey100"
    py={8}
    px={8}
    backgroundColor="#fff"
  >
    <Box pb={8} borderBottom="1px solid" borderColor="grey100">
      <ContentLoader height={10}>
        <rect x="0" y="0" rx="1" ry="1" width="30%" height="10" />
      </ContentLoader>
    </Box>
    <Box mt={8}>
      <ContentLoader height={rowCount * rowHeight}>
        {[...Array(rowCount)].map((__, index) => (
          <rect key={index} x="0" y={index * rowHeight} rx="1" ry="1" width="40%" height="10" />
        ))}
      </ContentLoader>
    </Box>
  </Card>
);

export default PanelPlaceholder;
