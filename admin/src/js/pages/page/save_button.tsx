import * as m from 'mithril';
import msx from 'lib/msx';
import * as API from 'lib/api';
import Page from 'lib/page';
import * as Toaster from 'components/toaster';
import { BaseComponent } from 'components/auth';

interface PageSaveButtonAttrs {
  page: Page;
  onsave: (page: Page) => void;
  classes?: string;
}

export default class PageSaveButtonComponent extends BaseComponent<PageSaveButtonAttrs> {
  constructor(v: m.CVnode<PageSaveButtonAttrs>) {
    super(v);
  }

  save() {
    let page = this.props.page;
    page.save().then((p: API.Page) => {
      page.uuid = p.uuid;
      window.history.replaceState(null, page.title, `/admin/pages/${p.uuid}`);
      return page.saveRoutes();
    })
      .then(() => {
        Toaster.add('Page successfully saved');
        this.props.onsave(page);
      })
      .catch((err: any) => {
        if (err.detail) {
          Toaster.add(err.detail, 'error');
        } else {
          Toaster.add('Internal server error.', 'error');
        }
      });
  }

  view() {
    return <a
      class={`button button--green ${this.props.classes || ''}`}
      onclick={(e: Event) => { e.stopPropagation(); this.save(); }}
    >
      Save
    </a>;
  }
}
