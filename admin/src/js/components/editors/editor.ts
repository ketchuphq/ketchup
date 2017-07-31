import * as m from 'mithril';
import * as API from 'lib/api';
import QuillComponent from 'components/editors/html';
import TextEditorComponent from 'components/editors/text';
import { BaseComponent } from 'components/auth';

type EditorType = 'quill' | 'cm';

interface EditorComponentAttrs {
  content: API.Content;
}

export default class EditorComponent extends BaseComponent<EditorComponentAttrs> {
  editor: EditorType;
  _id: string;

  constructor(v: m.CVnode<EditorComponentAttrs>) {
    super(v);
    this.editor = 'quill';
    this._id = Math.random().toString().slice(2, 10);
  }

  view() {
    let data = { elementId: this._id, content: this.props.content };
    if (this.editor == 'quill') {
      return m('.editor', m(QuillComponent, data));
    }
    return m('.editor', [
      m('div', m('a.small', { onclick: () => this.editor = 'quill' }, 'show editor')),
      m(TextEditorComponent, data)
    ]);
  }
}
