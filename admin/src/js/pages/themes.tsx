import msx from 'lib/msx';
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
    return Layout(<div class='themes'>
      <header>
        <a class='button button--green button--center'
          href='/admin/themes/install'
          config={m.route}
        >
          Get More
        </a>
        <h1>Themes</h1>
      </header>

      <h2>Installed themes</h2>
      <div class='table'>
        {ctrl.themes().map((theme) => {
          return <a class='tr'
            href={`/admin/themes/${theme.name}`}
            config={m.route}
          >
            {theme.name || 'untitled'}
          </a>;
        })}
      </div>
    </div>);
  }
}