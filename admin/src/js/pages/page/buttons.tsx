import * as m from 'mithril';
import msx from 'lib/msx';
import Page from 'lib/page';
import * as Toaster from 'components/toaster';
import PageSaveButtonComponent from 'pages/page/save_button';
import { BaseComponent } from '../../components/auth';

interface PageButtonsAttrs {
  page: Page;
  onsave: (page: Page) => void;
}

export default class PageButtonsComponent extends BaseComponent<PageButtonsAttrs> {
  constructor(v: m.CVnode<PageButtonsAttrs>) {
    super(v)
  }

  publish() {
    // todo: handle case where not saved yet
    this.props.page.publish().then(() => {
      Toaster.add('Page published');
      m.redraw();
    });
  }

  unpublish() {
    this.props.page.unpublish().then(() => {
      Toaster.add('Page unpublished');
      m.redraw();
    });
  }

  delete() {
    this.props.page.delete().then(() => {
      Toaster.add('Page deleted', 'error');
      m.route.set('/admin/pages');
    });
  }

  view() {
    let unpublishButton = <a
      class='button button--small button--blue'
      onclick={(e: Event) => { e.stopPropagation(); this.unpublish(); }}
    >
      Unpublish
    </a>;

    let deleteButton = <a
      class='button button--small button--red'
      onclick={(e: Event) => { e.stopPropagation(); this.delete(); }}
    >
      Delete
    </a>;

    let publishButton = <a
      class='button button--small button--blue'
      onclick={(e: Event) => { e.stopPropagation(); this.publish(); }}
    >
      Publish
    </a>;

    let page = this.props.page;
    return <div class='save-publish'>
      <PageSaveButtonComponent
        page={page}
        classes='button--small' onsave={this.props.onsave} />
      {!page.uuid ? '' : deleteButton}
      {page.isPublished ? unpublishButton : publishButton}
    </div>;
  }
}
