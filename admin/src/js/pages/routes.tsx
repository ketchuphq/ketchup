import msx from 'lib/msx';
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
      <div class='routes'>
        <h1>Routes</h1>
        <div class='table'>
          {ctrl.routes().map((r) =>
            <div class='tr'>
              <a href={r.path ? r.path : '#'}>{r.path}</a>
              {
                !r.pageUuid ? '' :
                  <a class='list-link'
                    href={`/admin/pages/${r.pageUuid}`}
                    config={m.route}
                  >
                    edit page
                  </a>
              }
            </div>)}
        </div>
      </div>
    );
  }
}