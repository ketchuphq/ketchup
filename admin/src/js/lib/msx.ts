import * as m from 'mithril';

export default function msx(element: string | Mithril.Component<any>, props: any, ...children: any[]) {
  if (typeof(element) === 'string') {
    return m(element, props, children);
  }
  return m.component(element, props, children);
}