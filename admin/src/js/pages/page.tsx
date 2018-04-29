import {PrivateRouteComponentProps} from 'components/auth';
import {ConfirmModalComponent} from 'components/modal';
import * as API from 'lib/api';
import * as Page from 'lib/page';
import Route from 'lib/route';
import Theme from 'lib/theme';
import cloneDeep from 'lodash-es/cloneDeep';
import isEqual from 'lodash-es/isEqual';
import PageControls from 'pages/page/controls';
import PageEditorsComponent from 'pages/page/editors';
import * as React from 'react';

interface State {
  page: API.Page;
  template: API.ThemeTemplate;
  routes: Route[];
  dirty: boolean;
  showSettings: boolean;
  showLeaveModal: boolean;
  pageUUID: string;
  nextRoute: boolean;
}

export default class PagePage extends React.Component<
  PrivateRouteComponentProps<{id: string}>,
  State
> {
  _clickStart: DOMTokenList; // keep track of click location to prevent firing on drag

  pageStore: Page.Store;
  initialContent: API.Content[];
  pageRef: React.RefObject<HTMLDivElement>;

  constructor(props: PrivateRouteComponentProps<{id: string}>) {
    super(props);
    this.pageStore = new Page.Store();
    this.state = {
      page: null,
      template: null,
      routes: [],
      dirty: false,
      showSettings: false,
      showLeaveModal: false,
      pageUUID: props.match.params.id,
      nextRoute: false,
    };
    this.pageRef = React.createRef();

    // when a page is updated, we need to make sure we process the new
    // theme and template
    let initial = true;
    this.pageStore.subscribe('index', (page) => {
      return this.fetchThemeTemplate(page.theme, page.template).then(() => {
        if (!this.state.page || initial) {
          this.initialContent = cloneDeep(page.contents);
        }
        this.setState({page: page});
        initial = false;
      });
    });
  }

  componentDidMount() {
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
        .then((routes) => this.setState({routes}));
    } else {
      // new page
      const page = Page.newPage();
      page.authors = [{uuid: this.props.user.uuid}];
      this.pageStore.set(page);
    }
  }

  toggleSettings = () => {
    this.setState({showSettings: !this.state.showSettings});
  };

  // fetchThemeTemplate fetches the Theme from the backend
  // Note: this gets called in the page update callback.
  fetchThemeTemplate = (theme: string, template: string): Promise<Theme | void> => {
    return Theme.get(theme).then(
      (t) => {
        this.setState({template: t.getTemplate(template)});
        return t;
      },
      () => {
        if (theme == 'none' && template == 'html') {
          return;
        }
        // catch deleted theme
        return Theme.get('none').then((t) => {
          this.pageStore.setThemeTemplate(t, t.getTemplate('html'));
          return t;
        });
      }
    );
  };

  confirmLeave = () => {
    if (this.state.showSettings) {
      this.toggleSettings();
      return;
    }
    this.setState((prev) => {
      let dirty = prev.dirty || !isEqual(this.initialContent, this.state.page.contents);
      let showLeaveModal = dirty || prev.showLeaveModal;
      return {dirty, showLeaveModal, nextRoute: !showLeaveModal};
    });
  };

  render() {
    if (!this.state.page) {
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
          if (validClick) {
            this.confirmLeave();
          }
        }}
        ref={this.pageRef}
      >
        <PageControls
          store={this.pageStore}
          routes={this.state.routes}
          toggleSettings={this.toggleSettings}
          showSettings={this.state.showSettings}
          leave={this.confirmLeave}
        />
        <div className="page-editor">
          <div className="controls">
            <input
              type="text"
              className="large"
              placeholder="title..."
              value={this.state.page.title || ''}
              onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
                e.preventDefault();
                this.pageStore.update((page) => {
                  page.title = e.target.value;
                });
              }}
            />
          </div>
          <PageEditorsComponent contents={this.state.page.contents} />
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
