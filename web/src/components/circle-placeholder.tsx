import React from 'react';
import ContentLoader from 'react-content-loader';

interface CirclePlaceholderProps {
  /** The radius of the circle. Defaults to 50px (100px total width and height) */
  size?: number;
}

const CirclePlaceholder: React.FC<CirclePlaceholderProps> = ({ size = 50 }) => (
  <ContentLoader height={2 * size} style={{ width: '100%' }}>
    <circle cx="50%" cy="50%" r={size} />
  </ContentLoader>
);

export default React.memo(CirclePlaceholder);
