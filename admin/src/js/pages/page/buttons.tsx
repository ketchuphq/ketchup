import * as m from 'mithril';
import msx from 'lib/msx';
import Page from 'lib/page';
import * as Toaster from 'components/toaster';
import PageSaveButtonComponent from 'pages/page/save_button';

interface PageButtonsAttrs {
  page: Page;
  onsave: (page: Page) => void;
}

export default class PageButtonsComponent {
  private readonly _attrs: PageButtonsAttrs;

  publish(page: Page) {
    // todo: handle case where not saved yet
    page.publish().then(() => {
      Toaster.add('Page published');
      m.redraw();
    });
  }

  unpublish(page: Page) {
    page.unpublish().then(() => {
      Toaster.add('Page unpublished');
      m.redraw();
    });
  }

  delete(page: Page) {
    page.delete().then(() => {
      Toaster.add('Page deleted', 'error');
      m.route.set('/admin/pages');
    });
  }

  static oninit(v: Mithril.Vnode<PageButtonsAttrs, PageButtonsComponent>) {
    v.state = new PageButtonsComponent();
  }
  static view({ attrs: { page, onsave }, state }: Mithril.Vnode<PageButtonsAttrs, PageButtonsComponent>) {
    let unpublishButton = <a
      class='button button--small button--blue'
      onclick={(e: Event) => { e.stopPropagation(); state.unpublish(page); }}
    >
      Unpublish
    </a>;

    let deleteButton = <a
      class='button button--small button--red'
      onclick={(e: Event) => { e.stopPropagation(); state.delete(page); }}
    >
      Delete
    </a>;

    let publishButton = <a
      class='button button--small button--blue'
      onclick={(e: Event) => { e.stopPropagation(); state.publish(page); }}
    >
      Publish
    </a>;

    return <div class='save-publish'>
      <PageSaveButtonComponent
        page={page}
        classes='button--small' onsave={onsave} />
      {!page.uuid ? '' : deleteButton}
      {page.isPublished ? unpublishButton : publishButton}
    </div>;
  }
}