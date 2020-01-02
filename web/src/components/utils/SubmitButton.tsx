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
