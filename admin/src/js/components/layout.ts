import Navigation from './navigation';

export default (content: string | Mithril.VirtualElement) =>
  m('.container', [
    m.component(Navigation),
    m('.container--body', content)
  ]);
