import Route from 'lib/route';
import Page from 'lib/page';

export default class EditRoutesComponent {
  dirty: boolean;

  constructor(public page: Page) {
    this.dirty = true;
    if (page.routes.length == 0) {
      page.routes.push(new Route());
      this.dirty = false;
    } else if (!page.routes[0].path) {
      this.dirty = false;
    }
  }

  infer() {
    if (this.dirty) {
      return;
    }
    if (this.page.routes.length < 1) {
      return;
    }
    if (!!this.page.routes[0].path) {
      return;
    }
    if (!this.page.name || this.page.name.trim() == '') {
      return;
    }
    this.page.routes[0].path = Route.format(this.page.name);
  };

  routeEditor(route: Route, i: number) {
    this.infer();
    return m('div', [
      m('input[type=text]', {
        placeholder: '/path/to/page',
        value: Route.format(route.path),
        onchange: m.withAttr('value', (v) => {
          this.dirty = true;
          route.path = Route.format(v);
        })
      }),
      i == 0 ? '' :
        m('a' as any as Mithril.Component<{}>, {
          onclick: () => this.page.routes.splice(i, 1)
        }, m.trust('&times;'))
    ]);
  }

  static controller = EditRoutesComponent;
  static view(ctrl: EditRoutesComponent) {
    return m('.edit-route.control', [
      m('.label', 'Permalink'),
      ctrl.page.routes.map(ctrl.routeEditor.bind(ctrl))
    ]);
  }
}