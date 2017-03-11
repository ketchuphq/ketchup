import msx from 'lib/msx';
import * as m from 'mithril';
import * as API from 'lib/api';
import Page from 'lib/page';
import Theme from 'lib/theme';
import { hidePopover, default as Popover } from 'components/popover';
import * as Toaster from 'components/toaster';
import { MustAuthController } from 'components/auth';

import PageThemePickerComponent from 'pages/page/theme_picker';
import PageEditRoutesComponent from 'pages/page/edit_route';
import PageButtonsComponent from 'pages/page/buttons';
import PageEditorsComponent from 'pages/page/editors';

export default class PagePage extends MustAuthController {
  page: Page;
  showControls: boolean;
  template: API.ThemeTemplate;
  maximized: boolean;
  minimiseToSettings: boolean;
  _nextRoute: string;
  _clickStart: DOMTokenList; // keep track of click location to prevent firing on drag

  constructor() {
    super();
    this.showControls = false;
    this.maximized = true;
    this.minimiseToSettings = false;
    let pageUUID = m.route.param('id');
    if (pageUUID) {
      Page.get(pageUUID).then((page) => {
        this.updateThemeTemplate(page.theme, page.template)
          .then(() => page.getRoutes())
          .then(() => {
            this.page = page;
            m.redraw();
          });
      });
    } else {
      this.page = new Page();
      this.updateThemeTemplate(this.page.theme, this.page.template);
    }
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

  save() {
    this.page.save().then((page: API.Page) => {
      this.page.uuid = page.uuid;
      window.history.replaceState(null, this.page.title, `/admin/pages/${page.uuid}`);
      return this.page.saveRoutes();
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
    this.page.publish().then(() => {
      Toaster.add('Page published');
      m.redraw();
    });
  }

  unpublish() {
    this.page.unpublish().then(() => {
      Toaster.add('Page unpublished');
      m.redraw();
    });
  }

  delete() {
    this.page.delete().then(() => {
      Toaster.add('Page deleted', 'error');
      m.route.set('/admin/pages');
    });
  }

  maximize() {
    this.minimiseToSettings = true;
    this.maximized = true;
  }

  renderSettings() {
    let path = ['Path: ', <strong>{this.page.defaultRoute}</strong>, ', '];
    return <div class={!this.showControls ? 'controlset' : 'controlset hidden'}>
      <div class='infoset' onclick={() => { this.showControls = true; }}>
        <div class='small black5'>
          {!this.page.defaultRoute ? '' : path}
          Theme: <strong>{this.page.theme}</strong>,
          Template: <strong>{this.page.template}</strong>
        </div>
        <PageButtonsComponent page={this.page} callbacks={this} />
      </div>
    </div>;
  }

  renderSettingsEditor(show: boolean = false) {
    return <div class={this.showControls || show ? 'controlset' : 'controlset hidden'}>
      <div class='settings'>
        <div class='controls'>
          <div class='control'>
            {this.page ? m(PageEditRoutesComponent, this.page) : null}
          </div>
          {this.maximized ? '' :
            <div class='control'
              onclick={() => { this.showControls = false; }}>
              close
            </div>}
        </div>
        <div class='controls'>
          {m(PageThemePickerComponent, {
            theme: this.page.theme,
            template: this.page.template,
            callback: this.updateThemeTemplate.bind(this)
          })}
        </div>
      </div>
      <PageButtonsComponent page={this.page} callbacks={this} />
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
    if (!this.minimiseToSettings) {
      controls.push(
        <a class='typcn typcn-times'
          href='/admin/pages'
          onclick={(e: MouseEvent) => {
            e.preventDefault();
            this._nextRoute = '/admin/pages';
          }}
        />
      );
    }

    let pageMaxClasses = 'page-max animate-fade-in';
    if (this._nextRoute) {
      pageMaxClasses = 'page-max animate-zoom-away';
    }

    return <div class={pageMaxClasses}
      onclick={(e: any) => {
        if (e.target.classList.contains('page-max') && this._clickStart && this._clickStart.contains('page-max')) {
          this._nextRoute = '/admin/pages';
        }
      }}
      oncreate={(v: Mithril.VnodeDOM<any, any>) => {
        v.dom.addEventListener('mousedown', (ev: any) => {
          this._clickStart = ev.target.classList;
        })
        v.dom.addEventListener('animationend', (ev: AnimationEvent) => {
          // old animation is removed, otherwise new animations won't fire.
          if (ev.animationName == 'fadeIn') {
            v.dom.classList.add('animate-fade-in-complete');
            v.dom.classList.remove('animate-fade-in');
          }
          // navigate away after zoomAway animation completes
          if (ev.animationName == 'zoomAway') {
            m.route.set(this._nextRoute);
          }
        });
      }}
    >
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
            value={this.page.title || ''}
            onchange={m.withAttr('value', (v: string) => {
              this.page.title = v;
            })}
          />
        </div>
        <PageEditorsComponent page={this.page} template={this.template} />
      </div>
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
    if (ctrl.maximized) {
      return ctrl.maxView();
    }
    return <div class='page-editor'>
      <div class='page-editor__icons'>
        <span class='typcn typcn-arrow-maximise' onclick={ctrl.maximize} />
      </div>
      <div class='controls'>
        <input
          class='large'
          type='text'
          placeholder='title...'
          oncreate={(v: Mithril.VnodeDOM<any, any>) => {
            (v.dom as HTMLInputElement).value = ctrl.page.title || '';
          }}
          onchange={m.withAttr('value', (v: string) => {
            ctrl.page.title = v;
          })}
        />
      </div>
      {ctrl.renderSettings()}
      {ctrl.renderSettingsEditor()}
      <PageEditorsComponent page={ctrl.page} template={ctrl.template} />
    </div>;
  }
}