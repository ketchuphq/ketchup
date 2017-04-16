import msx from 'lib/msx';
import * as m from 'mithril';
import * as API from 'lib/api';
import { MustAuthController } from 'components/auth';
import TLSNewComponent from 'pages/settings/tls-new';
import TLSComponent from 'pages/settings/tls';

let _: Mithril.Component<{}, SettingsPage> = SettingsPage;

// add redirect
// setup static upload
export default class SettingsPage extends MustAuthController {
  settings: API.TLSSettingsReponse;
  version: string;

  constructor() {
    super();
    this.settings = false;
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
    }).then(({ version }) => {
      this.version = version;
      m.redraw();
    });
  }

  static oninit(v: Mithril.Vnode<{}, SettingsPage>) {
    v.state = new SettingsPage();
  };

  static view(v: Mithril.Vnode<{}, SettingsPage>) {
    let settings = v.state.settings;
    let tlsSection = null;
    if (settings === false) {
      tlsSection = m('div', 'loading...');
    } else if (Object.keys(settings).length == 0 || !settings.hasCertificate) {
      tlsSection = m(TLSNewComponent, {
        user: v.state.user,
        email: settings.tlsEmail,
        domain: settings.tlsDomain
      });
    } else {
      tlsSection = m(TLSComponent, settings);
    }
    return <div class='settings'>
      <header>
        <h1>Settings</h1>
      </header>
      <h2>Ketchup</h2>
      <div class='table'>
        <div class='tr tr--center'>
          <label>Version</label>
          <div>{v.state.version}</div>
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