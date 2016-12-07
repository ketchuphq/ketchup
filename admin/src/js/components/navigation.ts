import { AuthController } from 'components/auth';

let collapsed = false;

export default class NavigationComponent extends AuthController {
  collapsed: Mithril.Property<boolean>;
  constructor() {
    super();
    this.collapsed = m.prop(collapsed);
  }

  toggle() {
    collapsed = !this.collapsed();
    this.collapsed(collapsed);
  }

  static controller = NavigationComponent;
  link(url: string, text: string, opts: { additionalClasses?: string, icon?: string } = {}) {
    return m(`a.nav-link${opts.additionalClasses || ''}`, {
      href: url,
      config: m.route
    }, [
        opts.icon != '' ? m(`span.typcn.typcn-${opts.icon}`) : '',
        m('span.nav-link__text', text)
      ]);
  }

  static view(ctrl: NavigationComponent) {
    if (!ctrl.user()) {
      return m('.container--navigation', [
        m('a.nav-title', {
          href: '/admin',
          config: m.route
        }, 'ketchup'),
        ctrl.link('/admin/login', 'Login')
      ]);
    }
    return m('.container--navigation', {
      className: ctrl.collapsed() ? 'container--navigation--hidden' : '',
    }, [
        ctrl.link('/admin', 'ketchup', { additionalClasses: '.nav-title' }),
        m('.nav-button',
          m('a.button.button--green.button--center', {
            href: '/admin/compose',
            config: m.route
          },
            m('span.typcn.typcn-edit'),
            m('span.nav-link__text', 'Compose')
          )
        ),
        ctrl.link('/admin/routes', 'Routes', { icon: 'flow-children' }),
        ctrl.link('/admin/pages', 'Pages', { icon: 'document-text' }),
        ctrl.link('/admin/themes', 'Theme', { icon: 'brush' }),
        ctrl.link('/admin/settings', 'Settings', { icon: 'spanner-outline' }),
        ctrl.link('/admin/logout', 'Logout', { icon: 'weather-night' }),
        m(`a.nav-link.nav-link--toggle`, {
          onclick: () => { ctrl.toggle(); }
        },
          m(`span.typcn.typcn-arrow-${ctrl.collapsed() ? 'maximise' : 'minimise'}`),
        )
      ]);
  };
}

