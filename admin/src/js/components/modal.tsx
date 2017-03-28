import msx from 'lib/msx';

export interface ModalAttrs {
  title: string;
  klass?: string;
  content: () => any;
}

export class ModalComponent {
  private readonly _attrs: ModalAttrs;
  static view(v: Mithril.Vnode<ModalAttrs, {}>) {
    let content = v.attrs;
    if (!content || Object.keys(content).length == 0) {
      return <div></div>;
    }
    return <div
      class='overlay'
      onclick={() => {
        delete content.title;
        delete content.content;
        delete content.klass;
      }}
    >
      <div class='modal-pad' />
      <div class={`modal ${content.klass || ''}`}>
        <div class='modal__title'>{content.title}</div>
        <div class='modal__content'>{content.content()}</div>
        {v.children}
      </div>
    </div>;
  };
}

interface ConfirmModalAttrs extends ModalAttrs {
  confirmText?: string;
  cancelText?: string;
}

export class ConfirmModalComponent {
  private static content: ConfirmModalAttrs;
  private static resolve: () => void;
  private static reject: () => void;

  static confirm(content: ConfirmModalAttrs) {
    this.reset();
    this.content = content;

    return new Promise((resolve, reject) => {
      this.resolve = resolve;
      this.reject = reject;
    });
  }

  static reset() {
    this.content = null;
    this.resolve = null;
    this.reject = null;
  }

  static view(_: Mithril.Vnode<{}, {}>) {
    let content = ConfirmModalComponent.content;
    if (!content || Object.keys(content).length == 0) {
      return <div></div>;
    }
    return <div
      class='overlay'
      onclick={() => {
        ConfirmModalComponent.reset();
      }}
    >
      <div class='modal-pad' />
      <div class={`modal ${content.klass || ''}`}>
        <div class='modal__title'>{content.title}</div>
        <div class='modal__content'>{content.content()}</div>
        <div class='modal__controls'>
          <div class='modal-button modal-button--green' onclick={() => {
            ConfirmModalComponent.resolve();
            ConfirmModalComponent.reset();
          }}>
            {content.confirmText || 'Yes'}
          </div>
          <div class='modal-button' onclick={() => {
            ConfirmModalComponent.reject();
            ConfirmModalComponent.reset();
          }}>
            {content.cancelText || 'Cancel'}
          </div>
        </div>
      </div>
    </div>;
  };
}