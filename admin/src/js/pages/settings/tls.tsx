import * as React from 'react';
import * as API from 'lib/api';

export default class TLSComponent extends React.Component<API.TLSSettingsResponse> {
  constructor(props: API.TLSSettingsResponse) {
    super(props);
  }

  render() {
    return (
      <div className="table">
        <div className="tr tr--center">
          <label>TLS Email</label>
          <input className="large" type="text" value={this.props.tlsEmail} />
        </div>
        <div className="tr tr--center">
          <label>TLS Domain</label>
          <input className="large" type="text" value={this.props.tlsDomain} />
        </div>
        <div className="tr tr--center">
          <label>TLS Agreement</label>
          <span className="input-text">
            <a href={this.props.termsOfService}>link</a>
          </span>
        </div>
        <div className="tr tr--center">
          <label>Agreed On</label>
          <span className="input-text">{this.props.agreedOn}</span>
        </div>
      </div>
    );
  }
}
