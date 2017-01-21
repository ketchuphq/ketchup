import * as API from 'lib/api';
import * as CodeMirror from 'CodeMirror';

export default class CodeMirrorComponent {
  codemirror: CodeMirror.Editor;
  content: API.Content;
  element: HTMLElement;
  dark: Mithril.BasicProperty<boolean>;
  id: string;

  constructor(_id: string, content: API.Content, readonly short: boolean = false) {
    this.content = content;
    this.id = `#codemirror-${_id}`;
    this.dark = m.prop(false);
  }

  get klass(): string {
    let k = '.codemirror';
    if (this.short) {
      return k + '.codemirror-short';
    }
    return k;
  }

  initializeCodeMirror(element: HTMLTextAreaElement) {
    element.value = this.content.value;
    this.codemirror = CodeMirror.fromTextArea(element, {
      mode: 'gfm',
      placeholder: 'start typing...',
      lineNumbers: false,
      theme: 'elegant',
      lineWrapping: true,
    });
    // this.codemirror.setOption('fullscreen', 'true')
    // this.codemirror.refresh()
    this.codemirror.on('change', (instance) => {
      this.content.value = instance.getValue();
    });
  }

  static controller = CodeMirrorComponent;
  static view(ctrl: CodeMirrorComponent) {
    return m(ctrl.id + ctrl.klass, m('textarea', {
      config: (el: HTMLTextAreaElement, isInitialized: boolean) => {
        if (!isInitialized) {
          ctrl.initializeCodeMirror(el);
        }
      }
    }));
  }
}