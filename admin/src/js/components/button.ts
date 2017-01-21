interface ButtonAttributes extends Mithril.Attributes {
  handler?: () => Mithril.Promise<any>;
}

export default class Button {
  loading: Mithril.Property<boolean>;
  handler: () => void;

  constructor(config: ButtonAttributes) {
    this.loading = m.prop(false);
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

  static controller = Button;
  static view = (ctrl: Button, config: ButtonAttributes, ...children: any[]) => {
    let c: ButtonAttributes = {
      onclick: () => ctrl.handler(),
      ...config
    };

    let loader = null;
    if (ctrl.loading()) {
      c = { class: config.class + ' button--loading' };
      loader = m('.loader', [m('.loading0'), m('.loading1'), m('.loading2')]);
    }

    return m('a.button', c, [
      loader,
      m('.button__inner', children)
    ]);
  }
}