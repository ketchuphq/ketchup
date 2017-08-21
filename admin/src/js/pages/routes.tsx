import msx from 'lib/msx';
import * as m from 'mithril';
import Route from 'lib/route';
import { BaseComponent } from 'components/auth';

export default class RoutesPage extends BaseComponent {
  routes: Route[];

  constructor(v: any) {
    super(v);
    this.routes = [];
    Route.list().then((data) => this.routes = data);
  }

  view() {
    return <div class='routes'>
      <h1>Routes</h1>
      <div class='table'>
        {this.routes.map((r) =>
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
