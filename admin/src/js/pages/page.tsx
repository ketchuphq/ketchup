import {PrivateRouteComponentProps} from 'components/auth';
import {ConfirmModalComponent} from 'components/modal';
import * as API from 'lib/api';
import * as Page from 'lib/page';
import Route from 'lib/route';
import GenericStore, {Data} from 'lib/store';
import cloneDeep from 'lodash-es/cloneDeep';
import isEqual from 'lodash-es/isEqual';
import PageControls from 'pages/page/controls';
import PageEditorsComponent from 'pages/page/editors';
import * as React from 'react';
let store = require('store/dist/store.modern') as StoreJsAPI;

interface State {
  pageTitle?: string;
  contents?: API.Content[];

  dirty: boolean;
  showSettings: boolean;
  showLeaveModal: boolean;
  showPreview: boolean;
  pageUUID: string;
  nextRoute: boolean;
}

export default class PagePage extends React.Component<
  PrivateRouteComponentProps<{id: string}>,
  State
> {
  _clickStart: DOMTokenList; // keep track of click location to prevent firing on drag

  pageStore: Page.Store;
  routesStore: GenericStore<Data<API.Route[]>>;
  pageRef: React.RefObject<HTMLDivElement>;

  constructor(props: PrivateRouteComponentProps<{id: string}>) {
    super(props);
    this.pageStore = new Page.Store();
    this.routesStore = new GenericStore<Data<API.Route[]>>(
      (from, _?) => ({
        initial: from.initial,
        current: from.current.slice(),
      }),
      {}
    );
    this.state = {
      dirty: false,
      showSettings: false,
      showLeaveModal: false,
      showPreview: store.get('showPreview', false),
      pageUUID: props.match.params.id,
      nextRoute: false,
    };
    this.pageRef = React.createRef();
  }

  componentDidMount() {
    // must subscribe to page store events before triggering events
    this.pageStore.subscribe('page-index', (page) => {
      this.setState({
        pageTitle: page.title,
        contents: page.contents,
      });
    });

    // handle animations
    this.pageRef.current.addEventListener('mousedown', (e: any) => {
      this._clickStart = e.target.classList;
    });
    this.pageRef.current.addEventListener('animationend', (e: AnimationEvent) => {
      // old animation is removed, otherwise new animations won't fire.
      if (e.animationName == 'fadeIn') {
        this.pageRef.current.classList.add('animate-fade-in-complete');
        this.pageRef.current.classList.remove('animate-fade-in');
      }
      // navigate away after zoomAway animation completes
      if (e.animationName == 'zoomAway') {
        this.props.history.push('/pages');
      }
    });

    // handle loading page
    if (this.state.pageUUID) {
      // load page
      this.pageStore
        .get(this.state.pageUUID)
        .then((page) => Route.getRoutes(page))
        .then((routes) => {
          if (!routes || routes.length == 0) {
            routes.push({});
          }
          this.routesStore.set({
            initial: cloneDeep(routes),
            current: routes,
          });
        });
    } else {
      // new page
      const page = Page.newPage();
      page.authors = [{uuid: this.props.user.uuid}];
      this.pageStore.set(page);
      this.routesStore.set({
        initial: [],
        current: [{}],
      });
    }
  }

  componentWillUnmount() {
    this.pageStore.unsubscribe('page-index');
  }

  toggleSettings = () => {
    this.setState({showSettings: !this.state.showSettings});
  };

  togglePreview = () => {
    this.setState((prev) => {
      store.set('showPreview', !prev.showPreview);
      return {showPreview: !prev.showPreview};
    });
  };

  confirmLeave = () => {
    if (this.state.showSettings) {
      this.toggleSettings();
      return;
    }
    this.setState((prev) => {
      let dirty =
        this.pageStore.hasChanges() ||
        !isEqual(this.routesStore.obj.initial, this.routesStore.obj.current);
      let showLeaveModal = dirty || prev.showLeaveModal;
      return {dirty, showLeaveModal, nextRoute: !showLeaveModal};
    });
  };

  render() {
    if (!this.pageStore.page) {
      return <div ref={this.pageRef} />;
    }

    let pageMaxClasses = 'page-max animate-fade-in';
    if (this.state.nextRoute) {
      pageMaxClasses = 'page-max animate-zoom-away animate-fill';
    }

    return (
      <div
        className={pageMaxClasses}
        onClick={(e: any) => {
          let validClick =
            e.target.classList.contains('page-max') &&
            this._clickStart &&
            this._clickStart.contains('page-max');
          if (validClick && !this.state.showPreview) {
            this.confirmLeave();
          }
        }}
        ref={this.pageRef}
      >
        <PageControls
          store={this.pageStore}
          routesStore={this.routesStore}
          toggleSettings={this.toggleSettings}
          togglePreview={this.togglePreview}
          showSettings={this.state.showSettings}
          leave={this.confirmLeave}
        />
        <div className={`page-editor ${this.state.showPreview ? 'page-editor--previews' : ''}`}>
          <div className="controls controls--title">
            <input
              type="text"
              className="large"
              placeholder="title..."
              value={this.state.pageTitle || ''}
              onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
                e.preventDefault();
                this.pageStore.update((page) => {
                  page.title = e.target.value;
                });
              }}
            />
          </div>
          <PageEditorsComponent
            showPreview={this.state.showPreview}
            contents={this.state.contents}
          />
        </div>
        <ConfirmModalComponent
          title="You are about to leave this page"
          visible={this.state.showLeaveModal}
          toggle={() => {
            this.setState((state) => ({
              showLeaveModal: !state.showLeaveModal,
            }));
          }}
          confirmText="Stay"
          cancelText="Leave"
          cancelColor="modal-button--red"
          reject={() => this.setState({nextRoute: true})}
        >
          <p>You have unsaved changes. Are you sure you want to leave this page?</p>
        </ConfirmModalComponent>
      </div>
    );
  }
}
