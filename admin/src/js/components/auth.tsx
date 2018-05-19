import {Loader} from 'components/loading';
import {get} from 'lib/requests';
import * as React from 'react';
import {Redirect, Route, RouteComponentProps, RouteProps} from 'react-router';

export const UserContext: React.Context<User> = React.createContext(null);

export interface User {
  email: string;
  uuid: string;
}

let cachedUser: User = null;

export let getUser = (force = false): Promise<User> => {
  if (cachedUser && !force) {
    return new Promise<User>((resolve) => {
      resolve(cachedUser);
    });
  }
  return get('/api/v1/user')
    .then((res) => res.json())
    .then((res: User) => {
      if (!res.uuid) {
        throw new Error('not logged in');
      }
      cachedUser = res;
      return res;
    })
    .catch(() => {
      cachedUser = null;
      throw new Error('not logged in');
    });
};

interface PrivateRouteProps extends RouteProps {
  user?: User;
  userLoaded: boolean;
}

export interface PrivateRouteComponentProps<T> extends RouteComponentProps<T> {
  user: User;
}

export class PrivateRoute extends React.Component<PrivateRouteProps> {
  render() {
    const {user, userLoaded, component: Komponent, ...rest} = this.props;

    if (user) {
      return <Route {...rest} render={(props) => <Komponent user={user} {...props} />} />;
    } else if (!userLoaded) {
      return (
        <Route
          {...rest}
          render={() => (
            <div style={{width: '100%'}}>
              <Loader show />
            </div>
          )}
        />
      );
    } else {
      return (
        <Route
          {...rest}
          render={(props) => <Redirect to={{pathname: '/login', state: {from: props.location}}} />}
        />
      );
    }
  }
}
