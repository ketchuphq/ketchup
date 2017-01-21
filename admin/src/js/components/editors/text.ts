import * as API from 'lib/api';

export default class TextEditorComponent {
  content: API.Content;
  element: HTMLElement;
  id: string;

  constructor(_id: string, content: API.Content, readonly short: boolean = false) {
    this.content = content;
    if (_id == null) {
      _id = Math.random().toString().slice(2, 10);
    }
    this.id = `#text-${_id}`;
  }

  get klass(): string {
    let k = '.text';
    if (this.short) {
      return k + '.text-short';
    }
    return k;
  }

  static controller = TextEditorComponent;
  static view(ctrl: TextEditorComponent) {
    return m(ctrl.id + ctrl.klass, m('textarea', {
      config: (el: HTMLTextAreaElement, isInitialized: boolean) => {
        if (!isInitialized) {
          el.value = ctrl.content.value;
        }
      },
      onchange: (el: Event) => {
        ctrl.content.value = (el.target as any).value;
      }
    }));
  }
}