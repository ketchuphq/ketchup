import * as React from 'react';
import * as API from 'lib/api';

interface Props {
  readonly content: API.Content;
  readonly short?: boolean;
}

export default class CodeMirrorComponent extends React.Component<Props> {
  codemirror: CodeMirror.Editor;
  textInput: React.RefObject<HTMLTextAreaElement>;
  constructor(props: Props) {
    super(props);
    this.textInput = React.createRef();
  }

  componentDidMount() {
    this.initializeCodeMirror(this.textInput.current);
  }

  insertImage(url: string, altText?: string) {
    this.codemirror.getDoc().replaceSelection(`![${altText || ''}](${url})`);
  }

  insertLink(text: string, url: string) {
    this.codemirror.getDoc().replaceSelection(`[${text}](${url})`);
  }

  get klass(): string {
    return ['codemirror', this.props.short ? 'codemirror-short' : 'codemirror-long'].join(' ');
  }

  async initializeCodeMirror(element: HTMLTextAreaElement) {
    const cmImports = await Promise.all([
      import(/* webpackChunkName: "codemirror" */ 'codemirror'),
      import(/* webpackChunkName: "codemirror" */ 'codemirror/mode/gfm/gfm'),
      import(/* webpackChunkName: "codemirror" */ 'codemirror/mode/markdown/markdown'),
      import(/* webpackChunkName: "codemirror" */ 'codemirror/addon/display/placeholder'),
      import(/* webpackChunkName: "codemirror" */ 'codemirror/addon/mode/overlay'),
    ]);

    let cm = cmImports[0];

    element.value = this.props.content.value || '';
    this.codemirror = cm.fromTextArea(element, {
      mode: 'gfm',
      placeholder: 'start typing...',
      lineNumbers: false,
      theme: 'elegant',
      lineWrapping: true,
    });
    this.codemirror.on('change', (instance) => {
      this.props.content.value = instance.getValue();
    });
  }

  render() {
    return (
      <div key="codemirror" className={this.klass}>
        <textarea ref={this.textInput} />
      </div>
    );
  }
}
