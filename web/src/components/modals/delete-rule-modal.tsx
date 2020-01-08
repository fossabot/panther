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
import { DeletePolicyInput, RuleSummary, RuleDetails } from 'Generated/schema';

import { useMutation, gql } from '@apollo/client';
import useRouter from 'Hooks/useRouter';
import urls from 'Source/urls';
import { getOperationName } from '@apollo/client/utilities/graphql/getFromAST';
import { LIST_RULES } from 'Pages/list-rules';
import BaseDeleteModal from 'Components/modals/base-delete-modal';

// Delete Rule and Delete Policy uses the same endpoint
const DELETE_RULE = gql`
  mutation DeletePolicy($input: DeletePolicyInput!) {
    deletePolicy(input: $input)
  }
`;

export interface DeleteRuleModalProps {
  rule: RuleDetails | RuleSummary;
}

const DeleteRuleModal: React.FC<DeleteRuleModalProps> = ({ rule }) => {
  const { location, history } = useRouter<{ id?: string }>();
  const ruleDisplayName = rule.displayName || rule.id;
  const mutation = useMutation<boolean, { input: DeletePolicyInput }>(DELETE_RULE, {
    awaitRefetchQueries: true,
    refetchQueries: [getOperationName(LIST_RULES)],
    variables: {
      input: {
        policies: [
          {
            id: rule.id,
          },
        ],
      },
    },
  });

  return (
    <BaseDeleteModal
      mutation={mutation}
      itemDisplayName={ruleDisplayName}
      onSuccess={() => {
        if (location.pathname.includes(rule.id)) {
          // if we were on the particular rule's details page or edit page --> redirect on delete
          history.push(urls.rules.list());
        }
      }}
    />
  );
};

export default DeleteRuleModal;
