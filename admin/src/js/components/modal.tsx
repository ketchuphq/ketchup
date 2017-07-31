import * as m from 'mithril'
import msx from 'lib/msx';
import { BaseComponent } from 'components/auth';

export interface ModalAttrs {
  title: string;
  klass?: string;
  content: () => any;
}

export class ModalComponent extends BaseComponent<ModalAttrs> {
  content: ModalAttrs;
  constructor(v: m.CVnode<ModalAttrs>) {
    super(v);
    this.content = v.attrs;
  }

  view(v: m.CVnode<ModalAttrs>) {
    if (!this.content || Object.keys(this.content).length == 0) {
      return <div></div>;
    }
    return <div
      class='overlay'
      onclick={() => {
        delete this.content.title;
        delete this.content.content;
        delete this.content.klass;
      }}
    >
      <div class='modal-pad' />
      <div class={`modal ${this.content.klass || ''}`}>
        <div class='modal__title'>{this.content.title}</div>
        <div class='modal__content'>{this.content.content()}</div>
        {v.children}
      </div>
    </div>;
  };
}

interface ConfirmModalAttrs extends ModalAttrs {
  confirmText?: string;
  cancelText?: string;
  confirmColor?: ModalButtonColor;
  cancelColor?: ModalButtonColor;
}

type ModalButtonColor = '' | 'modal-button--green' | 'modal-button--red'

export class ConfirmModalComponent extends BaseComponent {
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

  view() {
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
          <div class={`modal-button ${content.confirmColor || ''}`} onclick={() => {
            ConfirmModalComponent.resolve();
            ConfirmModalComponent.reset();
          }}>
            {content.confirmText || 'Yes'}
          </div>
          <div class={`modal-button ${content.cancelColor || ''}`} onclick={() => {
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
