import Navigation from 'components/navigation';
import * as Toaster from 'components/toaster';

export default (content: string | Mithril.VirtualElement, animate: boolean = false) =>
  m('.container', {
    class: animate ? 'container--in': '',
    key: 'container',
  }, [
    Toaster.render(),
    m.component(Navigation),
    m('.container__body', content)
  ]);
