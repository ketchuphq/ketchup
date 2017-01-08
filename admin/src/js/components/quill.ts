import * as API from 'lib/api';
import * as Quill from 'quill';

export default class QuillComponent {
  quill: any;
  content: API.Content;
  element: HTMLElement;
  _id: string;
  maximize: Mithril.BasicProperty<boolean>;
  dark: Mithril.BasicProperty<boolean>;

  constructor(content: API.Content, readonly short: boolean = false) {
    this.content = content;
    this._id = 'quill' + Math.random().toString().slice(2, 10);
    this.maximize = m.prop(false);
    this.dark = m.prop(false);
  }

  get id(): string {
    return `#${this._id}`;
  }

  get klass(): string {
    let k = '.quill';
    if (this.short) {
      return k + '.quill-short';
    }
    return k;
  }

  initializeQuill() {
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
  }

  static controller = QuillComponent;
  static view(ctrl: QuillComponent) {
    return m(ctrl.klass, {
      class: [
        ctrl.maximize() ? 'ql-container-full' : '',
        ctrl.dark() ? 'ql-container-full-dark' : ''
      ].join(' ')
    },
      m(`${ctrl.id}-toolbar`, [
        m('.ql-formats',
          m('select.ql-size'),
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
        m('.ql-formats.qlx-maximize',
          m('a', {
            onclick: () => {
              ctrl.maximize(!ctrl.maximize());
            }
          }, 'Zen mode')
        )
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
      }),
      ctrl.maximize() ? m('.qlx-controls', [
        m('span.typcn.typcn-adjust-contrast', {
          onclick: () => { ctrl.dark(!ctrl.dark()); }
        }),
        m('span.typcn.typcn-times', {
          onclick: () => { ctrl.maximize(false); }
        }),
      ]) : ''
    );
  }
}