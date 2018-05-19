import * as React from 'react';
import Toaster from 'components/toaster';

interface State {
  email: string;
  password: string;
  showReset: boolean;
}

export default class LoginPage extends React.Component<{}, State> {
  constructor(props: any, context?: any) {
    super(props, context);
    this.state = {
      email: '',
      password: '',
      showReset: false,
    };
  }

  handleLogin = (ev: React.FormEvent<any>) => {
    ev.preventDefault();
    let data = new FormData();
    data.append('email', this.state.email);
    data.append('password', this.state.password);
    fetch(`/api/v1/login`, {
      method: 'POST',
      credentials: 'same-origin',
      body: data,
    }).then((res) => {
      if (res.status == 200) {
        location.assign('/admin');
      }
    });
  };

  render() {
    return (
      <div className="login">
        <Toaster />
        <div className="login-logo">Ketchup</div>
        <div className="login-box">
          <h1>Log In</h1>
          <form onSubmit={this.handleLogin}>
            <div>
              <input
                type="text"
                placeholder="email"
                onChange={(ev) => {
                  this.setState({email: ev.target.value});
                }}
              />
            </div>
            <div>
              <input
                type="password"
                placeholder="password"
                onChange={(ev) => {
                  this.setState({password: ev.target.value});
                }}
              />
            </div>
            <input type="submit" value="Log In" className="button button--green" />
            <div
              className="button txt-small"
              onClick={() => {
                this.setState({showReset: true});
              }}
            >
              Forgot your password?
            </div>
          </form>
          {this.state.showReset ? (
            <div className="reset">
              <p>You can reset your password in the command-line using the following command:</p>
              <pre>./ketchup users:password youremail@example.com</pre>
            </div>
          ) : (
            ''
          )}
        </div>
      </div>
    );
  }
}
