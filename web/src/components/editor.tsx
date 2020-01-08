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
import { IAceEditorProps } from 'react-ace/lib/ace';

// Lazy-load the ace editor. Make sure that both editor and modes get bundled under the same chunk
const AceEditor = React.lazy(() => import(/* webpackChunkName: "ace-editor" */ 'react-ace'));

const baseAceEditorConfig = {
  fontSize: '16px',
  editorProps: {
    $blockScrolling: Infinity,
  },
  wrapEnabled: true,
  theme: 'cobalt',
  showPrintMargin: true,
  showGutter: true,
  highlightActiveLine: true,
  maxLines: Infinity,
  style: {
    zIndex: 0,
  },
};

export type EditorProps = IAceEditorProps;

const Editor: React.FC<EditorProps> = props => {
  // Asynchronously load (post-mount) all the mode & themes
  React.useEffect(() => {
    import(/* webpackChunkName: "ace-editor" */ 'brace/mode/json');
    import(/* webpackChunkName: "ace-editor" */ 'brace/mode/python');
    import(/* webpackChunkName: "ace-editor" */ 'brace/theme/cobalt');
  }, []);

  return (
    <React.Suspense fallback={null}>
      <AceEditor {...baseAceEditorConfig} {...props} />
    </React.Suspense>
  );
};

export default React.memo(Editor);
