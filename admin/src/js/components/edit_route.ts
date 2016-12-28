import Route from 'lib/route';

export default class EditRoutesComponent {
  routes: Route[];
  dirty: boolean;
  infer: () => string;

  // accept a parameter to watch for changes so we know when to save?
  constructor(routes: Route[], infer: () => string) {
    this.infer = () => {
      return this.dirty ? '' : infer();
    };
    this.routes = routes;
    this.dirty = true;
    if (routes.length == 0) {
      this.routes.push(new Route());
      this.dirty = false;
    } else if (!routes[0].path) {
      this.dirty = false;
    }
  }

  static controller = EditRoutesComponent;
  static view(ctrl: EditRoutesComponent) {
    return m('.edit-route.control', [
      m('.label', 'Permalink'),
      ctrl.routes.map((route: Route, i: number) => {
        return m('div', [
          m('input[type=text]', {
            placeholder: '/path/to/page',
            value: Route.format(ctrl.infer() || route.path || ''),
            onchange: m.withAttr('value', (v) => {
              ctrl.dirty = true;
              route.path = v;
            })
          }),
          i == 0 ? '' :
            m('a' as any as Mithril.Component<{}>, {
              onclick: () => ctrl.routes.splice(i, 1)
            }, m.trust('&times;'))
        ]);
      })
    ]);
  }
}