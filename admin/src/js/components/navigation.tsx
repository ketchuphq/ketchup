import msx from 'lib/msx';
import { AuthController } from 'components/auth';
import Button from 'components/button';
// https://github.com/DefinitelyTyped/DefinitelyTyped/issues/13679
var debounce: any = require('lodash/debounce');

let retain = (_: any, __: any, context: Mithril.Context) => context.retain = true;


export default class NavigationComponent extends AuthController {
  collapsed: Mithril.Property<boolean>;
  constructor() {
    super();
    this.collapsed = m.prop(this._userPromise.then(() =>
      this.pref('hideMenu') || window.innerWidth <= 480
    )) as Mithril.Property<boolean>;

    window.addEventListener('resize', this.resizeHandler);
    this.resizeHandler();
  }

  onunload() {
    window.removeEventListener('resize', this.resizeHandler);
  }

  resizeHandler = debounce(() => {
    if (window.innerWidth > 480) {
      return;
    }
    m.startComputation();
    this.setPref('hideMenu', true);
    this.collapsed(true);
    m.endComputation();
  }, 300);

  toggle() {
    let collapsed = !this.collapsed();
    this.setPref('hideMenu', collapsed);
    this.collapsed(collapsed);
  }

  static controller = NavigationComponent;
  link(url: string, text: string, opts: { onclick?: () => void, additionalClasses?: string, icon?: string } = {}) {
    return m('a.nav-link', {
      class: opts.additionalClasses || '',
      href: url,
      config: (a, b, c, d) => {
        c.retain = true;
        return m.route(a, b, c, d);
      },
      onclick: opts.onclick,
      key: `nav-${text.toLowerCase()}`,
    }, [
        !!opts.icon ? <span class={`typcn typcn-${opts.icon}`}></span> : '',
        <span class='nav-link__text'>{text}</span>
      ]);
  }

  static view(ctrl: NavigationComponent) {
    let navClass = 'navigation';
    if (!ctrl.user()) {
      return <div class={navClass} key='navigation' config={retain}>
        {ctrl.link('/admin', 'K', { additionalClasses: 'nav-title' })}
        {ctrl.link('/admin/login', 'Login')}
      </div>;
    }

    if (ctrl.collapsed()) {
      navClass += ' navigation--hidden';
    }

    return <div class={navClass} key='navigation' config={retain}>
      {ctrl.link('/admin', 'K', { additionalClasses: 'nav-title' })}
      <div class='nav-button'>
        <Button
          class='button--green button--center'
          href='/admin/compose'
          config={m.route}>
          <span class='typcn typcn-edit' />
          <span class='nav-link__text'>Compose</span>
        </Button>
      </div>
      {ctrl.link('/admin/pages', 'Pages', { icon: 'document-text' })}
      {ctrl.link('/admin/themes', 'Theme', { icon: 'brush' })}
      {ctrl.link('/admin/settings', 'Settings', { icon: 'spanner-outline' })}
      {ctrl.link('/admin/logout', 'Logout', { onclick: () => ctrl.logout(), icon: 'weather-night' })}
      <a class='nav-link nav-link--toggle' onclick={() => ctrl.toggle()}>
        <span class={`typcn typcn-arrow-${ctrl.collapsed() ? 'maximise' : 'minimise'}`} />
      </a>
    </div>;
  }
}
