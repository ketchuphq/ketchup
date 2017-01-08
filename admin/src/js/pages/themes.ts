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
      m('h1', 'Theme'),
      m('table',
        ctrl.themes().map((theme) => {
          return m('tr',
            m('td.link-cell',
              m('a', {
                href: `/admin/themes/${theme.name}`,
                config: m.route
              }, theme.name || 'untitled')
            )
          );
        })
      )
    ]));
  }
}