import msx from 'lib/msx';
import * as m from 'mithril';
import Page from 'lib/page';
import Route from 'lib/route';
import { BaseComponent } from 'components/auth';

interface PagePickerAttrs {
  onselect: (option: Page) => void;
}

export class PagePickerComponent extends BaseComponent<PagePickerAttrs> {
  pages: Page[];
  selected: string;
  onselect: (option: Page) => void;

  constructor(v: m.CVnode<PagePickerAttrs>) {
    super(v);
    this.pages = [];
    this.onselect = v.attrs.onselect;
    Page.list().then((pages) => {
      this.pages = pages;
      if (pages.length > 0) {
        this.selected = pages[0].uuid;
        this.onselect(pages[0]);
      }
    });
  }

  view() {
    return <select
      value={this.selected}
      onchange={m.withAttr('value', (v) => {
        this.selected = v;
        for (var i = 0; i < this.pages.length; i++) {
          let page = this.pages[i];
          if (page.uuid == v) {
            this.onselect(page);
            return;
          }
        }
      })}
    >
    {this.pages.map((page) =>
      <option>{page.uuid}</option>
    )}
    </select>;
  }
}

export class NewRouteComponent extends BaseComponent {
  route: Route;

  constructor(v: any) {
    super(v);
    this.route = new Route();
  }

  selectPage(page: Page) {
    this.route.pageUuid = page.uuid;
  }

  view() {
    return <div class='new-route'>
      <input type='text'
        placeholder='route name'
        onchange={m.withAttr('value', (v) => {
          this.route.path = v;
        })}
      />
      <PagePickerComponent onselect={ this.selectPage.bind(this) } />
    </div>;
  }
}
