import * as React from 'react';
import { css } from '@emotion/core';
import { Box, Card, Text } from 'pouncejs';

interface ItemCardProps {
  logo: string;
  title: string;
  onClick?: () => void;
}

const DestinationCard: React.FunctionComponent<ItemCardProps> = ({ logo, title, onClick }) => (
  <Card
    is="button"
    onClick={onClick}
    css={css`
      cursor: pointer;
      transition: transform 0.15s ease-in-out;
      &:hover {
        transform: scale3d(1.03, 1.03, 1.03);
      }
    `}
  >
    <Box height={92} px={10}>
      <img src={logo} alt={title} style={{ objectFit: 'contain' }} width="100%" height="100%" />
    </Box>
    <Box borderTopStyle="solid" borderTopWidth="1px" borderColor="grey50">
      <Text size="medium" px={4} py={3} color="grey500" textAlign="left">
        {title}
      </Text>
    </Box>
  </Card>
);

export default DestinationCard;
