import cloneDeep from 'lodash-es/cloneDeep';
import isEqual from 'lodash-es/isEqual';

import msx from 'lib/msx';
import * as m from 'mithril';
import * as API from 'lib/api';
import Page from 'lib/page';
import Theme from 'lib/theme';
import Popover from 'components/popover';
import { MustAuthController } from 'components/auth';
import { ConfirmModalComponent } from 'components/modal';

import PageThemePickerComponent from 'pages/page/theme_picker';
import PageEditRoutesComponent from 'pages/page/edit_route';
import PageButtonsComponent from 'pages/page/buttons';
import PageSaveButtonComponent from 'pages/page/save_button';
import PageEditorsComponent from 'pages/page/editors';

export default class PagePage extends MustAuthController {
  page: Page;
  showSettings: boolean;
  template: API.ThemeTemplate;
  _nextRoute: string;
  _clickStart: DOMTokenList; // keep track of click location to prevent firing on drag
  dirty: boolean;

  initialContent: API.Content[];

  constructor() {
    super();
    this.dirty = false;
    this.showSettings = false;
    let pageUUID = m.route.param('id');
    if (pageUUID) {
      Page.get(pageUUID)
        .then((page) => this.updatePage(page, true))
        .then((page) => page.getRoutes());
    } else {
      this.updatePage(new Page(), true);
    }
  }

  toggleSettings() {
    this.showSettings = !this.showSettings;
  }

  goToIndex() {
    this._nextRoute = '/admin/pages';
  }

  updatePage(page: Page, initial = false) {
    return this.updateThemeTemplate(page.theme, page.template)
      .then(() => {
        if (!this.page || initial) {
          this.initialContent = cloneDeep(page.contents);
          this.page = page;
        }
      })
      .then(() => this.updateContent(initial))
      .then(() => m.redraw())
      .then(() => page);
  }

  updateThemeTemplate(theme: string, template: string) {
    if (this.page) {
      this.page.theme = theme;
      this.page.template = template;
    }
    return Theme.get(theme).then((t) => {
      this.template = t.getTemplate(template);
      m.redraw();
    });
  }

  updateContent(initial = false) {
    // todo: keep old content temporarily for 'undo'
    let contentMap: { [key: string]: API.Content } = {};
    let placeholderContents: API.Content[] = [];
    let placeholderContentMap: { [key: string]: boolean } = {};
    (this.page.contents || []).forEach((c) => {
      // on initial load copy all fields
      if (c.uuid || c.value || initial) {
        contentMap[c.key] = c;
      }
    });
    (this.template.placeholders || []).forEach((p) => {
      if (contentMap[p.key]) {
        Object.keys(p).forEach((k: keyof API.ThemePlaceholder) => {
          if (p[k]) {
            contentMap[p.key][k] = p[k];
          }
        });
        placeholderContents.push(contentMap[p.key]);
      } else {
        placeholderContents.push(API.Content.copy(p, {}));
      }
      placeholderContentMap[p.key] = true;
    });

    let pageContents: API.Content[] = [];
    (this.page.contents || []).forEach((c) => {
      if (c.key == 'content' && this.template.hideContent) {
        return;
      }
      if (contentMap[c.key] && !placeholderContentMap[c.key]) {
        pageContents.push(c);
      }
    });
    this.page.contents = pageContents.concat(placeholderContents);
  }

  confirmLeave() {
    if (this.showSettings) {
      this.toggleSettings();
      return;
    }
    this.dirty = this.dirty || !isEqual(this.initialContent, this.page.contents);
    if (this.dirty) {
      ConfirmModalComponent.confirm({
        title: 'You are about to leave this page',
        content: () => <p>
          You have unsaved changes. Are you sure
        you want to leave this page?
      </p>,
        confirmText: 'Stay',
        cancelText: 'Leave',
        cancelColor: 'modal-button--red'
      })
        .then(() => { /* noop */ })
        .catch(() => { this.goToIndex(); });
      return;
    }
    this.goToIndex();
  }

  renderSettings() {
    return <div class='controlset'>
      <div class='settings'>
        <div class='controls'>
          <div class='control'>
            {this.page ? m(PageEditRoutesComponent, this.page) : null}
          </div>
        </div>
        <div class='controls'>
          <PageThemePickerComponent
            theme={this.page.theme}
            template={this.page.template}
            callback={(theme, template) => {
              this.updateThemeTemplate(theme, template)
                .then(() => this.updateContent());
            }}
          />
        </div>
      </div>
      <PageButtonsComponent page={this.page} onsave={(page: Page) => this.updatePage(page, true)} />
    </div>;
  }

  static oninit(v: Mithril.Vnode<{}, PagePage>) {
    v.state = new PagePage();
  };

  static view(v: Mithril.Vnode<{}, PagePage>) {
    let ctrl = v.state;
    if (!ctrl.page) {
      return;
    }

    let controls = [
      <PageSaveButtonComponent
        page={ctrl.page}
        onsave={(page: Page) => ctrl.updatePage(page, true)} />,
      <span class='typcn typcn-cog' onclick={() => ctrl.toggleSettings()} />,
      <Popover visible={ctrl.showSettings}>{ctrl.renderSettings()}</Popover>,
      <a class='typcn typcn-times'
        href='/admin/pages'
        onclick={(e: MouseEvent) => {
          e.preventDefault();
          ctrl.confirmLeave();
        }}
      />
    ];

    let pageMaxClasses = 'page-max animate-fade-in';
    if (ctrl._nextRoute) {
      pageMaxClasses = 'page-max animate-zoom-away animate-fill';
    }

    return <div class={pageMaxClasses}
      onclick={(e: any) => {
        let validClick = e.target.classList.contains('page-max') && ctrl._clickStart && ctrl._clickStart.contains('page-max');
        if (!validClick) {
          return;
        }
        ctrl.confirmLeave();
      }}
      oncreate={(v: Mithril.VnodeDOM<any, any>) => {
        v.dom.addEventListener('mousedown', (ev: any) => {
          ctrl._clickStart = ev.target.classList;
        });
        v.dom.addEventListener('animationend', (ev: AnimationEvent) => {
          // old animation is removed, otherwise new animations won't fire.
          if (ev.animationName == 'fadeIn') {
            v.dom.classList.add('animate-fade-in-complete');
            v.dom.classList.remove('animate-fade-in');
          }
          // navigate away after zoomAway animation completes
          if (ev.animationName == 'zoomAway') {
            m.route.set(ctrl._nextRoute);
          }
        });
      }}
    >
      <ConfirmModalComponent />
      <div class='page-max__controls'>{controls}</div>
      <div class='page-editor'
        oncreate={(v: Mithril.VnodeDOM<any, any>) => {
          v.dom.addEventListener('click', () => ctrl.showSettings = false);
        }}
      >
        <div class='controls'>
          <input
            type='text'
            class='large'
            placeholder='title...'
            value={ctrl.page.title || ''}
            onchange={m.withAttr('value', (v: string) => {
              ctrl.page.title = v;
            })}
          />
        </div>
        <PageEditorsComponent contents={ctrl.page.contents} />
      </div>
    </div>;
  }
}