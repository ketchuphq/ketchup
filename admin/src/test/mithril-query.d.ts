/**
 * Provide definitions for https://github.com/StephanHoyer/mithril-query
 */

declare namespace MithrilQuery {
  interface MithrilQueryStatic {
    <A, S>(element: Mithril.Component<A, S>, attrs?: A): AssertionStatic<A, S>;
    <A, S>(vnode: Mithril.Vnode<A, S>): AssertionStatic<A, S>;
  }

  type Vnode<A,S> = Mithril.Vnode<A, S> & { text: string }

  export interface AssertionStatic<A, S> {
    first(selector: string): Mithril.Vnode<A, S>;
    find(selector: string): Vnode<A,S>[];
    has(selector: string): boolean;
    contains(str: string): boolean;

    redraw(): void;
    click(selector: string): void;
    vnode: Mithril.Vnode<A, S>;
    log(selector: string): void;
    rootNode: any
  }

}
declare var mithrilQuery: MithrilQuery.MithrilQueryStatic;

declare module 'mithril-query' {
  export = mithrilQuery;
}