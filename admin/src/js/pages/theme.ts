import Theme from 'lib/theme';
import Layout from 'components/layout';

export default class ThemePage {
  theme: Mithril.Property<Theme>;

  constructor() {
    this.theme = m.prop<Theme>();
    let themeName = m.route.param('name');
    if (themeName) {
      Theme.get(themeName).then((theme) => {
        this.theme(theme);
      });
    }
  }

  static controller = ThemePage;
  static view(ctrl: ThemePage) {
    let templateKeys = Object.keys(ctrl.theme().templates);
    let assetKeys = Object.keys(ctrl.theme().assets);
    return Layout(
      !ctrl.theme() ? ''
        :
        m('.theme', [
          m('h1',
            m.trust('Theme &rsaquo; '),
            ctrl.theme().name
          ),
          m('h2', 'Templates'),
          m('table',
            templateKeys.sort().map((name) => {
              let t = ctrl.theme().templates[name];
              return m('tr', [
                m('td', t.name),
                m('td', t.engine)
              ]);
            })
          ),
          m('h2', 'Assets'),
          m('table',
            assetKeys.length == 0 ?
              m('tr', m('td', 'no assets'))
              :
              assetKeys.sort().map((asset) => {
                return m('tr', [
                  m('td.link-cell',
                    // todo: check for shadowed files
                    m('a', {
                      href: `/${asset}`
                    }, asset)
                  ),
                ]);
              })
          )
        ])
    );
  }
}