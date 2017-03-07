import msx from 'lib/msx';
import * as m from 'mithril';
import Navigation from 'components/navigation';
import * as Toaster from 'components/toaster';

export default (component: Mithril.Component<any, any>) => ({
  render: () => m('.container', [
    Toaster.render(),
    <Navigation />,
    <div class='container__body'>{m(component)}</div>
  ])
});