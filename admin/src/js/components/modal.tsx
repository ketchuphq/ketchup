import msx from 'lib/msx';

export interface ModalAttrs {
  title: string;
  klass?: string;
  content: () => any;
}

export class ModalComponent {
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
        <div class='.modal__contents'>
          <div class='modal__contents__title'>{content.title}</div>
          <div class='modal__contents__content'>{content.content()}</div>
          {v.children}
        </div>
      </div>
    </div>;
  };
}