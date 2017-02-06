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
          let status = page.isPublished ? '' : 'draft ';
          let time = page.formattedUpdatedAt;
          if (time && !page.isPublished) {
            time = '@ ' + time;
          }
          return <a class='tr'
            href={`/admin/pages/${page.uuid}`}
            config={m.route}
          >
            <div>{page.title || 'untitled'}</div>
            <div class='small black5'>{`${status} ${time}`}</div>
          </a>;
        })}
      </div>
    </div>);
  }
}