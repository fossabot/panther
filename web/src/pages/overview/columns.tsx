/* eslint-disable react/display-name */
import React from 'react';
import { generateEnumerationColumn } from 'Helpers/utils';
import { Badge, TableProps } from 'pouncejs';
import { SEVERITY_COLOR_MAP } from 'Source/constants';
import { TopFailingPolicy, TopFailingResource } from 'Pages/overview/index';

/**
 * The columns that the top failing policies table will show
 */
export const topFailingPoliciesColumns = [
  // add an enumeration column starting from 0
  generateEnumerationColumn(0),

  // The name is the `id` of the policy
  {
    key: 'id',
    header: 'Policy',
    flex: '1 0 0',
  },

  // Render badges to showcase severity
  {
    key: 'severity',
    flex: '0 0 100px',
    header: 'Severity',
    renderCell: ({ severity }) => <Badge color={SEVERITY_COLOR_MAP[severity]}>{severity}</Badge>,
  },
] as TableProps<TopFailingPolicy>['columns'];

/**
 * The columns that the top failing resources table will show
 */
export const topFailingResourcesColumns = [
  // add an enumeration column starting from 0
  generateEnumerationColumn(0),

  // The name is the `id` of the policy
  {
    key: 'id',
    header: 'Resource',
    flex: '1 0 0',
  },
] as TableProps<TopFailingResource>['columns'];
