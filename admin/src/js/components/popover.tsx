import msx from 'lib/msx';

export let hidePopover = () => {};

export default class Popover {
  visible: boolean;

  constructor() {
    this.visible = false;
    if (hidePopover != null) {
      hidePopover();
    }
    hidePopover = () => this.visible = false;
  }

  static oninit(v: Mithril.Vnode<{}, Popover>) {
    v.state = new Popover();
  };

  static view(v: Mithril.Vnode<{}, Popover>) {
    let ctrl = v.state;
    let content = [
      <a onclick={() => ctrl.visible = !ctrl.visible}>{v.children[0]}</a>
    ];
    if (ctrl.visible) {
      content.push(
        <div class='popover'>{v.children.slice(1)}</div>
      );
    }
    return <div class='popover-outer'>{content}</div>;
  }
}