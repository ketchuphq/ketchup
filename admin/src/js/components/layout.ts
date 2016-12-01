import Navigation from 'components/navigation';
import * as Toaster from 'components/toaster';

export default (content: string | Mithril.VirtualElement) =>
  m('.container', { key: 'container' }, [
    Toaster.render(),
    m.component(Navigation),
    m('.container--body', content)
  ]);
