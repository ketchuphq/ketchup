import msx from 'lib/msx';
import * as m from 'mithril';
import * as API from 'lib/api';
import { User } from 'components/auth';
import { ModalAttrs, ModalComponent } from 'components/modal';
import { add } from 'components/toaster';
import Button from 'components/button';

const leURL = 'https://acme-v01.api.letsencrypt.org/terms';
const ipRegex = /^(?:[0-9]{1,3}\.){3}[0-9]{1,3}$/;

interface TLSNewComponentAttrs {
  user: User;
  email: string;
  domain: string;
}

export default class TLSNewComponent {
  initialHost: string;
  tlsEmail: string;
  tlsDomain: string;
  modal: ModalAttrs;

  constructor(attrs: TLSNewComponentAttrs) {
    this.tlsEmail = attrs.email || attrs.user.email;
    this.initialHost = attrs.domain || window.location.hostname;
    if (this.initialHost.match(ipRegex)) {
      this.initialHost = '';
    }
    if (this.initialHost == 'localhost') {
      this.initialHost = '';
    }
    this.tlsDomain = this.initialHost;
  }

  register() {
    return m.request({
      url: '/api/v1/settings/tls',
      method: 'POST',
      data: {
        tlsEmail: this.tlsEmail,
        tlsDomain: this.tlsDomain,
        agreed: true,
      } as API.EnableTLSRequest
    })
      .catch((res: API.ErrorResponse) => {
        if (!res || !res.errors) {
          add('Unknown error', 'error');
          return;
        }
        this.modal = {
          title: 'Error',
          klass: 'modal--error',
          content: () => {
            return <p>{res.errors[0].detail}</p>;
          }
        };
      });
  }

  static oninit(v: Mithril.Vnode<TLSNewComponentAttrs, TLSNewComponent>) {
    v.state = new TLSNewComponent(v.attrs);
  };

  static view(v: Mithril.Vnode<TLSNewComponentAttrs, TLSNewComponent>) {
    let ctrl = v.state;
    let warning = null;
    if (ctrl.initialHost == '') {
      warning = <div class='tr'>
        It looks like you're not using a domain; please ensure that you've set up your DNS records correctly.
      </div>;
    }
    return <div class='table'>
      {warning}
      <div class='tr tr--center'>
        <label>TLS Email</label>
        <input
          class='large'
          type='text'
          value={ctrl.tlsEmail}
          onchange={m.withAttr('value', (e) => ctrl.tlsEmail = e)}
        />
      </div>
      <div class='tr tr--center'>
        <label>TLS Domain</label>
        <input
          class='large'
          type='text'
          config={(el: HTMLInputElement, isInitialized: boolean) => {
            if (!isInitialized) {
              el.value = ctrl.initialHost;
            }
          }}
          onchange={m.withAttr('value', (v) => ctrl.tlsDomain = v)}
        />
      </div>
      <div class='tr tr--right tr--tos'>
        <label for='tos'>
          I agree to <a href={leURL} target='_blank'>Let's Encrypt's Terms of Service</a>
        </label>
      </div>
      <div class='tr tr--right tr--no-border'>
        <Button
          class='button--green button--small'
          handler={() => ctrl.register()}>
          Enable TLS
        </Button>
        {m(ModalComponent, ctrl.modal)}
      </div>
    </div>;
  }
}