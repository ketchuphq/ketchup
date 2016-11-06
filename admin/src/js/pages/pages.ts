import Layout from '../components/layout';
import Page from '../lib/page';

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
            m('td',
              m('a', {
                href: `/admin/pages/${page.uuid}`,
                config: m.route
              }, page.name || 'untitled')
            ),
            m('td')
          );
        })
      )
    ]));
  }
}