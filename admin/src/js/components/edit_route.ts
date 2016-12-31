import Route from 'lib/route';
import Page from 'lib/page';

export default class EditRoutesComponent {
  dirty: boolean;
  infer: () => string;

  constructor(public page: Page) {
    this.infer = () => {
      return this.dirty ? '' : page.name;
    };
    this.dirty = true;
    if (page.routes.length == 0) {
      page.routes.push(new Route());
      this.dirty = false;
    } else if (!page.routes[0].path) {
      this.dirty = false;
    }
  }

  routeEditor(route: Route, i: number) {
    return m('div', [
      m('input[type=text]', {
        placeholder: '/path/to/page',
        value: Route.format(this.infer() || route.path || ''),
        onchange: m.withAttr('value', (v) => {
          this.dirty = true;
          route.path = v;
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