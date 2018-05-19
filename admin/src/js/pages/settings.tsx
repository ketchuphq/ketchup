import Layout from 'components/layout';
import {Loader, LoadingTable} from 'components/loading';
import * as API from 'lib/api';
import {get} from 'lib/requests';
import TLSComponent from 'pages/settings/tls';
import TLSNewComponent from 'pages/settings/tls-new';
import * as React from 'react';
import {PrivateRouteComponentProps} from 'components/auth';

// add redirect
// setup static upload
interface State {
  settings?: API.TLSSettingsResponse;
  version?: string;
  registryURL?: string;
}

export default class SettingsPage extends React.Component<PrivateRouteComponentProps<any>, State> {
  constructor(props: any) {
    super(props);
    this.state = {};
  }

  componentDidMount() {
    get('/api/v1/settings/tls')
      .then((res) => res.json())
      .then((settings: API.TLSSettingsResponse) => {
        this.setState({settings});
      });

    get('/api/v1/settings/info')
      .then((res) => res.json())
      .then(({version, registry_url}) => {
        this.setState({
          version,
          registryURL: registry_url,
        });
      });
  }

  render() {
    const {settings} = this.state;
    let tlsSection;
    if (!settings) {
      tlsSection = <LoadingTable />;
    } else if (Object.keys(settings).length == 0 || !settings.hasCertificate) {
      tlsSection = (
        <TLSNewComponent
          user={this.props.user}
          email={settings.tlsEmail}
          domain={settings.tlsDomain}
        />
      );
    } else {
      tlsSection = <TLSComponent {...settings} />;
    }
    return (
      <Layout className="settings">
        <header>
          <h1>Settings</h1>
        </header>
        <h2>Ketchup</h2>
        <div className="table">
          <div className="tr tr--center">
            <label>Version</label>
            <div>{this.state.version}</div>
          </div>
          <div className="tr tr--center">
            <label>Theme Registry</label>
            <div>{this.state.registryURL}</div>
          </div>
          <div className="tr tr--center">
            <label>Export your data as JSON</label>
            <a className="button button--green button--small" href="/api/v1/download-backup">
              Download backup
            </a>
          </div>
        </div>
        <h2>TLS</h2>
        {tlsSection}
      </Layout>
    );
  }
}
