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
