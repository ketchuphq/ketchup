import Route from 'lib/route';
import Layout from 'components/layout';

export default class RoutesPage {
  routes: Mithril.BasicProperty<Route[]>;

  constructor() {
    this.routes = m.prop([]);
    Route.list().then((data) => this.routes(data));
  }

  static controller = RoutesPage;
  static view(ctrl: RoutesPage) {
    return Layout(
      m('.routes',
        m('h1', 'Routes'),
        m('ul',
          ctrl.routes().map((r) =>
            m('li', [
              m('a', { href: r.path ? r.path : '#' }, r.path),
              !r.pageUuid ? '' :
                m('a.list-link', {
                  href: `/admin/pages/${r.pageUuid}`,
                  config: m.route
                }, 'edit page'),
              m('a.list-link', { href: r.path }, 'delete')
            ])
          )
        )
      )
    );
  }
}