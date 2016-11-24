import Theme from 'lib/theme';
import Layout from 'components/layout';

export default class ThemePage {
  themes: Mithril.Property<Theme[]>;

  constructor() {
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