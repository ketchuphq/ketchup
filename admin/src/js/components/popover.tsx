import * as m from 'mithril';
import msx from 'lib/msx';
import { BaseComponent } from 'components/auth';

interface PopoverAttrs {
  visible: boolean;
}

export default class Popover extends BaseComponent<PopoverAttrs> {
  view(v: m.CVnode<PopoverAttrs>) {
    if (!v.attrs.visible) {
      return;
    }
    return <div class='popover-outer'>
      <div class='popover'>{v.children}</div>
    </div>;
  }
}
