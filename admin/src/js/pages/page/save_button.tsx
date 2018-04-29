import * as React from 'react';
import * as API from 'lib/api';
import * as Page from 'lib/page';
import Route from 'lib/route';
import * as Toaster from 'components/toaster';

interface Props {
  store: Page.Store;
  routes: API.Route[];
  classes?: string;
}

export default class PageSaveButtonComponent extends React.Component<Props> {
  constructor(props: Props) {
    super(props);
  }

  save = (e: React.MouseEvent<any>) => {
    e.preventDefault();
    const {store, routes} = this.props;
    store
      .save()
      .then(() => {
        window.history.replaceState(null, store.page.title, `/admin/pages/${store.page.uuid}`);
        return Route.saveRoutes(store.page, routes);
      })
      .then(() => {
        Toaster.add('Page successfully saved');
      })
      .catch((err: any) => {
        if (err.detail) {
          Toaster.add(err.detail, 'error');
        } else {
          Toaster.add('Internal server error.', 'error');
          console.log(err);
        }
      });
  };

  render() {
    return (
      <a className={`button button--green ${this.props.classes || ''}`} onClick={this.save}>
        Save
      </a>
    );
  }
}
