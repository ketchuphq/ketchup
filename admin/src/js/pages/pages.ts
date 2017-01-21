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
    return Layout(m('.pages', [
      m('header',
        m('a.button.button--green.button--center', {
          href: '/admin/compose',
          config: m.route
        }, 'Compose'),
        m('h1', 'Pages')
      ),
      m('.table',
        ctrl.pages().map((page) => {
          let status = page.isPublished ? '' : 'draft ';
          let time = page.formattedUpdatedAt;
          if (time && !page.isPublished) {
            time = '@ ' + time;
          }
          return m('a.tr', {
            href: `/admin/pages/${page.uuid}`,
            config: m.route
          }, [
              m('div', page.title || 'untitled'),
              m('.small.black5', `${status} ${time}`)
            ]);
        })
      )
    ]));
  }
}