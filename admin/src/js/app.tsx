import * as React from 'react';
import * as ReactDOM from 'react-dom';
import * as ReactModal from 'react-modal';

import HomePage from 'pages/home';
// import RoutesPage from 'pages/routes';
import PagesPage from 'pages/pages';
import PagePage from 'pages/page';
import LoginPage from 'pages/login';
import ThemePage from 'pages/theme';
import ThemesPage from 'pages/themes';
import TemplatePage from 'pages/template';
import SettingsPage from 'pages/settings';
import DataPage from 'pages/data';
import InstallThemePage from 'pages/install-theme';
import * as WebFont from 'webfontloader';
import {BrowserRouter as Router, Route as PublicRoute} from 'react-router-dom';
import {getUser, User, UserContext, PrivateRoute as Route} from 'components/auth';

interface State {
  user?: User;
  userLoaded: boolean;
}

class App extends React.Component<{}, State> {
  constructor(props: any) {
    super(props);
    this.state = {userLoaded: false};
  }

  componentDidMount() {
    getUser().then(
      (user) => {
        this.setState({user, userLoaded: true});
      },
      () => {
        this.setState({userLoaded: true});
      }
    );
  }

  render() {
    return (
      <UserContext.Provider value={this.state.user}>
        <Router basename="/admin">
          <div id="app">
            <Route {...this.state} exact path="/" component={HomePage} />
            <PublicRoute path="/login" component={LoginPage} />
            {/* <Route path="/routes" component={RoutesPage} /> */}
            <Route {...this.state} path="/pages" component={PagesPage} />
            <Route {...this.state} path="/pages/:id" component={PagePage} />
            <Route {...this.state} path="/compose" component={PagePage} />
            <Route {...this.state} exact path="/themes" component={ThemesPage} />
            <Route {...this.state} path="/themes-install" component={InstallThemePage} />
            <Route {...this.state} exact path="/themes/:id" component={ThemePage} />
            <Route
              {...this.state}
              path="/themes/:name/templates/:template"
              component={TemplatePage}
            />
            <Route {...this.state} path="/settings" component={SettingsPage} />
            <Route {...this.state} path="/data" component={DataPage} />
          </div>
        </Router>
      </UserContext.Provider>
    );
  }
}

document.addEventListener('DOMContentLoaded', () => {
  WebFont.load({
    google: {families: ['Permanent Marker']},
  });
  ReactModal.setAppElement('#react-root');
  ReactDOM.render(<App />, document.getElementById('react-root'));
});
