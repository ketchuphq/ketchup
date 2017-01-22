import * as API from 'lib/api';

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
    require.ensure([
      'codemirror/lib/codemirror',
      'codemirror/mode/gfm/gfm',
      'codemirror/mode/markdown/markdown',
      'codemirror/addon/display/placeholder',
      'codemirror/addon/mode/overlay'
    ], () => {
      let cm = require<typeof CodeMirror>('codemirror');
      require('codemirror/mode/gfm/gfm');
      require('codemirror/mode/markdown/markdown');
      require('codemirror/addon/display/placeholder');
      require('codemirror/addon/mode/overlay');

      element.value = this.content.value;
      this.codemirror = cm.fromTextArea(element, {
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