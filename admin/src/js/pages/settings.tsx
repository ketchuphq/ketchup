import msx from 'lib/msx';
import * as m from 'mithril';
import * as API from 'lib/api';
import { MustAuthController } from 'components/auth';
import TLSNewComponent from 'pages/settings/tls-new';
import TLSComponent from 'pages/settings/tls';

// add redirect
// setup static upload
export default class SettingsPage extends MustAuthController {
  settings: API.TLSSettingsReponse;
  version: string;
  registryURL: string;

  constructor() {
    super();
    m.request({
      method: 'GET',
      url: '/api/v1/settings/tls',
    }).then((settings: API.TLSSettingsReponse) => {
      this.settings = settings;
      m.redraw();
    });
    m.request({
      method: 'GET',
      url: '/api/v1/settings/info',
    }).then(({ version, registry_url }) => {
      this.version = version;
      this.registryURL = registry_url;
      m.redraw();
    });
  }

  view() {
    let settings = this.settings;
    let tlsSection;
    if (!settings) {
      tlsSection = <div>loading...</div>;
    } else if (Object.keys(settings).length == 0 || !settings.hasCertificate) {
      tlsSection = <TLSNewComponent
        user={this.user}
        email={settings.tlsEmail}
        domain={settings.tlsDomain}
      />;
    } else {
      tlsSection = <TLSComponent {...settings} />;
    }
    return <div class='settings'>
      <header>
        <h1>Settings</h1>
      </header>
      <h2>Ketchup</h2>
      <div class='table'>
        <div class='tr tr--center'>
          <label>Version</label>
          <div>{this.version}</div>
        </div>
        <div class='tr tr--center'>
          <label>Theme Registry</label>
          <div>{this.registryURL}</div>
        </div>
        <div class='tr tr--center'>
          <label>Export your data as JSON</label>
          <a
            class='button button--green button--small'
            href='/api/v1/download-backup'
          >
            Download backup
          </a>
        </div>
      </div>
      <h2>TLS</h2>
      {tlsSection}
    </div>;
  }
}
