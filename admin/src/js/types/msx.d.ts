import * as m from 'mithril';

declare global {
  export namespace JSX {
    type Element = m.Vnode<any, any>;

    interface IntrinsicElements {
      [elemName: string]: any;
    }

    interface ElementAttributesProperty {
      // http://www.typescriptlang.org/docs/handbook/jsx.html
      props: any;
    }
  }
}
