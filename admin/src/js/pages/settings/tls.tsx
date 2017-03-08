import msx from 'lib/msx';
import * as API from 'lib/api';

export default {
  view: (v: Mithril.Vnode<API.TLSSettingsReponse, {}>) => <div class='table'>
    <div class='tr tr--center'>
      <label>TLS Email</label>
      <input class='large' type='text' value={v.attrs.tlsEmail} />
    </div>
    <div class='tr tr--center'>
      <label>TLS Domain</label>
      <input class='large' type='text' value={v.attrs.tlsDomain} />
    </div>
    <div class='tr tr--center'>
      <label>TLS Agreement</label>
      <span class='input-text'>
        <a href={v.attrs.termsOfService}>link</a>
      </span>
    </div>
    <div class='tr tr--center'>
      <label>Agreed On</label>
      <span class='input-text'>{v.attrs.agreedOn}</span>
    </div>
  </div>
};