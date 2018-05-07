import * as React from 'react';
import debounce from 'lodash-es/debounce';
import {UserContext} from 'components/auth';
import Button from 'components/button';
import {Link} from 'react-router-dom';
let store = require('store/dist/store.modern') as StoreJsAPI;

interface State {
  collapsed: boolean;
}

const NavLink: React.SFC<{
  url: string;
  additionalClasses?: string;
  icon?: string;
  children?: any;
}> = (props) => (
  <Link to={props.url} className={`nav-link ${props.additionalClasses}`}>
    {!!props.icon ? <span className={`typcn typcn-${props.icon}`} /> : ''}
    <span className="nav-link__text">{props.children}</span>
  </Link>
);

export default class Navigation extends React.Component<{}, State> {
  constructor(props: any, context?: any) {
    super(props, context);
    this.state = {
      collapsed: store.get('hideMenu') || window.innerWidth <= 480,
    };
    window.addEventListener('resize', this.resizeHandler);
    this.resizeHandler();
  }

  componentWillUnmount() {
    window.removeEventListener('resize', this.resizeHandler);
  }

  resizeHandler = debounce(() => {
    if (window.innerWidth > 480) {
      return;
    }
    store.set('hideMenu', true);
    this.setState({collapsed: true});
  }, 300);

  toggle = () => {
    this.setState((prev) => {
      let collapsed = !prev.collapsed;
      store.set('hideMenu', collapsed);
      return {collapsed};
    });
  };

  render() {
    let navClass = 'navigation';
    if (this.state.collapsed) {
      navClass += ' navigation--hidden';
    }

    return (
      <UserContext.Consumer>
        {(user) => {
          if (!user) {
            return (
              <div className={navClass}>
                <NavLink url="/" additionalClasses="nav-title">
                  K
                </NavLink>
                <NavLink url="/login">Login</NavLink>
              </div>
            );
          }

          return (
            <div className={navClass}>
              <NavLink url="/" additionalClasses="nav-title">
                K
              </NavLink>
              <div className="nav-button">
                <Button className="button--green button--center" href="/compose">
                  <span className="typcn typcn-edit" />
                  <span className="nav-link__text">Compose</span>
                </Button>
              </div>
              <NavLink url="/pages" icon="document-text">
                Pages
              </NavLink>
              <NavLink url="/themes" icon="brush">
                Theme
              </NavLink>
              <NavLink url="/data" icon="th-small">
                Data
              </NavLink>
              <NavLink url="/files" icon="document">
                Files
              </NavLink>
              <NavLink url="/settings" icon="spanner-outline">
                Settings
              </NavLink>
              <a className="nav-link" href="/admin/logout">
                <span className="typcn typcn-weather-night" />
                <span className="nav-link__text">Log Out</span>
              </a>
              <a className="nav-link nav-link--toggle" onClick={() => this.toggle()}>
                <span
                  className={`typcn typcn-arrow-${this.state.collapsed ? 'maximise' : 'minimise'}`}
                />
              </a>
            </div>
          );
        }}
      </UserContext.Consumer>
    );
  }
}
