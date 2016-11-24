import Page from 'lib/page';
import Layout from 'components/layout';

export default class PagesPage {
  pages: Mithril.Property<Page[]>;
  constructor() {
    this.pages = m.prop([]);
    Page.list().then((pages) => this.pages(pages));
  }
  static controller = PagesPage;
  static view(ctrl: PagesPage) {
    return Layout(m('.pages', [
      m('h1', 'Pages'),
      m('table',
        ctrl.pages().map((page) => {
          return m('tr',
            m('td.link-cell',
              m('a', {
                href: `/admin/pages/${page.uuid}`,
                config: m.route
              }, page.name || 'untitled')
            )
          );
        })
      )
    ]));
  }
}