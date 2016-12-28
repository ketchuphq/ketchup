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
  showControls: Mithril.Property<boolean>;
  constructor() {
    this.showControls = m.prop(false);
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
    return m('.infoset', {
      class: this.showControls() ? 'hidden' : '',
      onclick: () => { this.showControls(true); }
    }, [
        this.routes().length == 0 || !this.routes()[0].path ? '' :
          m('.small.black8', this.routes()[0].path),
        m('.small.black8', [
          'Theme: ', this.page().theme, ', ',
          'Template: ', this.page().template,
        ]),
        m('.save-publish', [
          m('a.button.button--small.button--green-outline', {
            onclick: (e: Event) => { e.stopPropagation(); this.save(); }
          }, 'Save'),
          m('a.button.button--small.button--blue-outline', {
            onclick: () => { this.save(); }
          }, 'Publish')
        ])
      ]);
  }

  renderSettingsEditor() {
    return m('.controlset', {
      class: this.showControls() ? '' : 'hidden'
    }, [
        m('.settings', [
          m('.controls',
            m('.control',
              m.component(EditRoutesComponent, this.routes(), () => this.page().name),
            ),
            m('.control', {
              onclick: () => { this.showControls(false) }
            }, 'close')
          ),
          m('.controls',
            m.component(
              ThemePickerComponent,
              this.page().theme,
              this.page().template,
              this.updateThemeTemplate.bind(this)
            )
          )
        ]),
        m('.save-publish', [
          m('a.button.button--small.button--green-outline', {
            onclick: (e: Event) => { e.stopPropagation(); this.save(); }
          }, 'Save'),
          m('a.button.button--small.button--blue-outline', {
            onclick: () => { this.save(); }
          }, 'Publish')
        ])
      ]);
  }
  renderEditors() {
    if (!this.page()) {
      return m('div', []);
    }
    var hideLabel = this.page().contents.length == 1;

    return m('.controls', this.page().contents.map((content) => {
      if (content.contentType == 'html') {
        return m('.control.control-full', [
          hideLabel ? '' : m('.label', content.key),
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
        ctrl.renderSettings(),
        ctrl.renderSettingsEditor(),
        m('.controls',
          m('input[type=text].large', {
            placeholder: 'title...',
            value: ctrl.page().name || '',
            onchange: m.withAttr('value', (v: string) => {
              ctrl.page().name = v;
            })
          }),
        ),
        ctrl.renderEditors()
      ])
    );
  }
}