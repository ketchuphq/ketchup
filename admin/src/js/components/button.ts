import * as m from 'mithril';
import * as stream from 'mithril/stream';

let _: Mithril.Component<ButtonAttrs, Button> = Button;
interface ButtonAttrs {
  handler?: () => Promise<any>;
  onclick?: () => any;

  class?: string;
  id?: string;
  href?: string;
}

export default class Button {
  loading: Mithril.Stream<boolean>;
  handler: () => void;

  constructor(readonly config: ButtonAttrs) {
    this.loading = stream(false);
    this.handler = () => {
      if (!config.handler) {
        return;
      }
      this.loading(true);
      m.redraw();
      config.handler()
        .then(() => this.loading(false))
        .catch(() => this.loading(false));
    };
  }

  static oninit(v: Mithril.Vnode<ButtonAttrs, Button>) {
    v.state = new Button(v.attrs);
  }
  static view = (v: Mithril.Vnode<ButtonAttrs, Button>) => {
    let ctrl = v.state;
    let c: ButtonAttrs = {
      class: v.attrs.class,
      id: v.attrs.id,
      href: v.attrs.href,
    };
    if (v.attrs.handler) {
      c.onclick = ctrl.handler;
    }

    let loader = null;
    if (ctrl.loading()) {
      c = { class: c.class + ' button--loading' };
      loader = m('.loader', [m('.loading0'), m('.loading1'), m('.loading2')]);
    }
    return m('a.button', c, [
      loader,
      m('.button__inner', v.children)
    ]);
  }
}
