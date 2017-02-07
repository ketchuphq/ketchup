import msx from 'lib/msx';
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

    let assetsList = [<a class='tr'>no assets</a>];
    if (assetKeys.length > 0) {
      assetsList = assetKeys.sort().map((asset) =>
        <a class='tr' href={`/${asset}`}>{asset}</a>
      );
    }

    return Layout(
      <div class='theme'>
        <header>
          <h1>
            <a href='/admin/themes' config={m.route}>Themes</a>
            {m.trust(' &rsaquo; ')}
            <span class='unbold'>{ctrl.theme().name}</span>
          </h1>
        </header>

        <h2>Templates</h2>
        <div class='table'>
          {templateKeys.sort().map((name) => {
            let t = ctrl.theme().templates[name];
            return <a class='tr'
              config={m.route}
              href={`/admin/themes/${ctrl.theme().name}/templates/${t.name}`}
            >
              <div>{t.name}</div>
              <div>{t.engine}</div>
            </a>;
          })}
        </div>

        <h2>Assets</h2>
        <div class='table'>
          {assetsList}
        </div>
      </div>
    );
  }
}