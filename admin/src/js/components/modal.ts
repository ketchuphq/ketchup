export interface ModalContent {
  title: string;
  klass?: string;
  content: () => any;
}

export class ModalComponent {
  constructor() {
  }

  static controller = ModalComponent;
  static view(ctrl: ModalComponent, content: ModalContent, ...children: any[]) {
    if (!content || Object.keys(content).length == 0) {
      return m('div');
    }
    return m('.overlay', {
      onclick: () => {
        delete content.title;
        delete content.content;
        delete content.klass;
      }
    }, [
        m('.modal-pad'),
        m(`.modal`, {
          class: content.klass || ''
        }, m('.modal__contents', [
          m('.modal__contents__title', content.title),
          m('.modal__contents__content', content.content()),
          ...children
        ]))
      ]);
  };
}