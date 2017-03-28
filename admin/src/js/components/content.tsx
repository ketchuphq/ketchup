import msx from 'lib/msx';
import * as API from 'lib/api';
import EditorComponent from 'components/editors/editor';
import QuillComponent from 'components/editors/html';
import TextEditorComponent from 'components/editors/text';
import CodeMirrorComponent from 'components/editors/markdown';

interface ContentEditor {
  shouldRender(content: API.Content): boolean;
  render(content: API.Content): Mithril.Vnode<any, any>;
}

let LongHTMLEditor: ContentEditor = {
  shouldRender: (content: API.Content) =>
    content.text != null && content.text.type == 'html',
  render: (content: API.Content) =>
    <EditorComponent content={content} />
};

let LongMarkdownEditor: ContentEditor = {
  shouldRender: (content: API.Content) =>
    content.text != null && content.text.type == 'markdown',
  render: (content: API.Content) =>
    <CodeMirrorComponent content={content} />
};

let LongTextEditor: ContentEditor = {
  shouldRender: (content: API.Content) =>
    content.text != null && content.text.type == 'text',
  render: (content: API.Content) =>
    <TextEditorComponent content={content} />
};

let ShortTextEditor: ContentEditor = {
  shouldRender: (content: API.Content) =>
    content.short != null && (content.short.type == 'text' || content.short.type == 'markdown'),
  render: (content: API.Content) =>
    <input type='text'
      oncreate={(v: Mithril.VnodeDOM<any, any>) => {
        (v.dom as HTMLTextAreaElement).value = content.value || '';
      }}
      onchange={(e: EventTarget) => {
        content.value = (e as any).target.value;
      }}
    />
};

let ShortHTMLEditor: ContentEditor = {
  shouldRender: (content: API.Content) =>
    content.short != null && content.short.type == 'html',
  render: (content: API.Content) =>
    <QuillComponent content={content} />
};

let editors: ContentEditor[] = [
  LongHTMLEditor,
  LongMarkdownEditor,
  LongTextEditor,
  ShortHTMLEditor,
  ShortTextEditor,
];

export function renderEditor(c: API.Content, hideLabel: boolean): Mithril.Vnode<any, any> {
  for (var i = 0; i < editors.length; i++) {
    let editor = editors[i];
    if (editor.shouldRender(c)) {
      return <div class='controls'>
        <div class='control control-full'>
          {hideLabel ? '' : <div class='label'>{c.key}</div>}
          {editor.render(c)}
        </div>
      </div>;
    }
  }

  console.log('no editor defined for object:', JSON.stringify(c));
}
