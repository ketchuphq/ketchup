import * as m from 'mithril';
import msx from 'lib/msx';
import { BaseComponent } from 'components/auth';

export interface ModalAttrs {
  title: string;
  klass?: string;
  visible: () => boolean;
  toggle: () => void;
}

export class ModalComponent<A extends ModalAttrs = ModalAttrs> extends BaseComponent<A> {
  protected controls: m.Children;

  toggle(e?: Event) {
    if (e && e.srcElement.className != 'overlay') {
      return;
    }
    this.props.toggle();
  }

  view(v: m.CVnode<ModalAttrs>) {
    if (!this.props.visible()) {
      return <div />;
    }
    return (
      <div class='overlay' onclick={(e: Event) => this.toggle(e)}>
        <div class='modal-pad' />
        <div class={`modal ${this.props.klass || ''}`}>
          <div class='modal__title'>{this.props.title}</div>
          <div class='modal__content'>{v.children}</div>
          {this.controls}
        </div>
      </div>
    );
  }
}

type ModalButtonColor = '' | 'modal-button--green' | 'modal-button--red';

interface ConfirmModalAttrs extends ModalAttrs {
  resolve?: () => Promise<any>;
  reject?: () => Promise<any>;
  confirmText?: string;
  cancelText?: string;
  confirmColor?: ModalButtonColor;
  cancelColor?: ModalButtonColor;
  hideCancel?: boolean;
}

export class ConfirmModalComponent extends ModalComponent<ConfirmModalAttrs> {
  confirmColor: ModalButtonColor;
  cancelColor: ModalButtonColor;
  constructor(v: m.CVnode<ConfirmModalAttrs>) {
    super(v);
    this.confirmColor = this.props.confirmColor || 'modal-button--green';
    this.cancelColor = this.props.cancelColor || '';
  }

  resolve() {
    let promise = Promise.resolve();
    if (this.props.resolve) {
      promise = this.props.resolve();
    }
    promise.then(() => {
      this.toggle();
    });
  }

  reject() {
    let promise = Promise.resolve();
    if (this.props.reject) {
      promise = this.props.reject();
    }
    promise.then(() => {
      this.toggle();
    });
  }

  view(v: m.CVnode<ConfirmModalAttrs>) {
    let confirm = (
      <div class={`modal-button ${this.confirmColor}`} onclick={() => this.resolve()}>
        {v.attrs.confirmText || 'Yes'}
      </div>
    );
    let cancel;
    if (!v.attrs.hideCancel) {
      cancel = (
        <div class={`modal-button ${this.cancelColor}`} onclick={() => this.reject()}>
          {v.attrs.cancelText || 'Cancel'}
        </div>
      );
    }
    this.controls = (
      <div class='modal__controls'>
        {confirm} {cancel}
      </div>
    );

    return super.view(v);
  }
}
