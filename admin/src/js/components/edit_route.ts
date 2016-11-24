import Route from 'lib/route';

export default class EditRoutesComponent {
  routes: Route[];
  allowMany: boolean;
  // hasUserEnteredRoute?

  // accept a parameter to watch for changes so we know when to save?
  constructor(routes: Route[], allowMany: boolean = false) {
    this.routes = routes;
    this.allowMany = allowMany;
  }

  static controller = EditRoutesComponent;
  static view(ctrl: EditRoutesComponent) {
    return m('.edit-route', [
      ctrl.routes.map((route: Route) => {
        return m('div',
          m('input[type=text]', {
            placeholder: '/path/to/page',
            value: route.path || '',
            onchange: m.withAttr('value', (v) => {
              route.path = v;
            })
          })
        );
      }),
      !(ctrl.allowMany || ctrl.routes.length == 0) ? '' :
        m('a.button.button--green', {
          onclick: () => { ctrl.routes.push(new Route()); }
        }, 'Add Permalink')
    ]);
  }
}