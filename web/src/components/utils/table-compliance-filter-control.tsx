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
