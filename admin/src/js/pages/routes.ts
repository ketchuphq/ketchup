import Route from 'lib/route';
import Page from 'lib/page';
import Layout from 'components/layout';

class PagePickerComponent {
  pages: Mithril.Property<Page[]>;
  selected: Mithril.Property<string>;
  onselect: (option: Page) => void;

  constructor(onselect: (option: Page) => void) {
    this.selected = m.prop<string>();
    this.pages = m.prop([]);
    this.onselect = onselect;
    Page.list().then((pages) => {
      this.pages(pages);
      if (pages.length > 0) {
        this.selected(pages[0].uuid);
        this.onselect(pages[0]);
      }
    });
  }

  static controller = PagePickerComponent;
  static view(ctrl: PagePickerComponent) {
    return m('select', {
      value: ctrl.selected(),
      onchange: m.withAttr('value', (v) => {
        ctrl.selected(v);
        for (var i = 0; i < ctrl.pages().length; i++) {
          let page = ctrl.pages()[i];
          if (page.uuid == v) {
            ctrl.onselect(page);
            return;
          }
        }
      })
    }, ctrl.pages().map((page) => {
      return m('option', page.uuid);
    }));
  }
}

class NewRouteComponent {
  route: Route;

  constructor() {
    this.route = new Route();
  }

  selectPage(page: Page) {
    this.route.page_uuid = page.uuid;
  }

  static controller = NewRouteComponent;
  static view(ctrl: NewRouteComponent) {
    return m('.new-route', [
      m('input[type=text]', {
        placeholder: 'route name',
        onchange: m.withAttr('value', (v) => {
          ctrl.route.path = v;
        })
      }),
      m.component(PagePickerComponent, ctrl.selectPage.bind(ctrl))
    ]);
  }
}

export default class RoutesPage {
  routes: Mithril.BasicProperty<Route[]>;

  constructor() {
    this.routes = m.prop([]);
    Route.list().then((data) => this.routes(data));
  }

  static controller = RoutesPage;
  static view(ctrl: RoutesPage) {
    return Layout(
      m('.routes',
        m('h1', 'Routes'),
        m.component(NewRouteComponent),
        m('ul',
          ctrl.routes().map((r) =>
            m('li',
              m('a', {
                href: r.path
              }, r.path)
            )
          )
        )
      )
    );
  }
}