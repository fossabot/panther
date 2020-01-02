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
