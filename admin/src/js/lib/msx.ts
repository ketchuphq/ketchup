import * as m from 'mithril';

export default function msx(element: string, props: any, ...children: any[]) {
  return m(element, props, children)
}