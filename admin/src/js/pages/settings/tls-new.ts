import * as API from 'lib/api';
import { User } from 'components/auth';
import { ModalContent, ModalComponent } from 'components/modal';
import { add } from 'components/toaster';
import Button from 'components/button';

const ipRegex = /^(?:[0-9]{1,3}\.){3}[0-9]{1,3}$/;

export default class TLSNewComponent {
  initialHost: string;
  tlsEmail: Mithril.BasicProperty<string>;
  tlsDomain: Mithril.BasicProperty<string>;
  modal: ModalContent;

  constructor(user: User, email: string, domain: string) {
    this.tlsEmail = m.prop<string>(email || user.email);

    this.initialHost = domain || window.location.hostname;
    if (this.initialHost.match(ipRegex)) {
      this.initialHost = '';
    }
    if (this.initialHost == 'localhost') {
      this.initialHost = '';
    }
    this.tlsDomain = m.prop<string>(this.initialHost);
  }

  register() {
    return m.request({
      url: '/api/v1/settings/tls',
      method: 'POST',
      data: <API.EnableTLSRequest>{
        tlsEmail: this.tlsEmail(),
        tlsDomain: this.tlsDomain(),
        agreed: true,
      }
    })
      .then((res) => {
        console.log(res);
      })
      .catch((res: API.ErrorResponse) => {
        if (!res || !res.errors) {
          console.log(res);
          add('Unknown error', 'error');
          return;
        }
        this.modal = {
          title: 'Error',
          klass: 'modal--error',
          content: () => {
            return m('p', res.errors[0].detail);
          }
        };
      });
  }

  static controller = TLSNewComponent;
  static view = (ctrl: TLSNewComponent) => {
    let warning = null;
    if (ctrl.initialHost == '') {
      warning = m('.tr', `It looks like you're not using a domain; please ensure that you've set up your DNS records correctly.`);
    }
    return m('.table', [
      warning,
      m('.tr.tr--center', [
        m('label', 'TLS Email'),
        m('input.large[type=text]', {
          value: ctrl.tlsEmail(),
          onchange: m.withAttr('value', ctrl.tlsEmail)
        })
      ]),
      m('.tr.tr--center', [
        m('label', 'TLS Domain'),
        m('input.large[type=text]', {
          config: (el: HTMLInputElement, isInitialized: boolean) => {
            if (!isInitialized) {
              el.value = ctrl.initialHost;
            }
          },
          onchange: m.withAttr('value', ctrl.tlsDomain)
        })
      ]),
      m('.tr.tr--right', [
        m('label[for=letos]', [
          'I agree to ',
          m('a', {
            href: 'https://acme-v01.api.letsencrypt.org/terms',
            target: '_blank'
          }, `Let's Encrypt's Terms of Service`)
        ]),
        m('input#letos[type=checkbox]')
      ]),
      m('.tr.tr--right.tr--no-border', [
        m.component(Button, {
          class: 'button--green button--small',
          handler: () => ctrl.register()
        }, 'Enable TLS'),
        m.component(ModalComponent, ctrl.modal),
      ])
    ]);
  }
}