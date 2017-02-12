import msx from 'lib/msx';
import * as API from 'lib/api';
import Page from 'lib/page';
import Layout from 'components/layout';
import { MustAuthController } from 'components/auth';

export default class PagesPage extends MustAuthController {
  pages: Mithril.Property<Page[]>;
  viewOption: Mithril.Property<API.ListPageRequest_ListPageFilter>;
  constructor() {
    super();
    this.pages = m.prop([]);
    this.viewOption = m.prop<API.ListPageRequest_ListPageFilter>('all');
    this.fetch(this.viewOption());
  }

  fetch(val: API.ListPageRequest_ListPageFilter) {
    this.viewOption(val);
    return Page.list(val)
      .then((pages) => this.pages(pages))
      .then(() => {
        m.redraw();
      });
  }

  static controller = PagesPage;
  static view(ctrl: PagesPage) {
    let tab = (v: API.ListPageRequest_ListPageFilter, desc?: string) => {
      let classes = 'tab-el';
      if (ctrl.viewOption() == v) {
        classes += ' tab-selected';
      }
      return <span class={classes} onclick={() => ctrl.fetch(v)}>
        {desc || v}
      </span>;
    };
    return Layout(<div class='pages'>
      <header>
        <a class='button button--green button--center' href='/admin/compose' config={m.route}>
          Compose
        </a>
        <h1>Pages</h1>
      </header>
      <h2 class='tabs'>
        {tab('all')}
        <span class='tab-divider'>|</span>
        {tab('draft')}
        <span class='tab-divider'>|</span>
        {tab('published')}
      </h2>
      <div class='table'>
        {ctrl.pages().map((page) => {
          let status = null;
          let klass = '';
          if (page.isPublished) {
            status = <div class='label small'>published</div>;
          } else {
            status = <div class='label label--gray small'>draft</div>;
            klass = 'page--draft';
          }

          return <a class='tr tr--center'
            href={`/admin/pages/${page.uuid}`}
            config={m.route}
          >
            <div class={`tr__expand ${klass}`}>{page.title || 'untitled'}</div>
            {status}
            <div class='page--date'>{`${page.formattedUpdatedAt}`}</div>
          </a>;
        })}
      </div>
    </div>);
  }
}