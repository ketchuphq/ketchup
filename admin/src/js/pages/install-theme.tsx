import msx from 'lib/msx';
import * as m from 'mithril';
import * as API from 'lib/api';
import Theme from 'lib/theme';
import { MustAuthController } from 'components/auth';

export default class InstallThemePage extends MustAuthController {
  installedThemes: { [key:string]: boolean };
  themes: API.Registry;
  installing: string;

  constructor() {
    super();
    this.installedThemes = {};
    Theme.list().then((themes) => {
      let installed: { [key:string]: boolean } = {};
      themes.forEach((theme) => {
        installed[theme.name] = true;
      });
      this.installedThemes = installed;
    });
    Theme.getAll()
      .then((registry: API.Registry) => this.themes = registry);
  }

  themeInstalled(name: string): boolean {
    return !!this.installedThemes[name];
  }

  static oninit(v: Mithril.Vnode<{}, InstallThemePage>) {
    v.state = new InstallThemePage();
  }

  static view(v: Mithril.Vnode<{}, InstallThemePage>) {
    let ctrl = v.state;
    let themes;
    let themeInstall = (p: API.Package) => {
      return () => {
        if (ctrl.installing) {
          return;
        }
        ctrl.installing = p.name;
        m.redraw();
        Theme.install(p).then(() => {
          ctrl.installing = null;
          m.redraw();
        });
      };
    };
    if (ctrl.themes) {
      themes = ctrl.themes.packages.map((p: API.Package) =>
        <div class='tr'>
          <div>{p.name}</div>
          <div>{p.vcsUrl}</div>
          {
          ctrl.themeInstalled(p.name) ? 'installed' :
            <a disabled={!!ctrl.installing}
              class={'button button--small' + (!!ctrl.installing ? 'button--disabled' : 'button--blue')}
              onclick={themeInstall(p)}
            >
              install
            </a>
          }
        </div>
      );
    }

    return <div>
      <h1>Theme Manager</h1>
      {!ctrl.installing ? '' : <div>{`Installing theme ${ctrl.installing}...`}</div>}
      <div class='table'>{themes}</div>
    </div>;
  }
}