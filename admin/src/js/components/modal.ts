export abstract class ModalComponent {
  constructor() {
  }

  abstract render(): Mithril.VirtualElement;

  static controller = ModalComponent;
  static view(ctrl: ModalComponent) {
    return m('.modal',
      m('.modal__contents', [
        m('.modal__contents__title'),
        m('.modal__contents__content', ctrl.render())
      ]),
    );
  };
}