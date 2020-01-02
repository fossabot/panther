import React from 'react';
import { Box, Card, Flex, Icon, IconProps, Label } from 'pouncejs';
import ErrorBoundary from 'Components/error-boundary';

interface DonutChartWrapperProps {
  title: string;
  icon: IconProps['type'];
}

const DonutChartWrapper: React.FC<DonutChartWrapperProps> = ({ children, title, icon }) => (
  <Card p={6} height={340}>
    <Flex alignItems="center" is="header" mb={6} color="grey500">
      <Icon size="small" type={icon} mr={4} />
      <Label size="large" is="h4">
        {title}
      </Label>
    </Flex>
    <Box height={250}>
      <ErrorBoundary>{children}</ErrorBoundary>
    </Box>
  </Card>
);

export default DonutChartWrapper;
