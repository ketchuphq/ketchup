import * as API from 'lib/api';
import * as Quill from 'quill';

export default class QuillComponent {
  quill: any;
  content: API.Content;
  element: HTMLElement;
  _id: string;

  constructor(content: API.Content, readonly short: boolean = false) {
    this.content = content;
    this._id = 'quill' + Math.random().toString().slice(2, 10);
  }

  get id(): string {
    return `#${this._id}`;
  }

  get klass(): string {
    let k = '.quill';
    if (this.short) {
      return k + '.quill-short';
    }
    return k;
  }

  initializeQuill() {
    this.quill = new Quill(this.id, {
      placeholder: 'start typing...',
      theme: 'snow'
    });
    this.quill.on('text-change', () => {
      let editor = this.element.getElementsByClassName('ql-editor')[0];
      this.content.value = editor.innerHTML;
    });
  }

  static controller = QuillComponent;
  static view(ctrl: QuillComponent) {
    return m(ctrl.klass,
      m(ctrl.id, {
        config: (el: HTMLElement, isInitialized: boolean) => {
          if (!isInitialized) {
            ctrl.element = el;
            if (ctrl.content.value) {
              ctrl.element.innerHTML = m.trust(ctrl.content.value) as any as string;
            }
            ctrl.initializeQuill();
          }
        }
      })
    );
  }
}