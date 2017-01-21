import * as m from 'mithril';

declare global {
  export namespace JSX {
    type Element = Mithril.VirtualElement;

    interface IntrinsicElements {
      [elemName: string]: Mithril.Attributes;
    }
    interface ElementAttributesProperty {
      // http://www.typescriptlang.org/docs/handbook/jsx.html
      props: any; // specify the property name on the element to use
    }
  }
}
