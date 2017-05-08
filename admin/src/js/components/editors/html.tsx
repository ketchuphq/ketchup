import msx from 'lib/msx';
import * as API from 'lib/api';

interface QuillAttrs {
  readonly elementId?: string;
  readonly content: API.Content;
  readonly short?: boolean;
}

export default class QuillComponent {
  private readonly _attrs: QuillAttrs;
  quill: Quill.Quill;
  content: API.Content;
  id: string;
  readonly short: boolean;

  constructor(attrs: QuillAttrs) {
    this.content = attrs.content;
    this.id = `quill-${attrs.elementId}`;
    this.short = attrs.short;
  }

  get klass(): string {
    let k = 'quill';
    if (this.short) {
      return k + 'quill-short';
    }
    return k;
  }

  initializeQuill(element: HTMLElement) {
    require.ensure(['quill'], (require) => {
      let Quill: Quill.Quill = require<Quill.Quill>('quill');
      this.quill = new Quill(`#${this.id}`, {
        placeholder: 'start typing...',
        theme: 'snow',
        modules: {
          toolbar: `#${this.id}-toolbar`
        }
      });
      let updateContent = () => {
        let editor = element.getElementsByClassName('ql-editor')[0];
        this.content.value = editor.innerHTML;
      };
      this.quill.on('text-change', updateContent.bind(this));
    }, 'quill');
  }

  static oninit(v: Mithril.Vnode<QuillAttrs, QuillComponent>) {
    v.state = new QuillComponent(v.attrs);
  }

  static view(v: Mithril.Vnode<QuillAttrs, QuillComponent>) {
    let ctrl = v.state;
    return <div class={ctrl.klass}>
      <div id={`${ctrl.id}-toolbar`}>
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
      <div id={ctrl.id}
        oncreate={(v: Mithril.VnodeDOM<any, any>) => {
          if (ctrl.content.value) {
            v.dom.innerHTML = ctrl.content.value;
          }
          ctrl.initializeQuill(v.dom as HTMLElement);
        }}
      >
      </div>
    </div>;
  }
};

let _: Mithril.Component<QuillAttrs, QuillComponent> = QuillComponent;
