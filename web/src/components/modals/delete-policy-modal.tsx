import React from 'react';
import { LIST_POLICIES } from 'Pages/list-policies';
import { DeletePolicyInput, PolicySummary, PolicyDetails } from 'Generated/schema';

import { useMutation, gql } from '@apollo/client';
import useRouter from 'Hooks/useRouter';
import urls from 'Source/urls';
import { getOperationName } from '@apollo/client/utilities/graphql/getFromAST';
import BaseDeleteModal from 'Components/modals/base-delete-modal';

const DELETE_POLICY = gql`
  mutation DeletePolicy($input: DeletePolicyInput!) {
    deletePolicy(input: $input)
  }
`;

export interface DeletePolicyModalProps {
  policy: PolicyDetails | PolicySummary;
}

const DeletePolicyModal: React.FC<DeletePolicyModalProps> = ({ policy }) => {
  const { location, history } = useRouter<{ id?: string }>();
  const policyDisplayName = policy.displayName || policy.id;
  const mutation = useMutation<boolean, { input: DeletePolicyInput }>(DELETE_POLICY, {
    awaitRefetchQueries: true,
    refetchQueries: [getOperationName(LIST_POLICIES)],
    variables: {
      input: {
        policies: [
          {
            id: policy.id,
          },
        ],
      },
    },
  });

  return (
    <BaseDeleteModal
      mutation={mutation}
      itemDisplayName={policyDisplayName}
      onSuccess={() => {
        if (location.pathname.includes(policy.id)) {
          // if we were on the particular policy's details page or edit page --> redirect on delete
          history.push(urls.policies.list());
        }
      }}
    />
  );
};

export default DeletePolicyModal;
