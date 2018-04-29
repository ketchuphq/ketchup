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
  constructor(props: Props) {
    super(props);
    this.state = {
      editor: 'quill',
    };
  }

  render() {
    let data = {content: this.props.content};
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
