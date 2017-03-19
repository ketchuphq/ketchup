import * as _ from 'mithril';

declare global {
  export namespace JSX {
    type Element = Mithril.Vnode<any, any>;

    interface IntrinsicElements {
      [elemName: string]: any;
    }

    interface ElementAttributesProperty {
      // http://www.typescriptlang.org/docs/handbook/jsx.html
      _attrs: any;
    }
  }
}
