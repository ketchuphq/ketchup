import * as API from 'lib/api';

export default class QuillComponent {
  quill: Quill.Quill;
  content: API.Content;
  element: HTMLElement;
  id: string;
  dark: Mithril.BasicProperty<boolean>;

  constructor(_id: string, content: API.Content, readonly short: boolean = false) {
    this.content = content;
    this.id = `#quill-${_id}`;
    this.dark = m.prop(false);
  }

  get klass(): string {
    let k = '.quill';
    if (this.short) {
      return k + '.quill-short';
    }
    return k;
  }

  initializeQuill() {
    require.ensure(['quill'], (require) => {
      let Quill: Quill.Quill = require<Quill.Quill>('quill');
      this.quill = new Quill(this.id, {
        placeholder: 'start typing...',
        theme: 'snow',
        modules: {
          toolbar: `${this.id}-toolbar`
        }
      });
      this.quill.on('text-change', () => {
        let editor = this.element.getElementsByClassName('ql-editor')[0];
        this.content.value = editor.innerHTML;
      });
    }); // chunkName?
  }

  static controller = QuillComponent;
  static view(ctrl: QuillComponent) {
    return m(ctrl.klass,
      m(`${ctrl.id}-toolbar`, [
        m('.ql-formats',
          m('select.ql-header'),
        ),
        m('.ql-formats',
          m('button.ql-bold'),
          m('button.ql-italic'),
          m('button.ql-underline'),
          m('button.ql-link'),
        ),
        m('.ql-formats',
          m('button.ql-list', { value: 'ordered' }),
          m('button.ql-list', { value: 'bullet' }),
        ),
        m('.ql-formats',
          m('button.ql-clean'),
        ),
      ]),
      m(ctrl.id, {
        config: (el: HTMLElement, isInitialized: boolean) => {
          if (!isInitialized) {
            ctrl.element = el;
            if (ctrl.content.value) {
              ctrl.element.innerHTML = m.trust(ctrl.content.value) as any as string;
            }
            ctrl.initializeQuill();
          }
        }
      })
    );
  }
}