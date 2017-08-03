import msx from 'lib/msx';
import * as m from 'mithril';
import * as API from 'lib/api';
import Theme from 'lib/theme';
import { MustAuthController } from 'components/auth';
import { LinkRow, Row } from 'components/table';
import { ModalComponent } from 'components/modal';
import { Loader } from 'components/loading';
import * as Toaster from 'components/toaster';

export default class InstallThemePage extends MustAuthController {
  installedThemes: { [key: string]: boolean };
  themes: API.Registry;
  installing: string;

  constructor() {
    super();
    this.installedThemes = {};
    this.loadThemes();
    Theme.getAll().then((registry: API.Registry) => (this.themes = registry));
  }

  loadThemes() {
    Theme.list().then((themes) => {
      let installed: { [key: string]: boolean } = {};
      themes.forEach((theme) => {
        installed[theme.name] = true;
      });
      this.installedThemes = installed;
      m.redraw();
    });
  }

  themeInstalled(name: string): boolean {
    return !!this.installedThemes[name];
  }

  installTheme(p: API.Package) {
    if (this.installing) {
      return;
    }
    this.installing = p.name;
    Theme.install(p).then(() => {
      Toaster.add(this.installing + ' theme installed.');
      this.installing = null;
      return this.loadThemes();
    });
  }

  view() {
    let packages = this.themes && this.themes.packages ? this.themes.packages : [];
    let themes = packages.map((p: API.Package) => {
      if (this.themeInstalled(p.name)) {
        return (
          <LinkRow href={`/admin/themes/${p.name}`} link>
            <div>{p.name}</div>
            <div>{p.vcsUrl}</div>
            <div>installed</div>
          </LinkRow>
        );
      }

      return (
        <Row center>
          <div>{p.name}</div>
          <div>{p.vcsUrl}</div>
          <a
            class='button button--small button--blue'
            disabled={!!this.installing}
            onclick={() => this.installTheme(p)}
          >
            install
          </a>
        </Row>
      );
    });

    return (
      <div>
        <h1>Theme Manager</h1>
        <div class='table'>{themes}</div>
        <ModalComponent
          title='Installing theme...'
          visible={() => !!this.installing}
          toggle={() => {}}
        >
          {Loader}
        </ModalComponent>
      </div>
    );
  }
}
