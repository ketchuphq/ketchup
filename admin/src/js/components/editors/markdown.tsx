import msx from 'lib/msx';
import * as API from 'lib/api';

let _: Mithril.Component<CodeMirrorAttrs, CodeMirrorComponent> = CodeMirrorComponent;

interface CodeMirrorAttrs {
  readonly content: API.Content;
  readonly short?: boolean;
}

export default class CodeMirrorComponent {
  private readonly _attrs: CodeMirrorAttrs;
  codemirror: CodeMirror.Editor;
  content: API.Content;
  element: HTMLElement;
  id: string;
  short: boolean;

  constructor(attrs: CodeMirrorAttrs) {
    this.content = attrs.content;
    this.short = attrs.short;
    this.id = `codemirror-${Math.random().toString().slice(2, 10)}`;
  }

  get klass(): string {
    return [
      'codemirror',
      this.short ? 'codemirror-short' : 'codemirror-long'
    ].join(' ');
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

      element.value = this.content.value || '';
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
    }, 'codemirror');
  }

  static oninit(v: Mithril.Vnode<CodeMirrorAttrs, CodeMirrorComponent>) {
    v.state = new CodeMirrorComponent(v.attrs);
  };

  static view(v: Mithril.Vnode<CodeMirrorAttrs, CodeMirrorComponent>) {
    let ctrl = v.state;
    ctrl.content = v.attrs.content; // for some reason we lose the reference
    return <div id={ctrl.id} class={ctrl.klass}>
      <textarea
        oncreate={(v: Mithril.VnodeDOM<any, any>) => {
          ctrl.initializeCodeMirror(v.dom as HTMLTextAreaElement);
        }}
      />
    </div>;
  }
}