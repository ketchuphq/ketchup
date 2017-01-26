import * as API from 'lib/api';
import QuillComponent from 'components/editors/html';
import TextEditorComponent from 'components/editors/text';

type EditorType = 'quill' | 'cm';

export default class EditorComponent {
  editor: Mithril.BasicProperty<EditorType>;
  _id: string;

  constructor() {
    this.editor = m.prop<EditorType>('quill');
    this._id = Math.random().toString().slice(2, 10);
  }

  static controller = EditorComponent;

  static view(ctrl: EditorComponent, content: API.Content): Mithril.VirtualElement {
    if (ctrl.editor() == 'quill') {
      return m.component(QuillComponent, ctrl._id, content) as any;
    }
    return m('div', [
      m('div', m('a.small', { onclick: () => ctrl.editor('quill') }, 'show editor')),
      m.component(TextEditorComponent, ctrl._id, content)
    ]);
  }
}