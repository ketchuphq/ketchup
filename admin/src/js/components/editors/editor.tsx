import * as React from 'react';
import * as API from 'lib/api';
import QuillComponent from 'components/editors/html';
import TextEditorComponent from 'components/editors/text';

type EditorType = 'quill' | 'cm';

interface Props {
  content: API.Content;
}

interface State {
  editor: EditorType;
}

export default class EditorComponent extends React.Component<Props, State> {
  editorRef: React.RefObject<any>;

  constructor(props: Props) {
    super(props);
    this.state = {
      editor: 'quill',
    };
    this.editorRef = React.createRef();
  }

  insertImage(url: string, altText?: string) {
    const editor = this.editorRef.current;
    if (editor.insertImage) {
      editor.insertImage(url, altText);
    }
  }

  insertText(text: string, url: string) {
    const editor = this.editorRef.current;
    if (editor.insertText) {
      editor.insertText(text, url);
    }
  }

  render() {
    let data = {content: this.props.content, ref: this.editorRef};
    if (this.state.editor == 'quill') {
      return (
        <div className="editor">
          <QuillComponent {...data} />
        </div>
      );
    }
    return (
      <div className="editor">
        <div>
          <a className="txt-small" onClick={() => this.setState({editor: 'quill'})}>
            show editor
          </a>
        </div>
        <TextEditorComponent {...data} />>
      </div>
    );
  }
}
