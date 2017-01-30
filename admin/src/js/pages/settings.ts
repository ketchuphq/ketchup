import * as API from 'lib/api';
import Layout from 'components/layout';
import { MustAuthController } from 'components/auth';
import TLSNewComponent from 'pages/settings/tls-new';
import TLSComponent from 'pages/settings/tls';

// add redirect
// setup static upload
export default class SettingsPage extends MustAuthController {
  settings: API.TLSSettingsReponse;

  constructor() {
    super();
    this.settings = false;
    m.startComputation();
    m.request({
      method: 'GET',
      url: '/api/v1/settings/tls',
    }).then((settings: API.TLSSettingsReponse) => {
      this.settings = settings;
      m.endComputation();
    });
  }

  static controller = SettingsPage;
  static view(ctrl: SettingsPage) {
    let tlsSection = null;
    if (ctrl.settings === false) {
      tlsSection = m('div', 'loading...');
    } else if (Object.keys(ctrl.settings).length == 0 || !ctrl.settings.hasCertificate) {
      tlsSection = m.component(TLSNewComponent, ctrl.user(), ctrl.settings.tlsEmail, ctrl.settings.tlsDomain);
    } else {
      tlsSection = m.component(TLSComponent, ctrl.settings);
    }

    return Layout(
      m('.settings', [
        m('header', m('h1', 'Settings')),
        // m('.table',
        //   m('.tr', [
        //     m('label', 'Ketchup Version'),
        //     m('div', '0.1')
        //   ]),
        //   m('.tr', [
        //     m('label', 'Default theme and template'),
        //     m('input.large', { type: 'text' })
        //   ])
        // ),
        m('h2', 'TLS'),
        tlsSection
      ])
    );
  }
}