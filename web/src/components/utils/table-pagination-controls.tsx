import React from 'react';
import { Flex, Icon, IconButton, Label } from 'pouncejs';

interface TablePaginationControls {
  page: number;
  onPageChange: (page: number) => void;
  totalPages: number;
}

const TablePaginationControls: React.FC<TablePaginationControls> = ({
  page,
  onPageChange,
  totalPages,
}) => {
  return (
    <Flex alignItems="center" justifyContent="center">
      <Flex mr={9} alignItems="center">
        <IconButton variant="default" disabled={page <= 1} onClick={() => onPageChange(page - 1)}>
          <Icon size="large" type="chevron-left" />
        </IconButton>
        <Label size="large" mx={4} color="grey400">
          {page} of {totalPages}
        </Label>
        <IconButton
          variant="default"
          disabled={page >= totalPages}
          onClick={() => onPageChange(page + 1)}
        >
          <Icon size="large" type="chevron-right" />
        </IconButton>
      </Flex>
    </Flex>
  );
};

export default React.memo(TablePaginationControls);
