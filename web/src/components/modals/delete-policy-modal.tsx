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
