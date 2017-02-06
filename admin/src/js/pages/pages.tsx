import msx from 'lib/msx';
import Page from 'lib/page';
import Layout from 'components/layout';
import { MustAuthController } from 'components/auth';

export default class PagesPage extends MustAuthController {
  pages: Mithril.Property<Page[]>;
  constructor() {
    super();
    this.pages = m.prop([]);
    Page.list().then((pages) => this.pages(pages));
  }
  static controller = PagesPage;
  static view(ctrl: PagesPage) {
    return Layout(<div class='pages'>
      <header>
        <a class='button button--green button--center'
          href='/admin/compose'
          config={m.route}
        >
          Compose
        </a>
        <h1>Pages</h1>
      </header>
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