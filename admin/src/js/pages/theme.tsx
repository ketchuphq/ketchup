import Button from 'components/button';
import Layout from 'components/layout';
import {Loader} from 'components/loading';
import {ConfirmModalComponent} from 'components/modal';
import {LinkRow, Row, Table} from 'components/table';
import * as API from 'lib/api';
import {Package} from 'lib/api';
import Theme from 'lib/theme';
import * as React from 'react';
import {RouteComponentProps} from 'react-router';
import {Link} from 'react-router-dom';

interface ThemeProps {
  theme: API.Theme;
  ref?: string;
}

interface VersionState {
  showUpdateModal: boolean;
  hasUpdate: boolean;
  latestRef?: string;
}

class VersionSection extends React.Component<ThemeProps, VersionState> {
  constructor(v: ThemeProps) {
    super(v);
    this.state = {
      showUpdateModal: false,
      hasUpdate: false,
    };
  }

  checkUpdates() {
    return Theme.checkForUpdates(this.props.theme.name).then(({currentRef}) => {
      if (!currentRef) {
        this.setState({
          latestRef: 'No updates found.',
          hasUpdate: false,
          showUpdateModal: true,
        });
      } else {
        this.setState({
          latestRef: `Latest ref: ${currentRef.slice(0, 6)}`,
          hasUpdate: true,
          showUpdateModal: true,
        });
      }
    });
  }

  render() {
    if (!this.props.ref) {
      return null;
    }
    return (
      <div>
        <h2>Version</h2>
        <div className="table">
          <div className="tr tr--center">
            <code>{this.props.ref}</code>
            <Button className="button--green button--small" handler={() => this.checkUpdates()}>
              Check for Updates
            </Button>
          </div>
        </div>
        <ConfirmModalComponent
          title="Updates"
          visible={this.state.showUpdateModal}
          toggle={() => {
            this.setState((state) => ({showUpdateModal: !state.showUpdateModal}));
          }}
          confirmText={this.state.hasUpdate ? 'Update' : 'Okay'}
          hideCancel={this.state.hasUpdate}
        >
          <p>{this.state.latestRef}</p>
        </ConfirmModalComponent>
      </div>
    );
  }
}

class TemplatesSection extends React.PureComponent<ThemeProps> {
  render() {
    let theme = this.props.theme;
    let templateKeys = Object.keys(theme.templates);
    let templates = templateKeys
      .sort()
      .map((name) => theme.templates[name])
      .map((t) => (
        <LinkRow key={t.name} href={`/themes/${theme.name}/templates/${t.name}`} link>
          <div>{t.name}</div>
          <div>{t.engine}</div>
        </LinkRow>
      ));

    return (
      <div>
        <h2>Templates</h2> <div className="table">{templates}</div>
      </div>
    );
  }
}

class AssetsSection extends React.PureComponent<ThemeProps> {
  render() {
    let assetKeys = Object.keys(this.props.theme.assets);
    let assetsList: React.ReactNode = <a className="tr">no assets</a>;
    if (assetKeys.length > 0) {
      assetsList = assetKeys.sort().map((asset) => (
        <LinkRow key={asset} href={`/${asset}`}>
          {asset}
        </LinkRow>
      ));
    }
    return (
      <div>
        <h2>Assets</h2>
        <Table>{assetsList}</Table>
      </div>
    );
  }
}

interface PackageProps {
  pkg: Package;
}

class PackageSection extends React.PureComponent<PackageProps> {
  render() {
    if (!this.props.pkg) {
      return null;
    }
    let pkg = this.props.pkg;
    let fields = [];
    if (pkg.authors && pkg.authors.length > 0) {
      fields.push(
        <Row>
          <div>Authors</div>
          <div>{pkg.authors.join(', ')}</div>
        </Row>
      );
    }

    if (pkg.homepage) {
      fields.push(
        <Row>
          <div>Homepage</div>
          <div>{pkg.homepage}</div>
        </Row>
      );
    }
    if (pkg.vcsUrl) {
      fields.push(
        <Row>
          <div>Source</div>
          <div>{pkg.vcsUrl}</div>
        </Row>
      );
    }
    if (pkg.tags && pkg.tags.length > 0) {
      fields.push(
        <Row>
          <div>Tags</div>
          <div>{pkg.tags.join(', ')}</div>
        </Row>
      );
    }

    if (fields.length == 0) {
      return null;
    }

    return (
      <div>
        <h2>Package</h2>
        <Table>{fields}</Table>
      </div>
    );
  }
}

interface State {
  theme?: API.Theme;
}

export default class ThemePage extends React.Component<RouteComponentProps<{id: string}>, State> {
  constructor(props: any) {
    super(props);
    this.state = {};
  }

  componentDidMount() {
    Theme.get(this.props.match.params.id).then((theme) => {
      this.setState({theme});
    });
  }

  render() {
    if (!this.state.theme) {
      return (
        <Layout className="theme">
          <header>
            <h1>
              <Link to="/themes">Themes</Link> &rsaquo;{' '}
            </h1>
          </header>
          <Loader show />;
        </Layout>
      );
    }

    const theme = this.state.theme;

    return (
      <Layout className="theme">
        <header>
          <h1>
            <Link to="/themes">Themes</Link> &rsaquo; <span className="unbold">{theme.name}</span>
          </h1>
        </header>
        <p className="txt-large">{theme.description}</p>
        <VersionSection theme={theme} />
        <PackageSection pkg={theme.package} />
        <TemplatesSection theme={theme} />
        <AssetsSection theme={theme} />
      </Layout>
    );
  }
}
