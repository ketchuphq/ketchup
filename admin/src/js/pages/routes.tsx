import msx from 'lib/msx';
import * as m from 'mithril';
import Route from 'lib/route';

export default class RoutesPage {
  routes: Route[];

  constructor() {
    this.routes = [];
    Route.list().then((data) => this.routes = data);
  }

  static oninit(v: Mithril.Vnode<{}, RoutesPage>) {
    v.state = new RoutesPage();
  }

  static view(v: Mithril.Vnode<{}, RoutesPage>) {
    let ctrl = v.state;
    return <div class='routes'>
      <h1>Routes</h1>
      <div class='table'>
        {ctrl.routes.map((r) =>
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
    </div>;
  }
}