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
import { PoliciesForResourceInput, ResourcesForPolicyInput } from 'Generated/schema';
import { Text, TextProps, defaultTheme } from 'pouncejs';

type Filters = PoliciesForResourceInput & ResourcesForPolicyInput;

interface TableComplianceFilterControlProps<FilterKey extends keyof Filters>
  extends Omit<TextProps, 'size'> {
  text: string;
  updateFilter: (filters: { [key: string]: Filters[FilterKey] }) => void;
  filterKey: FilterKey;
  filterValue: Filters[FilterKey];
  activeFilterValue?: Filters[FilterKey];
  count?: number;
  countColor?: keyof typeof defaultTheme.colors;
}

function TableComplianceFilterControl<FilterKey extends keyof Filters>({
  filterKey,
  filterValue,
  updateFilter,
  activeFilterValue,
  text,
  count,
  countColor,
  ...rest
}: TableComplianceFilterControlProps<FilterKey>): React.ReactElement {
  return (
    <Text
      {...rest}
      size="medium"
      p={2}
      color="grey300"
      is="button"
      borderRadius="medium"
      onClick={() => updateFilter({ [filterKey]: filterValue })}
      backgroundColor={filterValue === activeFilterValue ? 'grey50' : undefined}
    >
      {text}{' '}
      <Text size="medium" color={countColor} is="span">
        {count}
      </Text>
    </Text>
  );
}

export default TableComplianceFilterControl;
