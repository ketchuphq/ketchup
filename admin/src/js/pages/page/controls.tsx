import Popover from 'components/popover';
import * as API from 'lib/api';
import * as Page from 'lib/page';
import PageButtonsComponent from 'pages/page/buttons';
import PageSaveButtonComponent from 'pages/page/save_button';
import PageThemePickerComponent from 'pages/page/theme_picker';
import PageEditRoutesComponent from 'pages/page/edit_route';
import * as React from 'react';
import GenericStore, {Data} from 'lib/store';

interface ControlsProps {
  store: Page.Store;
  routesStore: GenericStore<Data<API.Route[]>>;
  toggleSettings: () => void;
  showSettings: boolean;
  leave: () => void;
}

export default class PageControls extends React.Component<ControlsProps, {}> {
  render() {
    const stores = {
      store: this.props.store,
      routesStore: this.props.routesStore,
    };
    return (
      <div className="page-max__controls">
        <PageSaveButtonComponent {...stores} />
        <span className="typcn typcn-cog" onClick={() => this.props.toggleSettings()} />
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
