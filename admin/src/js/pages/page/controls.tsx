import Popover from 'components/popover';
import * as API from 'lib/api';
import * as Page from 'lib/page';
import GenericStore, {Data} from 'lib/store';
import PageButtonsComponent from 'pages/page/buttons';
import PageEditRoutesComponent from 'pages/page/edit_route';
import PageSaveButtonComponent from 'pages/page/save_button';
import PageThemePickerComponent from 'pages/page/theme_picker';
import * as React from 'react';

interface Props {
  store: Page.Store;
  routesStore: GenericStore<Data<API.Route[]>>;
  toggleSettings: () => void;
  togglePreview: () => void;
  showSettings: boolean;
  leave: () => void;
}

interface State {
  route?: string;
  published: boolean;
}

export default class PageControls extends React.Component<Props, State> {
  constructor(props: Props) {
    super(props);
    let route;
    let routes = this.props.routesStore.obj.initial;
    if (routes && routes.length > 0) {
      route = routes[0].path;
    }
    this.state = {
      route,
      published: Page.isPublished(this.props.store.page),
    };
  }

  componentDidMount() {
    this.props.store.subscribe('page-controls', (page) => {
      this.setState({
        published: Page.isPublished(page),
      });
    });

    this.props.routesStore.subscribe('page-controls', (data) => {
      let routes = data.initial;
      if (routes.length > 0) {
        this.setState({
          route: routes[0].path,
        });
      }
    });
  }

  componentWillUnmount() {
    this.props.store.unsubscribe('page-controls');
    this.props.routesStore.unsubscribe('page-controls');
  }

  render() {
    const stores = {
      store: this.props.store,
      routesStore: this.props.routesStore,
    };

    return (
      <div className="page-max__controls">
        <PageSaveButtonComponent {...stores} />
        {this.state.route && this.state.published ? (
          <a
            title="Open page in new tab"
            className="typcn typcn-link"
            href={this.state.route}
            target="_blank"
          />
        ) : null}
        <span
          title="Toggle preview"
          className="typcn typcn-zoom"
          onClick={() => this.props.togglePreview()}
        />
        <span
          title="Settings"
          className="typcn typcn-cog"
          onClick={() => this.props.toggleSettings()}
        />
        <Popover visible={this.props.showSettings}>
          <div className="controlset">
            <div className="settings">
              <div className="controls">
                <div className="control">
                  {this.props.store.page ? <PageEditRoutesComponent {...stores} /> : null}
                </div>
              </div>
              <div className="controls">
                <PageThemePickerComponent store={this.props.store} />
              </div>
            </div>
            <PageButtonsComponent {...stores} />
          </div>
        </Popover>
        <a
          className="typcn typcn-times"
          href="/admin/pages"
          onClick={(e: React.MouseEvent<HTMLAnchorElement>) => {
            e.preventDefault();
            this.props.leave();
          }}
        />
      </div>
    );
  }
}
