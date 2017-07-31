import msx from 'lib/msx';
import * as m from 'mithril';
import Route from 'lib/route';
import Page from 'lib/page';
import { BaseComponent } from 'components/auth';

interface EditRoutesAttr {
  page: Page;
}

export default class EditRoutesComponent extends BaseComponent<EditRoutesAttr> {
  dirty: boolean;
  page: Page;

  constructor(v: m.CVnode<EditRoutesAttr>) {
    super(v)
    this.page = v.attrs.page;
    this.dirty = true;
    if (this.page.routes.length == 0) {
      this.page.routes.push(new Route());
      this.dirty = false;
    } else if (!this.page.routes[0].path) {
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
    if (!this.page.title || this.page.title.trim() == '') {
      return;
    }
    this.page.routes[0].path = Route.format(this.page.title);
  };

  routeEditor(route: Route, i: number) {
    this.infer();
    let r;
    if (i > 0) {
      r = <a onclick={() => this.page.routes.splice(i, 1)}>
        {m.trust('&times;')}
      </a>;
    }
    return <div>
      <input type='text'
        placeholder='/path/to/page'
        value={Route.format(route.path)}
        onchange={m.withAttr('value', (v) => {
          this.dirty = true;
          route.path = Route.format(v);
        })}
      />
      {r}
    </div>;
  }

  view() {
    return <div class='edit-route control'>
      <div class='label'>Permalink</div>
      {this.page.routes.map(this.routeEditor.bind(this))}
    </div>;
  }
}
