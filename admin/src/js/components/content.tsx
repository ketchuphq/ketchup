import EditorComponent from 'components/editors/editor';
import QuillComponent from 'components/editors/html';
import CodeMirrorComponent from 'components/editors/markdown';
import ShortTextEditorComponent from 'components/editors/short_text';
import TextEditorComponent from 'components/editors/text';
import * as API from 'lib/api';
import * as React from 'react';

interface ContentEditor {
  shouldRender(content: API.Content): boolean;
  render(content: API.Content, ref: React.RefObject<any>): React.ReactElement<any>;
}

const LongHTMLEditor: ContentEditor = {
  shouldRender: (content: API.Content) => content.text != null && content.text.type == 'html',
  render: (content: API.Content, ref: React.RefObject<any>) => (
    <EditorComponent ref={ref} content={content} />
  ),
};

const LongMarkdownEditor: ContentEditor = {
  shouldRender: (content: API.Content) => content.text != null && content.text.type == 'markdown',
  render: (content: API.Content, ref: React.RefObject<any>) => (
    <CodeMirrorComponent ref={ref} content={content} />
  ),
};

const LongTextEditor: ContentEditor = {
  shouldRender: (content: API.Content) => content.text != null && content.text.type == 'text',
  render: (content: API.Content, ref: React.RefObject<any>) => (
    <TextEditorComponent ref={ref} content={content} />
  ),
};

const ShortTextEditor: ContentEditor = {
  shouldRender: (content: API.Content) =>
    content.short != null && (content.short.type == 'text' || content.short.type == 'markdown'),
  render: (content: API.Content, ref: React.RefObject<any>) => (
    <ShortTextEditorComponent ref={ref} content={content} />
  ),
};

const ShortHTMLEditor: ContentEditor = {
  shouldRender: (content: API.Content) => content.short != null && content.short.type == 'html',
  render: (content: API.Content, ref: React.RefObject<any>) => (
    <QuillComponent ref={ref} content={content} />
  ),
};

const editors: ContentEditor[] = [
  LongHTMLEditor,
  LongMarkdownEditor,
  LongTextEditor,
  ShortHTMLEditor,
  ShortTextEditor,
];

interface Props {
  content: API.Content | API.Data;
  hideLabel: boolean;
}

const imageExtensions = ['png', 'jpg', 'jpeg', 'gif'];

export class Editor extends React.PureComponent<Props> {
  editorRef: React.RefObject<any>;

  constructor(props: Props) {
    super(props);
    this.editorRef = React.createRef();
  }

  insertFile(filename: string, url: string) {
    const editor = this.editorRef.current;
    let isImage = imageExtensions.some((v) => filename.substr(-v.length) == v);
    if (isImage) {
      if (editor.insertImage) {
        editor.insertImage(url);
      }
    } else if (editor.insertText) {
      editor.insertText(filename, url);
    }
  }

  render() {
    for (var i = 0; i < editors.length; i++) {
      let editor = editors[i];
      if (editor.shouldRender(this.props.content)) {
        return (
          <div className="controls">
            <div className="control control-full">
              {this.props.hideLabel ? '' : <div className="label">{this.props.content.key}</div>}
              {editor.render(this.props.content, this.editorRef)}
            </div>
          </div>
        );
      }
    }
  }
}
