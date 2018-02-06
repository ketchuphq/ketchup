import * as m from 'mithril';
import msx from 'lib/msx';
import * as API from 'lib/api';
import { BaseComponent } from 'components/auth';

interface QuillAttrs {
  readonly elementId?: string;
  readonly content: API.Content;
  readonly short?: boolean;
}

export default class QuillComponent extends BaseComponent<QuillAttrs> {
  quill: any; // figure out how to set this correctly
  content: API.Content;
  id: string;
  readonly short: boolean;

  constructor(v: m.CVnode<QuillAttrs>) {
    super(v);
    this.content = v.attrs.content;
    this.id = `quill-${v.attrs.elementId}`;
    this.short = v.attrs.short;
  }

  get klass(): string {
    let k = 'quill';
    if (this.short) {
      return k + 'quill-short';
    }
    return k;
  }

  async initializeQuill(element: HTMLElement) {
    // todo: figure out how to type this correctly
    const Quill: any = await import(/* webpackChunkName: "quill" */ 'quill');
    this.quill = new Quill(`#${this.id}`, {
      placeholder: 'start typing...',
      theme: 'snow',
      modules: {
        toolbar: `#${this.id}-toolbar`,
      },
    });
    let updateContent = () => {
      let editor = element.getElementsByClassName('ql-editor')[0];
      this.content.value = editor.innerHTML;
    };
    this.quill.on('text-change', updateContent.bind(this));
  }

  view() {
    return (
      <div class={this.klass}>
        <div id={`${this.id}-toolbar`}>
          <div class='ql-formats'>
            <select class='ql-header' />
          </div>
          <div class='ql-formats'>
            <button class='ql-bold' />
            <button class='ql-italic' />
            <button class='ql-underline' />
            <button class='ql-link' />
          </div>
          <div class='ql-formats'>
            <button class='ql-list' value='ordered' />
            <button class='ql-list' value='bullet' />
          </div>
          <div class='ql-formats'>
            <button class='ql-clean' />
          </div>
        </div>
        <div
          id={this.id}
          oncreate={(v: m.VnodeDOM<any, any>) => {
            if (this.content.value) {
              v.dom.innerHTML = this.content.value;
            }
            this.initializeQuill(v.dom as HTMLElement);
          }}
        />
      </div>
    );
  }
}
