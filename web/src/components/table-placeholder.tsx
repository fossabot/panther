import React from 'react';
import ContentLoader from 'react-content-loader';

interface TablePlaceholderProps {
  /** The number of rows that the placeholder component should render. Defaults to 5 */
  rowCount?: number;

  /** The height of each row. Defaults to 10px */
  rowHeight?: number;

  /** The vertical gap between each row. Defaults to 5px */
  rowGap?: number;
}

const TablePlaceholder: React.FC<TablePlaceholderProps> = ({
  rowCount = 5,
  rowHeight = 10,
  rowGap = 5,
}) => (
  <ContentLoader height={rowCount * (rowHeight + rowGap)}>
    {[...Array(rowCount)].map((__, index) => (
      <rect
        key={index}
        x="0"
        y={index * (rowHeight + rowGap)}
        rx="1"
        ry="1"
        width="100%"
        height={rowHeight}
      />
    ))}
  </ContentLoader>
);

export default TablePlaceholder;
