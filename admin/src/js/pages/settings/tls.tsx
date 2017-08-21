import * as m from 'mithril';
import msx from 'lib/msx';
import * as API from 'lib/api';
import { BaseComponent } from 'components/auth';

export default class TLSComponent extends BaseComponent<API.TLSSettingsResponse> {
  constructor(v: m.CVnode<API.TLSSettingsResponse>) {
    super(v);
  }

  view() {
    return (
      <div class='table'>
        <div class='tr tr--center'>
          <label>TLS Email</label>
          <input class='large' type='text' value={this.props.tlsEmail} />
        </div>
        <div class='tr tr--center'>
          <label>TLS Domain</label>
          <input class='large' type='text' value={this.props.tlsDomain} />
        </div>
        <div class='tr tr--center'>
          <label>TLS Agreement</label>
          <span class='input-text'>
            <a href={this.props.termsOfService}>link</a>
          </span>
        </div>
        <div class='tr tr--center'>
          <label>Agreed On</label>
          <span class='input-text'>
            {this.props.agreedOn}
          </span>
        </div>
      </div>
    );
  }
}
