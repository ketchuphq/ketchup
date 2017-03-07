import * as m from 'mithril';

export default function msx(element: string | Mithril.Component<any, any>, props: any, ...children: any[]): Mithril.Vnode<any, any> {
  if (typeof(element) === 'string') {
    return m(element, props, children);
  }
  return m(element, props, children);
}