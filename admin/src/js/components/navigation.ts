import { AuthController } from 'components/auth';

export default class NavigationComponent extends AuthController {
  constructor() {
    super();
  }

  static controller = NavigationComponent;
  link(url: string, text: string, additionalClasses: string = '') {
    return m(`a.nav-link${additionalClasses}`, {
      href: url,
      config: m.route
    }, text);
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
    return m('.container--navigation', [
      ctrl.link('/admin', 'ketchup', '.nav-title'),
      ctrl.link('/admin/compose', 'Compose'),
      ctrl.link('/admin/routes', 'Routes'),
      ctrl.link('/admin/pages', 'Pages'),
      ctrl.link('/admin/themes', 'Theme'),
      ctrl.link('/admin/settings', 'Settings'),
      ctrl.link('/admin/logout', 'Logout')
    ]);
  };
}

