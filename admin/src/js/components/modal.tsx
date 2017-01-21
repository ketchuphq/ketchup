import msx from 'lib/msx';

export interface ModalContent {
  title: string;
  klass?: string;
  content: () => any;
}

export class ModalComponent {
  static controller = ModalComponent;
  static view(_: ModalComponent, content: ModalContent, ...children: any[]) {
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
        <div class='.modal__contents'>
          <div class='modal__contents__title'>{content.title}</div>
          <div class='modal__contents__content'>{content.content()}</div>
          {children}
        </div>
      </div>
    </div>;
  };
}