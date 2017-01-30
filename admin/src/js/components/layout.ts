import Navigation from 'components/navigation';
import * as Toaster from 'components/toaster';

let retain = (_: any, __: any, context: Mithril.Context) => context.retain = true;

export default (content: string | Mithril.VirtualElement, animate: boolean = false) =>
  m('.container', {
    class: animate ? 'container--in': '',
    key: 'container',
    config: retain
  }, [
    Toaster.render(),
    m.component(Navigation),
    m('.container__body', content)
  ]);
