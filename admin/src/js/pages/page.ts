import Page from 'lib/page';
import Route from 'lib/route';
import Layout from 'components/layout';
import QuillComponent from 'components/quill';
import EditRoutesComponent from 'components/edit_route';
import ThemePickerComponent from 'components/theme_picker';

export default class PagePage {
  page: Mithril.Property<Page>;
  routes: Mithril.Property<Route[]>;
  showRouteEditor: boolean;
  constructor() {
    this.showRouteEditor = false;
    this.page = m.prop<Page>();
    this.routes = m.prop([]);
    let pageUUID = m.route.param('id');
    if (pageUUID) {
      Page.get(pageUUID).then((page) => {
        this.page(page);
        page.getRoutes().then((r) => {
          this.routes(r);
        });
      });
    } else {
      this.page(new Page());
    }
  }


  updateThemeTemplate(theme: string, template: string) {
    let page = this.page()
    page.theme = theme;
    page.template = template;
    this.page(page);
  }

  renderEditors() {
    if (!this.page()) {
      return [];
    }

    return this.page().contents.map((content) => {
      if (content.contentType == 'html') {
        return m.component(QuillComponent, content);
      }
      return null;
    });
  }

  save() {
    this.page()
      .save()
      .then((page: Page) => {
        this.page().uuid = page.uuid;
        Route.saveList(this.routes(), page.uuid);
      });
  }

  static controller = PagePage;
  static view(ctrl: PagePage) {
    return Layout(
      m('.page-editor', [
        m('.controls',
          m('input[type=text].large', {
            placeholder: 'title...',
            value: ctrl.page().name || '',
            onchange: m.withAttr('value', (v: string) => {
              ctrl.page().name = v;
            })
          }),
          m('a.button.button--green', {
            onclick: () => { ctrl.save(); }
          }, 'Save Changes')
        ),
        m('.controls',
          m('.control',
            m.component(EditRoutesComponent, ctrl.routes()),
          ),
        ),
        m('.controls',
          m.component(
            ThemePickerComponent,
            ctrl.page().theme,
            ctrl.page().template,
            ctrl.updateThemeTemplate.bind(ctrl)
          )
        ),
        ctrl.renderEditors()
      ])
    );
  }
}