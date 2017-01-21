import * as API from 'lib/api';
import Theme from 'lib/theme';
import Layout from 'components/layout';
import { MustAuthController } from 'components/auth';

export default class InstallThemePage extends MustAuthController {
  installedThemes: Mithril.Property<{ [key:string]: boolean }>;
  themes: Mithril.Property<API.Registry>;
  installing: Mithril.Property<string>;

  constructor() {
    super();
    this.themes = m.prop(null);
    this.installing = m.prop<string>();
    this.installedThemes = m.prop({});
    Theme.list().then((themes) => {
      let installed: { [key:string]: boolean } = {};
      themes.forEach((theme) => {
        installed[theme.name] = true;
      });
      this.installedThemes(installed);
    });
    Theme.getAll()
      .then((registry: API.Registry) => this.themes(registry));
  }

  themeInstalled(name: string): boolean {
    return !!this.installedThemes()[name];
  }

  static controller = InstallThemePage;
  static view(ctrl: InstallThemePage) {
    return Layout(m('div', [
      m('h1', 'Theme Manager'),
      !ctrl.installing() ? '' :
        m('div', `Installing theme ${ctrl.installing()}...`),
      m('.table',
        !ctrl.themes() ? null :
          ctrl.themes().packages.map((p: API.Package) => m('.tr', [
            m('div', p.name),
            m('div', p.vcsUrl),
            ctrl.themeInstalled(p.name) ? 'installed' :
              m('a.button.button--small', {
                disabled: !!ctrl.installing(),
                class: !!ctrl.installing() ? 'button--disabled' : 'button--blue',
                onclick: () => {
                  if (ctrl.installing()) {
                    return;
                  }
                  ctrl.installing(p.name);
                  m.redraw();
                  Theme.install(p).then(() => {
                    ctrl.installing(null);
                    m.redraw();
                  });
                }
              }, 'install')
          ]))
      )
    ]));
  }
}