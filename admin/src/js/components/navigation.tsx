import msx from 'lib/msx';
import { AuthController } from 'components/auth';
import Button from 'components/button';

export default class NavigationComponent extends AuthController {
  collapsed: Mithril.Property<boolean>;
  constructor() {
    super();
    this.collapsed = m.prop(this._userPromise.then(() =>
      this.pref('hideMenu') || false
    )) as Mithril.Property<boolean>;
  }

  toggle() {
    let collapsed = !this.collapsed();
    this.setPref('hideMenu', collapsed);
    this.collapsed(collapsed);
  }

  static controller = NavigationComponent;
  link(url: string, text: string, opts: { onclick?: () => void, additionalClasses?: string, icon?: string } = {}) {
    return m(`a.nav-link${opts.additionalClasses || ''}`, {
      href: url,
      config: m.route,
      onclick: opts.onclick
    }, [
        !!opts.icon ? <span class={`typcn typcn-${opts.icon}`}></span> : '',
        <span class='nav-link__text'>{text}</span>
      ]);
  }

  static view(ctrl: NavigationComponent) {
    if (!ctrl.user()) {
      return m('.navigation', [
        ctrl.link('/admin', 'K', { additionalClasses: '.nav-title' }),
        ctrl.link('/admin/login', 'Login')
      ]);
    }

    let navClass = 'navigation';
    navClass += ctrl.collapsed() ? ' navigation--hidden' : '';
    return <div class={navClass}>
      {ctrl.link('/admin', 'K', { additionalClasses: '.nav-title' })}
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
