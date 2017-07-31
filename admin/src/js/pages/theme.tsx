import msx from 'lib/msx';
import * as m from 'mithril';
import Theme from 'lib/theme';
import { MustAuthController } from 'components/auth';

export default class ThemePage extends MustAuthController {
  theme: Theme;

  constructor() {
    super();
    let themeName = m.route.param('name');
    if (themeName) {
      Theme.get(themeName).then((theme) => {
        this.theme = theme;
        m.redraw();
      });
    }
  }

  view() {
    if (!this.theme) {
      return;
    }
    let templateKeys = Object.keys(this.theme.templates);
    let assetKeys = Object.keys(this.theme.assets);

    let assetsList = [<a class='tr'>no assets</a>];
    if (assetKeys.length > 0) {
      assetsList = assetKeys.sort().map((asset) =>
        <a class='tr' href={`/${asset}`}>{asset}</a>
      );
    }

    return <div class='theme'>
      <header>
        <h1>
          <a href='/admin/themes' oncreate={m.route.link}>
            Themes
          </a>
          {m.trust(' &rsaquo; ')}
          <span class='unbold'>{this.theme.name}</span>
        </h1>
      </header>

      <h2>Templates</h2>
      <div class='table'>
        {templateKeys.sort().map((name) => {
          let t = this.theme.templates[name];
          return <a class='tr'
            oncreate={m.route.link}
            href={`/admin/themes/${this.theme.name}/templates/${t.name}`}
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
    </div>;
  }
}
