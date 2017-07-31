import * as m from 'mithril'
import msx from 'lib/msx';
import { BaseComponent } from 'components/auth';

export interface ModalAttrs {
  title: string;
  klass?: string;
  visible: () => boolean;
  toggle: () => void;
}

export class ModalComponent extends BaseComponent<ModalAttrs> {
  view(v: m.CVnode<ModalAttrs>) {
    // need to get attrs here, not in constructor which only
    // gets called the first time.
    if (!v.attrs.visible()) {
      return <div></div>;
    }
    return <div
      class='overlay'
      onclick={v.attrs.toggle}
    >
      <div class='modal-pad' />
      <div class={`modal ${v.attrs.klass || ''}`}>
        <div class='modal__title'>{v.attrs.title}</div>
        <div class='modal__content'>
          {v.children}
        </div>
      </div>
    </div>;
  };
}

interface ConfirmModalAttrs {
  title: string;
  klass?: string;
  content: () => any;
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
