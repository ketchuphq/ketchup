import msx from 'lib/msx';
import * as m from 'mithril';
import * as API from 'lib/api';
import { User, BaseComponent } from 'components/auth';
import { ModalComponent } from 'components/modal';
import { add } from 'components/toaster';
import Button from 'components/button';

const leURL = 'https://acme-v01.api.letsencrypt.org/terms';
const ipRegex = /^(?:[0-9]{1,3}\.){3}[0-9]{1,3}$/;

interface TLSNewComponentAttrs {
  user: User;
  email: string;
  domain: string;
}

export default class TLSNewComponent extends BaseComponent<TLSNewComponentAttrs> {
  initialHost: string;
  tlsEmail: string;
  tlsDomain: string;
  showErrorModal: boolean;
  errors: string;

  constructor(v: m.CVnode<TLSNewComponentAttrs>) {
    super(v)
    this.tlsEmail = v.attrs.email || v.attrs.user.email;
    this.initialHost = v.attrs.domain || window.location.hostname;
    this.showErrorModal = false;
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
        this.errors = res.errors[0].detail;
        this.showErrorModal = true
        m.redraw();
      });
  }

  view() {
    let warning = null;
    if (this.initialHost == '') {
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
          value={this.tlsEmail}
          onchange={m.withAttr('value', (e) => this.tlsEmail = e)}
        />
      </div>
      <div class='tr tr--center'>
        <label>TLS Domain</label>
        <input
          class='large'
          type='text'
          config={(el: HTMLInputElement, isInitialized: boolean) => {
            if (!isInitialized) {
              el.value = this.initialHost;
            }
          }}
          onchange={m.withAttr('value', (v) => this.tlsDomain = v)}
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
          handler={() => this.register()}>
          Enable TLS
        </Button>
        <ModalComponent
          title='Error'
          klass='modal--error'
          visible={() => this.showErrorModal}
          toggle={() => {
            this.showErrorModal = !this.showErrorModal;
            m.redraw();
          }}
          >
            <p>{this.errors}</p>
        </ModalComponent>
      </div>
    </div>;
  }
}
