import * as m from 'mithril';
import msx from 'lib/msx';
import * as API from 'lib/api';
import { BaseComponent } from 'components/auth';

interface TextEditorAttrs {
  elementId?: string;
  readonly content: API.Content;
  readonly short?: boolean;
}

export default class TextEditorComponent extends BaseComponent<TextEditorAttrs> {
  private readonly _attrs: TextEditorAttrs;
  content: API.Content;
  element: HTMLElement;
  id: string;
  readonly short: boolean;

  constructor(v: m.CVnode<TextEditorAttrs>) {
    super(v);
    this.content = v.attrs.content;
    if (v.attrs.elementId == null) {
      v.attrs.elementId = Math.random().toString().slice(2, 10);
    }
    this.id = `#text-${v.attrs.elementId}`;
    this.short = v.attrs.short;
  }

  get klass(): string {
    let k = '.text';
    if (this.short) {
      return k + '.text-short';
    }
    return k;
  }

  view() {
    return <div id={this.id} class={this.klass}>
      <textarea
        onchange={(el: Event) => {
          this.content.value = (el.target as any).value;
        }}
        oncreate={(v: m.VnodeDOM<any, any>) => {
          (v.dom as HTMLTextAreaElement).value = this.content.value;
        }}
      />
    </div>;
  }
}
