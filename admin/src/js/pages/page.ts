import { BasePage, default as Page } from 'lib/page';
import Route from 'lib/route';
import Layout from 'components/layout';
import QuillComponent from 'components/quill';
import EditRoutesComponent from 'components/edit_route';
import ThemePickerComponent from 'components/theme_picker';
import * as Toaster from 'components/toaster';

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
    let page = this.page();
    page.theme = theme;
    page.template = template;
    this.page(page);
  }

  renderSettings() {
    return [
      m('.controls',
        m('.control', this.page().name),
      ),
      m('.controls', [
        this.page().theme,
        this.page().template,
      ]
      )
    ];
  }

  renderSettingsEditor() {
    return [
      m('.controls',
        m('.control',
          m.component(EditRoutesComponent, this.routes(), () => this.page().name),
        ),
      ),
      m('.controls',
        m.component(
          ThemePickerComponent,
          this.page().theme,
          this.page().template,
          this.updateThemeTemplate.bind(this)
        )
      )
    ];
  }

  renderEditors() {
    if (!this.page()) {
      return m('div', []);
    }

    return m('.controls', this.page().contents.map((content) => {
      if (content.contentType == 'html') {
        return m('.control.control-full', [
          m('.label', content.key),
          m.component(QuillComponent, content)
        ]);
      }
      return null;
    }));
  }

  save() {
    this.page()
      .save()
      .then((page: BasePage) => {
        this.page().uuid = page.uuid;
        window.history.replaceState(null, this.page().name, `/admin/pages/${page.uuid}`);
        return this.page().saveRoutes(this.routes());
      })
      .then(() => {
        Toaster.add('Page successfully saved');
      })
      .catch((err: any) => {
        if (err.detail) {
          Toaster.add(err.detail, 'error');
        } else {
          Toaster.add('Internal server error.', 'error');
        }
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
        ctrl.renderSettingsEditor(),
        ctrl.renderEditors()
      ])
    );
  }
}