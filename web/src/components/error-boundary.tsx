import React from 'react';
import { logError } from 'Helpers/loggers';
import { Card, Text, Flex } from 'pouncejs';

interface ErrorBoundaryProps {
  fallbackStrategy?: 'passthrough' | 'invisible' | 'placeholder';
}

interface ErrorBoundaryState {
  hasError: boolean;
}

class ErrorBoundary extends React.Component<ErrorBoundaryProps, ErrorBoundaryState> {
  static defaultProps: Partial<ErrorBoundaryProps>;

  state = { hasError: false };

  componentDidCatch(error: Error, info: React.ErrorInfo) {
    // Display fallback UI
    this.setState({ hasError: true });

    // send this error to our logging service.
    logError(error, { extras: info });
  }

  render() {
    const { fallbackStrategy, children } = this.props;

    // if no error occurs -> render as normal
    if (!this.state.hasError) {
      return children;
    }

    // else decide what to show based on the selected strategy
    switch (fallbackStrategy) {
      case 'invisible':
        return null;
      case 'placeholder':
        return (
          <Card bg="red100" width="100%" height="100%">
            <Flex alignItems="center" justifyContent="center" width="100%" height="100%">
              <Text size="large" color="red300" py={5} px={4}>
                Something went wrong and we could not correctly display this content
              </Text>
            </Flex>
          </Card>
        );
      case 'passthrough':
      default:
        return children;
    }
  }
}

ErrorBoundary.defaultProps = {
  fallbackStrategy: 'placeholder',
};

export default ErrorBoundary;
