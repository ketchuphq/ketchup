import msx from 'lib/msx';
import * as m from 'mithril';
import Page from 'lib/page';

interface Callbacks {
  save: () => void;
  unpublish: () => void;
  delete: () => void;
  publish: () => void;
}

interface PageButtonsAttrs {
  page: Page;
  callbacks: Callbacks;
}

export default class PageButtonsComponent {
  private readonly _attrs: PageButtonsAttrs;
  static view(v: Mithril.Vnode<PageButtonsAttrs, {}>) {
    let { page, callbacks } = v.attrs;
    let saveButton = <a
      class='button button--small button--green'
      onclick={(e: Event) => { e.stopPropagation(); callbacks.save(); }}
    >
      Save
    </a>;

    let unpublishButton = <a
      class='button button--small button--blue'
      onclick={(e: Event) => { e.stopPropagation(); callbacks.unpublish(); }}
    >
      Unpublish
    </a>;

    let deleteButton = <a
      class='button button--small button--red'
      onclick={(e: Event) => { e.stopPropagation(); callbacks.delete(); }}
    >
      Delete
    </a>;

    let publishButton = <a
      class='button button--small button--blue'
      onclick={(e: Event) => { e.stopPropagation(); callbacks.publish(); }}
    >
      Publish
    </a>;

    return <div class='save-publish'>
      {saveButton}
      {!page.uuid ? '' : deleteButton}
      {page.isPublished ? [unpublishButton] : publishButton}
    </div>;
  }
}