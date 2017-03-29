import msx from 'lib/msx';

interface PopoverAttrs {
  visible: boolean;
}

export default class Popover {
  private readonly _attrs: PopoverAttrs;

  static view({ attrs: {visible}, children }: Mithril.Vnode<PopoverAttrs, {}>) {
    if (!visible) {
      return;
    }
    return <div class='popover-outer'>
      <div class='popover'>{children}</div>
    </div>;
  }
}