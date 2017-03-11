import msx from 'lib/msx';
import * as m from 'mithril';
import Page from 'lib/page';
import Route from 'lib/route';

interface PagePickerAttrs {
  onselect: (option: Page) => void;
}

class PagePickerComponent {
  _attrs: PagePickerAttrs;
  pages: Page[];
  selected: string;
  onselect: (option: Page) => void;

  constructor(attrs: PagePickerAttrs) {
    this.pages = [];
    this.onselect = attrs.onselect;
    Page.list().then((pages) => {
      this.pages = pages;
      if (pages.length > 0) {
        this.selected = pages[0].uuid;
        this.onselect(pages[0]);
      }
    });
  }


  static oninit(v: Mithril.Vnode<PagePickerAttrs, PagePickerComponent>) {
    v.state = new PagePickerComponent(v.attrs);
  };

  static view(v: Mithril.Vnode<PagePickerAttrs, PagePickerComponent>) {
    let ctrl = v.state;
    return <select
      value={ctrl.selected}
      onchange={m.withAttr('value', (v) => {
        ctrl.selected = v;
        for (var i = 0; i < ctrl.pages.length; i++) {
          let page = ctrl.pages[i];
          if (page.uuid == v) {
            ctrl.onselect(page);
            return;
          }
        }
      })}
    >
    {ctrl.pages.map((page) =>
      <option>{page.uuid}</option>
    )}
    </select>;
  }
}

class NewRouteComponent {
  route: Route;

  constructor() {
    this.route = new Route();
  }

  selectPage(page: Page) {
    this.route.pageUuid = page.uuid;
  }

  static oninit(v: Mithril.Vnode<{}, NewRouteComponent>) {
    v.state = new NewRouteComponent();
  };

  static view(v: Mithril.Vnode<{}, NewRouteComponent>) {
    let ctrl = v.state;
    return <div class='new-route'>
      <input type='text'
        placeholder='route name'
        onchange={m.withAttr('value', (v) => {
          ctrl.route.path = v;
        })}
      />
      <PagePickerComponent onselect={ ctrl.selectPage.bind(ctrl) } />
    </div>;
  }
}