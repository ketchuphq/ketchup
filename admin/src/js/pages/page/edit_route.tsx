import msx from 'lib/msx';
import * as m from 'mithril';
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

  static oninit(v: Mithril.Vnode<Page, EditRoutesComponent>) {
    v.state = new EditRoutesComponent(v.attrs);
  };

  static view(v: Mithril.Vnode<Page, EditRoutesComponent>) {
    let ctrl = v.state;
    return <div class='edit-route control'>
      <div class='label'>Permalink</div>
      {ctrl.page.routes.map(ctrl.routeEditor.bind(ctrl))}
    </div>;
  }
}