import msx from 'lib/msx';

export default class Popover {
  visible: Mithril.Property<boolean>;

  constructor() {
    this.visible = m.prop(false);
  }

  static controller = Popover;
  static view(ctrl: Popover, _: any, children: any[]) {
    let content = [
      <a onclick={() => ctrl.visible(!ctrl.visible())}>{children[0]}</a>
    ];
    if (ctrl.visible()) {
      content.push(
        <div class='popover'>{children.slice(1)}</div>
      );
    }
    return <div class='popover-outer'>{content}</div>;
  }
}