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
import { Button, BoxProps, Spinner, Flex } from 'pouncejs';

const SubmitButton: React.FC<BoxProps<HTMLButtonElement> & { submitting: boolean }> = ({
  submitting,
  disabled,
  children,
  ...rest
}) => (
  <Button {...rest} type="submit" size="large" variant="primary" disabled={disabled}>
    {submitting ? (
      <Flex alignItems="center" justifyContent="center">
        <Spinner size="small" mr={2} />
        {children}
      </Flex>
    ) : (
      children
    )}
  </Button>
);

export default React.memo(SubmitButton);
