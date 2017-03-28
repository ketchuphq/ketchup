import cloneDeep from 'lodash-es/cloneDeep';
import isEqual from 'lodash-es/isEqual';

import msx from 'lib/msx';
import * as m from 'mithril';
import * as API from 'lib/api';
import Page from 'lib/page';
import Theme from 'lib/theme';
import { hidePopover, default as Popover } from 'components/popover';
import { MustAuthController } from 'components/auth';
import { ConfirmModalComponent } from 'components/modal';

import PageThemePickerComponent from 'pages/page/theme_picker';
import PageEditRoutesComponent from 'pages/page/edit_route';
import PageButtonsComponent from 'pages/page/buttons';
import PageSaveButtonComponent from 'pages/page/save_button';
import PageEditorsComponent from 'pages/page/editors';

export default class PagePage extends MustAuthController {
  page: Page;
  showControls: boolean;
  template: API.ThemeTemplate;
  _nextRoute: string;
  _clickStart: DOMTokenList; // keep track of click location to prevent firing on drag
  dirty: boolean;

  initialContent: API.Content[];

  constructor() {
    super();
    this.showControls = false;
    this.dirty = false;
    let pageUUID = m.route.param('id');
    if (pageUUID) {
      Page.get(pageUUID)
        .then((page) => this.updatePage(page))
        .then((page) => {
          this.initialContent = cloneDeep(page.contents);
          return page;
        })
        .then((page) => page.getRoutes());
    } else {
      this.updatePage(new Page());
    }
  }

  updatePage(page: Page) {
    return this.updateThemeTemplate(page.theme, page.template)
      .then(() => this.page = page)
      .then(() => this.updateContent())
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

  updateContent() {
    // todo: keep old content temporarily
    let contentMap: { [key: string]: API.Content } = {};
    let placeholderContents: API.Content[] = [];
    let placeholderContentMap: { [key: string]: boolean } = {};

    (this.page.contents || []).forEach((c) => {
      if (c.uuid || c.value) {
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

  renderSettings(show: boolean = false) {
    return <div class={this.showControls || show ? 'controlset' : 'controlset hidden'}>
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
      <PageButtonsComponent page={this.page} />
    </div>;
  }

  static oninit(v: Mithril.Vnode<{}, PagePage>) {
    v.state = new PagePage();
  };

  static view(v: Mithril.Vnode<{}, PagePage>) {
    let ctrl = v.state;
    if (!ctrl.page) {
      return; // loading
    }

    let controls = [
      <PageSaveButtonComponent page={ctrl.page} />,
      <Popover>
        <span class='typcn typcn-cog' />
        <div>{ctrl.renderSettings(true)}</div>
      </Popover>,
      <a class='typcn typcn-times'
        href='/admin/pages'
        onclick={(e: MouseEvent) => {
          e.preventDefault();
          ctrl._nextRoute = '/admin/pages';
        }}
      />
    ];

    let pageMaxClasses = 'page-max animate-fade-in';
    if (ctrl._nextRoute) {
      pageMaxClasses = 'page-max animate-zoom-away';
    }

    return <div class={pageMaxClasses}
      onclick={(e: any) => {
        let validClick = e.target.classList.contains('page-max') && ctrl._clickStart && ctrl._clickStart.contains('page-max');
        ctrl.dirty = ctrl.dirty || !isEqual(ctrl.initialContent, ctrl.page.contents);
        if (validClick) {
          if (ctrl.dirty) {
            ConfirmModalComponent.confirm({
              title: 'You have unsaved changes',
              content: () => 'do you want to save or close?',
              confirmText: 'Save',
              cancelText: 'Cancel'
            })
              .then(() => { console.log('woot'); })
              .catch(() => { console.log('woops'); });
            return;
            // confirm
          }
          ctrl._nextRoute = '/admin/pages';
        }
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
          v.dom.addEventListener('click', () => hidePopover());
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