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
