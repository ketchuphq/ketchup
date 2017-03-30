import msx from 'lib/msx';
import * as m from 'mithril';
import Theme from 'lib/theme';
import { MustAuthController } from 'components/auth';
import { loading } from 'components/loading';

let _: Mithril.Component<{}, ThemesPage> = ThemesPage;

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

  static oninit(v: Mithril.Vnode<{}, ThemesPage>) {
    v.state = new ThemesPage();
  };

  static view(v: Mithril.Vnode<{}, ThemesPage>) {
    let ctrl = v.state;
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
        {loading(ctrl.loading)}
        {ctrl.themes.map((theme) => {
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