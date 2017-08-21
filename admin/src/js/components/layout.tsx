import msx from 'lib/msx';
import * as m from 'mithril';
import Navigation from 'components/navigation';
import * as Toaster from 'components/toaster';
import { getUser } from 'components/auth';

export default (component: m.ComponentTypes<any, any>) => ({
  onmatch: (_: any, requestedPath: string) => {
    getUser()
      .catch(() => {
        if (requestedPath.match('^/admin/?$')) {
          m.route.set('/admin/login'); // default path is admin
        } else {
          m.route.set(`/admin/login?next=${requestedPath}`);
        }
      });
  },
  render: () => <div class='container'>
    {Toaster.render()}
    <Navigation />
    <div class='container__body'>{m(component)}</div>
  </div>
});
