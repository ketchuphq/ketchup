import * as m from 'mithril';
import * as API from 'lib/api';
import QuillComponent from 'components/editors/html';
import TextEditorComponent from 'components/editors/text';

type EditorType = 'quill' | 'cm';

interface EditorComponentAttrs {
  content: API.Content;
}

export default class EditorComponent {
  editor: EditorType;
  _id: string;

  constructor() {
    this.editor = 'quill';
    this._id = Math.random().toString().slice(2, 10);
  }

  static oninit(v: Mithril.Vnode<EditorComponentAttrs, EditorComponent>) {
    v.state = new EditorComponent();
  }

  static view(v: Mithril.Vnode<EditorComponentAttrs, EditorComponent>) {
    let ctrl = v.state;
    let data = { elementId: ctrl._id, content: v.attrs.content };
    if (ctrl.editor == 'quill') {
      return m(QuillComponent, data);
    }
    return m('div', [
      m('div', m('a.small', { onclick: () => ctrl.editor = 'quill' }, 'show editor')),
      m(TextEditorComponent, data)
    ]);
  }
}
