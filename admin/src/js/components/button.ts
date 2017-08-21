import * as m from 'mithril';
import { BaseComponent } from 'components/auth';

interface ButtonAttrs {
  handler?: () => Promise<any>;
  onclick?: () => any;

  class?: string;
  id?: string;
  href?: string;
}

export default class Button extends BaseComponent<ButtonAttrs> {
  loading: boolean;
  handler: () => void;

  constructor(v: m.CVnode<ButtonAttrs>) {
    super(v);
    this.loading = false;
    this.handler = () => {
      if (!v.attrs.handler) {
        return;
      }
      this.loading = true;
      m.redraw();
      v.attrs
        .handler()
        .then(() => (this.loading = false))
        .catch(() => (this.loading = false));
    };
  }

  view(v: m.Vnode<ButtonAttrs, Button>) {
    let c: ButtonAttrs = {
      class: v.attrs.class,
      id: v.attrs.id,
      href: v.attrs.href
    };
    if (v.attrs.handler) {
      c.onclick = this.handler;
    }

    let loader = null;
    if (this.loading) {
      c = { class: c.class + ' button--loading' };
      loader = m('.loader', [m('.loading0'), m('.loading1'), m('.loading2')]);
    }
    return m('a.button', c, [loader, m('.button__inner', v.children)]);
  }
}
