import {Loader} from 'components/loading';
import {ModalComponent} from 'components/modal';
import {LinkRow, Row} from 'components/table';
import * as Toaster from 'components/toaster';
import * as API from 'lib/api';
import Theme from 'lib/theme';
import * as React from 'react';
import Layout from 'components/layout';

interface State {
  installedThemes: {[key: string]: boolean};
  themes?: API.Registry;
  installing?: string;
}

export default class InstallThemePage extends React.Component<{}, State> {
  constructor(props: any) {
    super(props);
    this.state = {
      installedThemes: {},
    };
  }

  componentDidMount() {
    Theme.getAll().then((registry: API.Registry) => this.setState({themes: registry}));
    this.loadThemes();
  }

  loadThemes() {
    Theme.list().then((themes) => {
      let installed: {[key: string]: boolean} = {};
      themes.forEach((theme) => {
        installed[theme.name] = true;
      });
      this.setState({installedThemes: installed});
    });
  }

  themeInstalled(name: string): boolean {
    return !!this.state.installedThemes[name];
  }

  installTheme(p: API.Package) {
    if (this.state.installing) {
      return;
    }
    this.setState({installing: p.name});
    Theme.install(p).then(() => {
      Toaster.add(this.state.installing + ' theme installed.');
      this.setState({installing: null});
      return this.loadThemes();
    });
  }

  render() {
    const packages =
      this.state.themes && this.state.themes.packages ? this.state.themes.packages : [];
    const themes = packages.map((p: API.Package) => {
      if (this.themeInstalled(p.name)) {
        return (
          <LinkRow key={p.name} href={`/themes/${p.name}`} link>
            <div>{p.name}</div>
            <div>{p.vcsUrl}</div>
            <div>installed</div>
          </LinkRow>
        );
      }

      let linkClasses = 'button button--small button--blue';
      if (!!this.state.installing) {
        linkClasses += ' button--disabled';
      }
      return (
        <Row key={p.name} center>
          <div>{p.name}</div>
          <div>{p.vcsUrl}</div>
          <a
            className={linkClasses}
            onClick={(e) => {
              e.preventDefault();
              this.installTheme(p);
            }}
          >
            install
          </a>
        </Row>
      );
    });

    return (
      <Layout className="install-theme">
        <h1>Theme Manager</h1>
        <div className="table">{themes}</div>
        <ModalComponent
          title="Installing theme..."
          visible={!!this.state.installing}
          toggle={() => {}}
        >
          {Loader}
        </ModalComponent>
      </Layout>
    );
  }
}
