import * as m from 'mithril';
import msx from 'lib/msx';
import { BaseComponent } from 'components/auth';

interface PopoverAttrs {
  visible: boolean;
}

export default class Popover extends BaseComponent<PopoverAttrs> {
  view(v: m.CVnode<PopoverAttrs>) {
    let klass = 'popover-outer';
    if (!v.attrs.visible) {
      klass += ' popover-outer-hidden';
    }
    return (
      <div class={klass}>
        <div class='popover'>{v.children}</div>
      </div>
    );
  }
}
