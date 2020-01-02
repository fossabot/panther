import React from 'react';
import { useTheme } from 'pouncejs';

const ReactJSONView = React.lazy(() =>
  import(/* webpackChunkName: 'react-json-view' */ 'react-json-view')
);

interface JsonViewerProps {
  data: { [key: string]: string | object | number };
  collapsed?: boolean;
}

const JsonViewer: React.FC<JsonViewerProps> = ({ data, collapsed }) => {
  const theme = useTheme();
  return (
    <React.Suspense fallback={null}>
      <ReactJSONView
        src={data}
        name={false}
        theme="grayscale:inverted"
        iconStyle="triangle"
        displayObjectSize={false}
        displayDataTypes={false}
        collapsed={collapsed || 1}
        style={{ fontSize: theme.fontSizes[2] }}
        sortKeys
      />
    </React.Suspense>
  );
};

export default JsonViewer;
