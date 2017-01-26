import * as API from 'lib/api';
import Page from 'lib/page';
import Theme from 'lib/theme';
import Layout from 'components/layout';
import EditRoutesComponent from 'components/edit_route';
import ThemePickerComponent from 'components/theme_picker';
import * as Toaster from 'components/toaster';
import { renderEditor } from 'components/content';
import { MustAuthController } from 'components/auth';

export default class PagePage extends MustAuthController {
  page: Mithril.Property<Page>;
  showControls: Mithril.Property<boolean>;
  template: Mithril.Property<API.ThemeTemplate>;
  maximize: Mithril.Property<boolean>;
  minimiseToSettings: Mithril.Property<boolean>;

  constructor() {
    super();
    this.showControls = m.prop(false);
    this.maximize = m.prop(true);
    this.page = m.prop<Page>();
    this.template = m.prop<API.ThemeTemplate>();
    this.minimiseToSettings = m.prop(false);
    let pageUUID = m.route.param('id');
    if (pageUUID) {
      Page.get(pageUUID)
        .then((page) => {
          if (!!page.theme) {
            Theme.get(page.theme).then((t) => {
              this.template(t.getTemplate(page.template));
            });
          }
          page.getRoutes()
            .then(() => {
              this.page(page);
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
    this.page().save().then((page: API.Page) => {
      this.page().uuid = page.uuid;
      window.history.replaceState(null, this.page().title, `/admin/pages/${page.uuid}`);
      return this.page().saveRoutes();
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
    // todo: handle case where not saved yet
    this.page().publish().then(() => {
      Toaster.add('Page published');
      m.redraw();
    });
  }

  unpublish() {
    this.page().unpublish().then(() => {
      Toaster.add('Page unpublished');
      m.redraw();
    });
  }

  delete() {
    this.page().delete().then(() => {
      Toaster.add('Page deleted', 'error');
      m.route('/admin/pages');
    });
  }

  renderSavePublish() {
    return m('.save-publish', [
      m('a.button.button--small.button--green', {
        onclick: (e: Event) => { e.stopPropagation(); this.save(); }
      }, 'Save'),
      !this.page().uuid ? '' :
        m('a.button.button--small.button--red', {
          onclick: (e: Event) => { e.stopPropagation(); this.delete(); }
        }, 'Delete'),
      this.page().isPublished ?
        [
          m('a.button.button--small.button--blue', {
            onclick: (e: Event) => { e.stopPropagation(); this.unpublish(); }
          }, 'Unpublish'),
          m('a.button.button--small.button--blue', {
            href: this.page().defaultRoute
          }, 'View'),
        ]
        :
        m('a.button.button--small.button--blue', {
          onclick: (e: Event) => { e.stopPropagation(); this.publish(); }
        }, 'Publish')
    ]);
  }
  renderSettings() {
    return m('.controlset', {
      class: this.showControls() ? 'hidden' : '',
    }, [
        m('.infoset', {
          onclick: () => { this.showControls(true); }
        }, [
            m('.small.black5', [
              !this.page().defaultRoute ? '' :
                ['Path: ', m('strong', this.page().defaultRoute), ', '],
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
              this.page() ? m.component(EditRoutesComponent, this.page()) : null,
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

  // todo: if not part of the theme, show delete button
  renderEditors() {
    if (!this.page()) {
      return m('div', []);
    }

    let contentMap: { [key: string]: API.Content } = {};
    let mainContent: API.Content;
    let contents: API.Content[] = [];
    (this.page().contents || []).forEach((c) => {
      contentMap[c.key] = c;
      if (c.key != 'content') {
        contents.push(c);
      } else {
        mainContent = c;
      }
    });

    let placeholders: API.ThemePlaceholder[] = [];
    let filteredPlaceholders: API.ThemePlaceholder[] = [];
    let hideContent = false;

    if (this.template()) {
      hideContent = this.template().hideContent;
      placeholders = (this.template().placeholders || []);
      let placeholderOrder: { [key: string]: number } = {};

      placeholders.forEach((p, i) => {
        placeholderOrder[p.key] = i;
        if (p.key == 'content' && !mainContent.value) {
          API.ContentText.copy(p.text, mainContent.text);
          // todo: if p.text is not set, set mainContent.text to html?
        }
        if (!contentMap[p.key]) {
          filteredPlaceholders.push(p);
        }
      });

      contents.sort((a, b) => {
        if (placeholderOrder[a.key] < placeholderOrder[b.key]) {
          return - 1;
        }
        if (placeholderOrder[a.key] > placeholderOrder[b.key]) {
          return 1;
        }
        return 0;
      });
    }

    return m('div', [
      contents.map((c) =>
        renderEditor(this.page(), c, false)
      ),
      filteredPlaceholders.map((p) =>
        renderEditor(this.page(), p, false)
      ),
      !hideContent && mainContent ? renderEditor(this.page(), mainContent, true) : null
    ]);
  }

  maxView() {
    let controls = [
      m('a.button.button--green', {
        onclick: (e: Event) => { e.stopPropagation(); this.save(); }
      }, 'Save'),
      m('span.typcn.typcn-cog', {
        onclick: () => { this.maximize(false); }
      }),
    ];
    if (!this.minimiseToSettings()) {
      controls.push(m('a.typcn.typcn-times', {
        href: '/admin/pages',
        config: m.route
      }));
    }

    return Layout(
      m('.page-max',
        m('.page-max__controls', controls),
        m('.page-editor', [
          m('.controls',
            m('input[type=text].large', {
              placeholder: 'title...',
              value: this.page().title || '',
              onchange: m.withAttr('value', (v: string) => {
                this.page().title = v;
              })
            }),
          ),
          this.renderEditors()
        ])
      )
    );
  }

  static controller = PagePage;
  static view(ctrl: PagePage) {
    if (ctrl.maximize()) {
      return ctrl.maxView();
    }
    return Layout(
      m('.page-editor', [
        m('.page-editor__icons', [
          m('span.typcn.typcn-arrow-maximise', {
            onclick: () => {
              ctrl.minimiseToSettings(true);
              ctrl.maximize(true);
            }
          }),
        ]),
        m('.controls',
          m('input[type=text].large', {
            placeholder: 'title...',
            value: ctrl.page().title || '',
            onchange: m.withAttr('value', (v: string) => {
              ctrl.page().title = v;
            })
          }),
        ),
        ctrl.renderSettings(),
        ctrl.renderSettingsEditor(),
        ctrl.renderEditors()
      ])
    );
  }
}