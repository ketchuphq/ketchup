export default class NavigationComponent {
  constructor() {
  }

  static controller = NavigationComponent;
  static view(ctrl: NavigationComponent) {
    return m('.container--navigation', [
      m('.nav-title', 'ketchup'),
      m('a.nav-link', {
        href: '/admin/compose',
        config: m.route
      }, 'Compose'),
      m('a.nav-link', {
        href: '/admin/routes',
        config: m.route
      }, 'Routes'),
      m('a.nav-link', {
        href: '/admin/pages',
        config: m.route
      }, 'Pages')
    ]);
  };
}

