import msx from 'lib/msx';
import * as m from 'mithril';
import * as API from 'lib/api';
import Theme from 'lib/theme';
import { MustAuthController } from 'components/auth';

export default class InstallThemePage extends MustAuthController {
  installedThemes: { [key: string]: boolean };
  themes: API.Registry;
  installing: string;

  constructor() {
    super();
    this.installedThemes = {};
    Theme.list().then((themes) => {
      let installed: { [key: string]: boolean } = {};
      themes.forEach((theme) => {
        installed[theme.name] = true;
      });
      this.installedThemes = installed;
    });
    Theme.getAll().then((registry: API.Registry) => (this.themes = registry));
  }

  themeInstalled(name: string): boolean {
    return !!this.installedThemes[name];
  }

  view() {
    let themeInstall = (p: API.Package) => {
      return () => {
        if (this.installing) {
          return;
        }
        this.installing = p.name;
        m.redraw();
        Theme.install(p).then(() => {
          this.installing = null;
          m.redraw();
        });
      };
    };
    let packages = this.themes && this.themes.packages ? this.themes.packages : [];
    let themes = packages.map((p: API.Package) => {
      let themeControl;
      if (this.themeInstalled(p.name)) {
        themeControl = (
          <a class='tr' href={`/admin/themes/${p.name}`} oncreate={m.route.link}>
            installed
          </a>
        );
      } else {
        let themeControlClasses = ['button', 'button--small'];
        if (!!this.installedThemes) {
          themeControlClasses.push('button--disabled');
        } else {
          themeControlClasses.push('button--blue');
        }
        themeControl = (
          <a
            disabled={!!this.installing}
            class={themeControlClasses.join(' ')}
            onclick={themeInstall(p)}
          >
            install
          </a>
        );
        return (
          <div class='tr'>
            <div>{p.name}</div>
            <div>{p.vcsUrl}</div>
            {themeControl}
          </div>
        );
      }
    });

    return (
      <div>
        <h1>Theme Manager</h1>
        {!this.installing ? '' : <div>{`Installing theme ${this.installing}...`}</div>}
        <div class='table'>{themes}</div>
      </div>
    );
  }
}
