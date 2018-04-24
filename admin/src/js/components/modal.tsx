import * as React from 'react';
import * as ReactModal from 'react-modal';

export interface ModalProps {
  title: string;
  klass?: string;
  visible: boolean;
  toggle: () => void;
}

export class ModalComponent<A extends ModalProps = ModalProps> extends React.Component<A> {
  constructor(props: any) {
    super(props);
  }

  render() {
    return (
      <ReactModal isOpen={this.props.visible}>
        <div className={`modal ${this.props.klass || ''}`}>
          <div className="modal__title">{this.props.title}</div>
          <div className="modal__content">{this.props.children}</div>
        </div>
      </ReactModal>
    );
  }
}

type ModalButtonColor = '' | 'modal-button--green' | 'modal-button--red';

interface ConfirmModalProps extends ModalProps {
  resolve?: () => void;
  reject?: () => void;
  confirmText?: string;
  cancelText?: string;
  confirmColor?: ModalButtonColor;
  cancelColor?: ModalButtonColor;
  hideCancel?: boolean;
}
export class ConfirmModalComponent<
  A extends ConfirmModalProps = ConfirmModalProps
> extends React.Component<A> {
  constructor(props: any) {
    super(props);
  }

  resolve = (_: React.MouseEvent<HTMLElement>) => {
    this.props.toggle();
    if (this.props.resolve) {
      this.props.resolve();
    }
  };

  reject = (_: React.MouseEvent<HTMLElement>) => {
    this.props.toggle();
    if (this.props.reject) {
      this.props.reject();
    }
  };

  render() {
    let className = `modal ${this.props.klass || ''}`;
    return (
      <ReactModal
        portalClassName="modal-portal"
        overlayClassName="overlay"
        className={className}
        isOpen={this.props.visible}
      >
        <div className="modal__title">{this.props.title}</div>
        <div className="modal__content">{this.props.children}</div>
        <div className="modal__controls">
          <div
            className={`modal-button ${this.props.confirmColor || 'modal-button--green'}`}
            onClick={this.resolve}
          >
            {this.props.confirmText || 'Yes'}
          </div>{' '}
          {!this.props.hideCancel ? (
            <div className={`modal-button ${this.props.cancelColor}`} onClick={this.reject}>
              {this.props.cancelText || 'Cancel'}
            </div>
          ) : null}
        </div>
      </ReactModal>
    );
  }
}
