import { Content } from '../lib/page';
import * as Quill from 'quill';

export default class QuillComponent {
  quill: any;
  content: Content;
  element: HTMLElement;
  constructor(content: Content) {
    this.content = content;
  }

  initializeQuill() {
    this.quill = new Quill('#editor', { theme: 'snow' });
    this.quill.on('text-change', () => {
      let editor = this.element.getElementsByClassName('ql-editor')[0];
      this.content.value = editor.innerHTML;
    });
  }

  static controller = QuillComponent;
  static view(ctrl: QuillComponent) {
    return m('.quill',
      m('#editor', {
        config: (el: HTMLElement, isInitialized: boolean) => {
          if (!isInitialized) {
            ctrl.element = el;
            ctrl.element.innerHTML = m.trust(ctrl.content.value) as any as string;
            ctrl.initializeQuill();
          }
        }
      })
    );
  }
}