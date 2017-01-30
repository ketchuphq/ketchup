import Theme from 'lib/theme';
import Layout from 'components/layout';
import { MustAuthController } from 'components/auth';

export default class ThemePage extends MustAuthController {
  theme: Mithril.Property<Theme>;

  constructor() {
    super();
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
    if (!ctrl.theme()) {
      return Layout('');
    }
    let templateKeys = Object.keys(ctrl.theme().templates);
    let assetKeys = Object.keys(ctrl.theme().assets);
    return Layout(
      m('.theme', [
        m('header',
          m('h1', [
            m('a[href=/admin/themes]', { config: m.route }, 'Themes'),
            m.trust(' &rsaquo; '),
            m('span.unbold', ctrl.theme().name)
          ]),
        ),
        m('h2', 'Templates'),
        m('.table',
          templateKeys.sort().map((name) => {
            let t = ctrl.theme().templates[name];
            return m('a.tr', {
              config: m.route,
              href: `/admin/themes/${ctrl.theme().name}/templates/${t.name}`
            }, [
                m('div', t.name),
                m('div', t.engine)
              ]);
          })
        ),
        m('h2', 'Assets'),
        m('.table',
          assetKeys.length == 0 ?
            m('.tr', 'no assets') :
            assetKeys.sort().map((asset) => {
              // todo: check for shadowed files
              return m('a.tr', {
                href: `/${asset}`
              }, asset);
            })
        )
      ])
    );
  }
}