import * as API from 'lib/api';
import QuillComponent from 'components/editors/html';
import TextEditorComponent from 'components/editors/text';

type EditorType = 'quill' | 'cm';

export default class EditorComponent {
  maximize: Mithril.BasicProperty<boolean>;
  dark: Mithril.BasicProperty<boolean>;
  editor: Mithril.BasicProperty<EditorType>;
  _id: string;

  constructor() {
    this.maximize = m.prop(false);
    this.dark = m.prop(false);
    this.editor = m.prop<EditorType>('quill');
    this._id = Math.random().toString().slice(2, 10);
  }

  static controller = EditorComponent;

  static view(ctrl: EditorComponent, content: API.Content) {
    let editor: Mithril.Component<{}>;
    let toggle: Mithril.VirtualElement;
    if (ctrl.editor() == 'quill') {
      toggle = m('a.small', { onclick: () => ctrl.editor('cm') }, 'show html');
      editor = m.component(QuillComponent, ctrl._id, content);
    } else {
      toggle = m('a.small', { onclick: () => ctrl.editor('quill') }, 'show editor');
      editor = m.component(TextEditorComponent, ctrl._id, content);
    }

    return m('div', {
      class: [
        ctrl.maximize() ? 'ql-container-full' : '',
        ctrl.dark() ? 'ql-container-full-dark' : ''
      ].join(' ')
    }, [
        m('div', [
          toggle,
          m('.ql-formats.qlx-maximize',
            m('a', {
              onclick: () => {
                ctrl.maximize(!ctrl.maximize());
              }
            }, 'Zen mode')
          )
        ]),
        editor,
        ctrl.maximize() ? m('.qlx-controls', [
          m('span.typcn.typcn-adjust-contrast', {
            onclick: () => { ctrl.dark(!ctrl.dark()); }
          }),
          m('span.typcn.typcn-times', {
            onclick: () => { ctrl.maximize(false); }
          }),
        ]) : ''
      ]);
  }
}