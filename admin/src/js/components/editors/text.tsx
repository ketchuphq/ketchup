import msx from 'lib/msx';
import * as API from 'lib/api';

let _: Mithril.Component<TextEditorAttrs, TextEditorComponent> = TextEditorComponent;

interface TextEditorAttrs {
  elementId?: string;
  readonly content: API.Content;
  readonly short?: boolean;
}

export default class TextEditorComponent {
  content: API.Content;
  element: HTMLElement;
  id: string;
  readonly short: boolean;

  constructor(attrs: TextEditorAttrs) {
    this.content = attrs.content;
    if (attrs.elementId == null) {
      attrs.elementId = Math.random().toString().slice(2, 10);
    }
    this.id = `#text-${attrs.elementId}`;
    this.short = attrs.short;
  }

  get klass(): string {
    let k = '.text';
    if (this.short) {
      return k + '.text-short';
    }
    return k;
  }

  static oninit(v: Mithril.Vnode<TextEditorAttrs, TextEditorComponent>) {
    v.state = new TextEditorComponent(v.attrs);
  };

  static view(v: Mithril.Vnode<TextEditorAttrs, TextEditorComponent>) {
    let ctrl = v.state;
    ctrl.content = v.attrs.content; // for some reason we lose the reference
    return <div id={ctrl.id} class={ctrl.klass}>
      <textarea
        onchange={(el: Event) => {
          ctrl.content.value = (el.target as any).value;
        }}
        oncreate={(v: Mithril.VnodeDOM<any, any>) => {
          (v.dom as HTMLTextAreaElement).value = ctrl.content.value;
        }}
      />
    </div>;
  }
}