import Theme from 'lib/theme';
import Layout from 'components/layout';
import { MustAuthController } from 'components/auth';

export default class ThemePage extends MustAuthController {
  themes: Mithril.Property<Theme[]>;

  constructor() {
    super();
    this.themes = m.prop([]);
    Theme.list().then((themes) => this.themes(themes));
  }

  static controller = ThemePage;
  static view(ctrl: ThemePage) {
    return Layout(m('.themes', [
      m('header',
        m('a.button.button--green.button--center', {
          href: '/admin/themes/install',
          config: m.route
        }, 'Get More'),
        m('h1', 'Themes')
      ),
      m('h2', 'Installed themes'),
      m('.table',
        ctrl.themes().map((theme) => {
          return m('a.tr', {
            href: `/admin/themes/${theme.name}`,
            config: m.route
          }, theme.name || 'untitled');
        })
      )
    ]));
  }
}