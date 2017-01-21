import * as API from 'lib/api';

export default {
  controller: () => {},
  view: (_: any, settings: API.TLSSettingsReponse) => m('.table', [
    m('.tr.tr--center', [
      m('label', 'TLS Email'),
      m('input.large[type=text]', { value: settings.tlsEmail })
    ]),
    m('.tr.tr--center', [
      m('label', 'TLS Domain'),
      m('input.large[type=text]', { value: settings.tlsDomain })
    ]),
    m('.tr.tr--center', [
      m('label', 'TLS Agreement'),
      m('a', { href: settings.termsOfService }, 'link')
    ]),
    m('.tr.tr--center', [
      m('label', 'Agreed On'),
      m('span', settings.agreedOn)
    ]),
    m('.tr.tr--right', [
      m('.button.button--green.button--small', 'Save')
    ])
  ])
}
