import msx from 'lib/msx';
import * as API from 'lib/api';
import Page from 'lib/page';
import Theme from 'lib/theme';
import Layout from 'components/layout';
import { hidePopover, default as Popover } from 'components/popover';
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
    let saveButton = <a
      class='button button--small button--green'
      onclick={(e: Event) => { e.stopPropagation(); this.save(); }}
    >
      Save
    </a>;

    let unpublishButton = <a
      class='button button--small button--blue'
      onclick={(e: Event) => { e.stopPropagation(); this.unpublish(); }}
    >
      Unpublish
    </a>;

    let viewButton = <a
      class='button button--small button--blue'
      href={this.page().defaultRoute}
    >
      View
    </a>;

    let deleteButton = <a
      class='button button--small button--red'
      onclick={(e: Event) => { e.stopPropagation(); this.delete(); }}
    >
      Delete
    </a>;

    let publishButton = <a
      class='button button--small button--blue'
      onclick={(e: Event) => { e.stopPropagation(); this.publish(); }}
    >
      Publish
    </a>;

    return <div class='save-publish'>
      {saveButton}
      {!this.page().uuid ? '' : deleteButton}
      {this.page().isPublished ? [unpublishButton] : publishButton}
    </div>;
  }
  renderSettings() {
    let path = ['Path: ', <strong>{this.page().defaultRoute}</strong>, ', '];
    return <div class={!this.showControls() ? 'controlset' : 'controlset hidden'}>
      <div class='infoset' onclick={() => { this.showControls(true); }}>
        <div class='small black5'>
          {!this.page().defaultRoute ? '' : path}
          Theme: <strong>{this.page().theme}</strong>,
          Template: <strong>{this.page().template}</strong>
        </div>
        {this.renderSavePublish()}
      </div>
    </div>;
  }

  renderSettingsEditor(show: boolean = false) {
    return <div class={this.showControls() || show ? 'controlset' : 'controlset hidden'}>
      <div class='settings'>
        <div class='controls'>
          <div class='control'>
            {this.page() ? m.component(EditRoutesComponent, this.page()) : null}
          </div>
          {this.maximize() ? '' :
            <div class='control'
              onclick={() => { this.showControls(false); }}>
              close
            </div>}
        </div>
        <div class='controls'>
          {m.component(
            ThemePickerComponent,
            this.page().theme,
            this.page().template,
            this.updateThemeTemplate.bind(this)
          )}
        </div>
      </div>
      {this.renderSavePublish()}
    </div>;
  }

  // todo: if not part of the theme, show delete button
  renderEditors() {
    if (!this.page()) {
      return <div></div>;
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
        } else {
          // update the oneof type of the content from the placeholder
          // todo: exhaustively map fields
          // todo: convert markdown <> css more elegantly
          contentMap[p.key].multiple = p.multiple;
          contentMap[p.key].text = p.text;
          contentMap[p.key].short = p.short;
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

    return <div>
      {contents.map((c) => renderEditor(this.page(), c, false))}
      {filteredPlaceholders.map((p) => renderEditor(this.page(), p, false))}
      {!hideContent && mainContent ? renderEditor(this.page(), mainContent, true) : null}
    </div>;
  }

  maxView() {
    let controls = [
      <a class='button button--green'
        onclick={(e: Event) => { e.stopPropagation(); this.save(); }}>
        Save
      </a>,
      <Popover>
        <span class='typcn typcn-cog' />
        <div>{this.renderSettingsEditor(true)}</div>
      </Popover>
    ];
    if (!this.minimiseToSettings()) {
      controls.push(
        <a class='typcn typcn-times' href='/admin/pages' config={m.route} />
      );
    }

    return Layout(
      <div class='page-max' onclick={(e: any) => {
        if (e.target.className == 'page-max') {
          m.route('/admin/pages');
        }
      }}>
        <div class='page-max__controls'>{controls}</div>
        <div class='page-editor'
          config={(el, isInitialized) => {
            if (!isInitialized) {
              el.addEventListener('focus', () => hidePopover(), true);
            }
          }}>
          <div class='controls'>
            <input
              type='text'
              class='large'
              placeholder='title...'
              value={this.page().title || ''}
              onchange={this.updatePageTitle}
            />
          </div>
          {this.renderEditors()}
        </div>
      </div>
    );
  }

  updatePageTitle = m.withAttr('value', (v: string) => {
    this.page().title = v;
  });

  static controller = PagePage;
  static view(ctrl: PagePage) {
    if (ctrl.maximize()) {
      return ctrl.maxView();
    }
    let maximize = () => {
      ctrl.minimiseToSettings(true);
      ctrl.maximize(true);
    }
    return Layout(
      <div class='page-editor'>
        <div class='page-editor__icons'>
          <span class='typcn typcn-arrow-maximise' onclick={maximize} />
        </div>
        <div class='controls'>
          <input
            class='large'
            type='text'
            placeholder='title...'
            value={ctrl.page().title || ''}
            onchange={ctrl.updatePageTitle}
          />
        </div>
        {ctrl.renderSettings()}
        {ctrl.renderSettingsEditor()}
        {ctrl.renderEditors()}
      </div>
    );
  }
}