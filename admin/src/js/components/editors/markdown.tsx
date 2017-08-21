import * as m from 'mithril';
import msx from 'lib/msx';
import * as API from 'lib/api';
import { BaseComponent } from '../auth';

interface CodeMirrorAttrs {
  readonly content: API.Content;
  readonly short?: boolean;
}

export default class CodeMirrorComponent extends BaseComponent<CodeMirrorAttrs> {
  codemirror: CodeMirror.Editor;
  content: API.Content;
  element: HTMLElement;
  id: string;
  short: boolean;

  constructor(v: any) {
    super(v);
    this.content = v.attrs.content;
    this.short = v.attrs.short;
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
    ], (require) => {
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

  view() {
    return <div id={this.id} class={this.klass}>
      <textarea
        oncreate={(v: m.VnodeDOM<any, any>) => {
          this.initializeCodeMirror(v.dom as HTMLTextAreaElement);
        }}
      />
    </div>;
  }
}
