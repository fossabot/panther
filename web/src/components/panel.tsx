import React from 'react';
import { Box, Flex, Heading, Card, Label } from 'pouncejs';

interface PanelProps {
  title: string;
  size: 'small' | 'large';
  actions?: React.ReactNode;
}

const Panel: React.FC<PanelProps> = ({ title, actions, size, children }) => {
  return (
    <Card
      is="section"
      width={1}
      borderBottom="1px solid"
      borderColor="grey100"
      p={size === 'large' ? 8 : 6}
    >
      <Flex
        pb={size === 'large' ? 8 : 6}
        borderBottom="1px solid"
        borderColor="grey100"
        justifyContent="space-between"
        alignItems="center"
      >
        {size === 'large' ? (
          <Heading size="medium" is="h2">
            {title}
          </Heading>
        ) : (
          <Label size="large" is="h4">
            {title}
          </Label>
        )}
        {actions}
      </Flex>
      <Box mt={size === 'large' ? 8 : 6}>{children}</Box>
    </Card>
  );
};

export default Panel;
