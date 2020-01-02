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
