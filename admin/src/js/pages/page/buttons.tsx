import * as React from 'react';
import * as Page from 'lib/page';
import * as Toaster from 'components/toaster';
import PageSaveButtonComponent from 'pages/page/save_button';
import * as API from 'lib/api';
import GenericStore, {Data} from 'lib/store';

interface Props {
  store: Page.Store;
  routesStore: GenericStore<Data<API.Route[]>>;
}

interface State {
  isPublished: boolean;
}

export default class PageButtonsComponent extends React.Component<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = {
      isPublished: Page.isPublished(props.store.page),
    };
  }

  publish = (e: React.MouseEvent<any>) => {
    e.preventDefault();
    // todo: handle case where not saved yet
    this.props.store.publish().then(() => {
      this.setState({isPublished: true});
      Toaster.add('Page published');
    });
  };

  unpublish = (e: React.MouseEvent<any>) => {
    e.preventDefault();
    this.props.store.unpublish().then(() => {
      Toaster.add('Page unpublished');
      this.setState({isPublished: false});
    });
  };

  delete = (e: React.MouseEvent<any>) => {
    e.preventDefault();
    Page.deletePage(this.props.store.page).then(() => {
      Toaster.add('Page deleted');
      location.assign('/admin/pages');
    });
  };

  render() {
    let unpublishButton = (
      <a className="button button--small button--blue" onClick={this.unpublish}>
        Unpublish
      </a>
    );

    let deleteButton = (
      <a className="button button--small button--red" onClick={this.delete}>
        Delete
      </a>
    );

    let publishButton = (
      <a className="button button--small button--blue" onClick={this.publish}>
        Publish
      </a>
    );

    let page = this.props.store.page;
    return (
      <div className="save-publish">
        <PageSaveButtonComponent classes="button--small" {...this.props} />
        {!page.uuid ? '' : deleteButton}
        {Page.isPublished(page) ? unpublishButton : publishButton}
      </div>
    );
  }
}
