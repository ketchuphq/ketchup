import * as API from 'lib/api';
import Page from 'lib/page';
import Route from 'lib/route';
import Theme from 'lib/theme';
import Layout from 'components/layout';
import QuillComponent from 'components/quill';
import EditRoutesComponent from 'components/edit_route';
import ThemePickerComponent from 'components/theme_picker';
import * as Toaster from 'components/toaster';

export default class PagePage {
  page: Mithril.Property<Page>;
  routes: Mithril.Property<Route[]>;
  showControls: Mithril.Property<boolean>;
  template: Mithril.Property<API.ThemeTemplate>;

  constructor() {
    this.showControls = m.prop(false);
    this.page = m.prop<Page>();
    this.routes = m.prop([]);
    this.template = m.prop<API.ThemeTemplate>();
    let pageUUID = m.route.param('id');
    if (pageUUID) {
      Page.get(pageUUID).then((page) => {
        this.page(page);
        Theme.get(page.theme).then((t) => {
          this.template(t.getTemplate(page.template));
        });
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
    Theme.get(theme).then((t) => {
      this.template(t.getTemplate(template));
    });
  }

  save() {
    this.page()
      .save()
      .then((page: API.Page) => {
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

  publish() {
    this.page()
      .publish()
      .then(() => {
        Toaster.add('Page published');
        m.redraw();
      });
  }

  renderSavePublish() {
    return m('.save-publish', [
      m('a.button.button--small.button--green', {
        onclick: (e: Event) => { e.stopPropagation(); this.save(); }
      }, 'Save'),
      this.page().isPublished ?
        m('a.button.button--small.button--blue', {
          href: this.routes()[0].path
        }, 'View')
        :
        m('a.button.button--small.button--blue', {
          onclick: (e: Event) => { e.stopPropagation(); this.publish(); }
        }, 'Publish')
    ]);
  }
  renderSettings() {
    return m('.controls', {
      class: this.showControls() ? 'hidden' : '',
    }, [
        m('.infoset', {
          onclick: () => { this.showControls(true); }
        }, [
            m('.small.black8', [
              this.routes().length == 0 || !this.routes()[0].path ? '' :
                ['Path: ', m('strong', this.routes()[0].path), ', '],
              'Theme: ', m('strong', this.page().theme), ', ',
              'Template: ', m('strong', this.page().template),
            ]),
          ]),
        this.renderSavePublish()
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
              onclick: () => { this.showControls(false); }
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
        this.renderSavePublish()
      ]);
  }
  renderEditors() {
    if (!this.page()) {
      return m('div', []);
    }
    let placeholderKeys = Object.keys(this.template().placeholders);
    let contentKeys = this.page().contents.map((c) => c.key);

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