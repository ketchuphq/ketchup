import * as React from 'react';
import * as API from 'lib/api';
import * as Page from 'lib/page';
import Route from 'lib/route';
import * as Toaster from 'components/toaster';
import GenericStore, {Data} from 'lib/store';
import cloneDeep from 'lodash-es/cloneDeep';

interface Props {
  store: Page.Store;
  routesStore: GenericStore<Data<API.Route[]>>;
  classes?: string;
}

export default class PageSaveButtonComponent extends React.Component<Props> {
  constructor(props: Props) {
    super(props);
  }

  save = (e: React.MouseEvent<any>) => {
    e.preventDefault();
    const {store, routesStore} = this.props;
    const currentRoutes = routesStore.obj.current;
    let isNewPage = !store.obj.uuid;
    store
      .save()
      .then((page) => {
        window.history.replaceState(null, page.title, `/admin/pages/${page.uuid}`);
        return Route.saveRoutes(page, currentRoutes, isNewPage);
      })
      .then((routes) => {
        routesStore.update((data) => {
          data.initial = cloneDeep(routes);
          data.current = cloneDeep(routes);
        });
        Toaster.add('Page successfully saved');
      })
      .catch((err: any) => {
        if (err.detail) {
          Toaster.add(err.detail, 'error');
        } else {
          Toaster.add('Internal server error.', 'error');
          console.error(err);
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
