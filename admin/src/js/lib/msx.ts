import * as m from 'mithril';

export default function msx<A, S>(
  element: string | Mithril.Component<A, S>,
  props: A & Mithril.Lifecycle<A,S>,
  ...children: any[]
): Mithril.Vnode<any, any> {
  if (typeof (element) === 'string') {
    return m(element, props, children);
  }
  return m(element, props, children);
}