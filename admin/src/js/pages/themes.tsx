import msx from 'lib/msx';
import * as m from 'mithril';
import Theme from 'lib/theme';
import { MustAuthController } from 'components/auth';
import { loading } from 'components/loading';

export default class ThemesPage extends MustAuthController {
  themes: Theme[];
  loading: boolean;

  constructor() {
    super();
    this.themes = [];
    this.loading = true;
    Theme.list().then((themes) => {
      this.loading = false;
      this.themes = themes;
    });
  }

  view() {
    return <div class='themes'>
      <header>
        <a class='button button--green button--center'
          href='/admin/themes-install'
          oncreate={m.route.link}
        >
          Get More
        </a>
        <h1>Themes</h1>
      </header>

      <h2>Installed themes</h2>
      <div class='table'>
        {loading(this.loading)}
        {this.themes.map((theme) => {
          return <a class='tr'
            href={`/admin/themes/${theme.name}`}
            oncreate={m.route.link}
          >
            {theme.name || 'untitled'}
          </a>;
        })}
      </div>
    </div>;
  }
}
