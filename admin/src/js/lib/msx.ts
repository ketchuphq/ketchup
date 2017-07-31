import * as m from 'mithril';

export default function msx<A, S>(
  element: string | m.Component<A, S>,
  props: A & m.Lifecycle<A,S>,
  ...children: any[]
): m.Vnode<any, any> {
  if (!element) {
    return null;
  }
  if (typeof (element) === 'string') {
    return m(element, props, children);
  }
  return m(element, props, children);
}
