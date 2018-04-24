import {User} from 'components/auth';
import Button from 'components/button';
import {ModalComponent} from 'components/modal';
import * as API from 'lib/api';
import {post} from 'lib/requests';
import * as React from 'react';

const leURL = 'https://acme-v01.api.letsencrypt.org/terms';
const ipRegex = /^(?:[0-9]{1,3}\.){3}[0-9]{1,3}$/;

interface Props {
  user: User;
  email: string;
  domain: string;
}

interface State {
  initialHost: string;
  tlsEmail: string;
  tlsDomain: string;
  showErrorModal: boolean;
  errors?: string;
}

export default class TLSNewComponent extends React.Component<Props, State> {
  constructor(props: Props) {
    super(props);
    let initialHost = props.domain || window.location.hostname;
    if (initialHost.match(ipRegex)) {
      initialHost = '';
    }
    if (initialHost == 'localhost') {
      initialHost = '';
    }
    this.state = {
      tlsEmail: props.email || props.user.email,
      tlsDomain: initialHost,
      initialHost: initialHost,
      showErrorModal: false,
    };
  }

  register() {
    return post('/api/v1/settings/tls', {
      tlsEmail: this.state.tlsEmail,
      tlsDomain: this.state.tlsDomain,
      agreed: true,
    } as API.EnableTLSRequest)
      .then((res) => res.json())
      .then((res: API.ErrorResponse) => {
        if (res && res.errors) {
          // todo: check res code
          // add('Unknown error', 'error');
          // return;
          this.setState({
            errors: res.errors[0].detail,
            showErrorModal: true,
          });
        }
      });
  }

  render() {
    let warning = null;
    if (this.state.initialHost == '') {
      warning = (
        <div className="tr">
          It looks like you're not using a domain; please ensure that you've set up your DNS records
          correctly.
        </div>
      );
    }
    return (
      <div className="table">
        {warning}
        <div className="tr tr--center">
          <label>TLS Email</label>
          <input
            className="large"
            type="text"
            value={this.state.tlsEmail}
            onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
              this.setState({tlsEmail: e.target.value});
            }}
          />
        </div>
        <div className="tr tr--center">
          <label>TLS Domain</label>
          <input
            className="large"
            type="text"
            value={this.state.tlsDomain}
            onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
              this.setState({tlsDomain: e.target.value});
            }}
          />
        </div>
        <div className="tr tr--right tr--tos">
          <label htmlFor="tos">
            I agree to{' '}
            <a href={leURL} target="_blank">
              Let's Encrypt's Terms of Service
            </a>
          </label>
        </div>
        <div className="tr tr--right tr--no-border">
          <Button className="button--green button--small" handler={() => this.register()}>
            Enable TLS
          </Button>
          <ModalComponent
            title="Error"
            klass="modal--error"
            visible={this.state.showErrorModal}
            toggle={() => {
              this.setState((prev) => ({showErrorModal: !prev.showErrorModal}));
            }}
          >
            <p>{this.state.errors}</p>
          </ModalComponent>
        </div>
      </div>
    );
  }
}
