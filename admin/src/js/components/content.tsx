import * as React from 'react';
import * as API from 'lib/api';
import EditorComponent from 'components/editors/editor';
import QuillComponent from 'components/editors/html';
import TextEditorComponent from 'components/editors/text';
import CodeMirrorComponent from 'components/editors/markdown';

interface ContentEditor {
  shouldRender(content: API.Content): boolean;
  render(content: API.Content): React.ReactElement<any>;
}

let LongHTMLEditor: ContentEditor = {
  shouldRender: (content: API.Content) => content.text != null && content.text.type == 'html',
  render: (content: API.Content) => <EditorComponent content={content} />,
};

let LongMarkdownEditor: ContentEditor = {
  shouldRender: (content: API.Content) => content.text != null && content.text.type == 'markdown',
  render: (content: API.Content) => <CodeMirrorComponent content={content} />,
};

let LongTextEditor: ContentEditor = {
  shouldRender: (content: API.Content) => content.text != null && content.text.type == 'text',
  render: (content: API.Content) => <TextEditorComponent content={content} />,
};

// exported for test
export class ShortTextEditorComponent extends React.PureComponent<{content: API.Content}> {
  textInput: React.RefObject<HTMLInputElement>;
  constructor(props: any) {
    super(props);
    this.textInput = React.createRef();
  }
  componentDidMount() {
    this.textInput.current.value = this.props.content.value;
  }

  render() {
    return (
      <input
        type="text"
        ref={this.textInput}
        onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
          this.props.content.value = e.target.value;
        }}
      />
    );
  }
}

let ShortTextEditor: ContentEditor = {
  shouldRender: (content: API.Content) =>
    content.short != null && (content.short.type == 'text' || content.short.type == 'markdown'),
  render: (content: API.Content) => <ShortTextEditorComponent content={content} />,
};

let ShortHTMLEditor: ContentEditor = {
  shouldRender: (content: API.Content) => content.short != null && content.short.type == 'html',
  render: (content: API.Content) => <QuillComponent content={content} />,
};

let editors: ContentEditor[] = [
  LongHTMLEditor,
  LongMarkdownEditor,
  LongTextEditor,
  ShortHTMLEditor,
  ShortTextEditor,
];

export let Editor: React.SFC<{
  content: API.Content | API.Data;
  hideLabel: boolean;
}> = (props) => {
  for (var i = 0; i < editors.length; i++) {
    let editor = editors[i];
    if (editor.shouldRender(props.content)) {
      return (
        <div className="controls">
          <div className="control control-full">
            {props.hideLabel ? '' : <div className="label">{props.content.key}</div>}
            {editor.render(props.content)}
          </div>
        </div>
      );
    }
  }
};
