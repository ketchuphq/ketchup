import * as m from 'mithril';
import * as Toaster from 'components/toaster';
import { AuthController, User } from 'components/auth';

export default class LoginPage extends AuthController {
  email: string;
  password: string;
  showReset: boolean;

  constructor() {
    super();
    this._userPromise
      .then((user?: User) => {
        if (!!user) {
          m.route.set('/admin');
        }
      });
    this.email = '';
    this.password = '';
    this.showReset = false;
  }

  login() {
    let data = new FormData();
    data.append('email', this.email);
    data.append('password', this.password);
    m.request({
      method: 'POST',
      url: '/api/v1/login',
      serialize: (x) => x,
      data: data
    }).then(() => {
      location.reload();
    });
  }

  view() {
    return m('.login',
      Toaster.render(),
      m('.login-logo', 'Ketchup'),
      m('.login-box',
        m('h1', 'Log In'),
        m('form', {
          onsubmit: (event: Event) => {
            event.preventDefault();
            this.login();
          }
        }, [
            m('div',
              m('input[type=text]', {
                placeholder: 'email',
                onchange: m.withAttr('value', (val) => this.email = val)
              }),
            ),
            m('div',
              m('input[type=password]', {
                placeholder: 'password',
                onchange: m.withAttr('value', (val) => this.password = val)
              })
            ),
            m('button.button.button--green', 'Log In'),
            m('.button.txt-small', {
              onclick: () => this.showReset = true
            }, 'Forgot your password?')
          ]),
        !this.showReset ? '' : m('.reset',
          m('p', 'You can reset your password in the command-line using the following command:'),
          m('pre', './ketchup users:password youremail@gmail.com')
        )
      )
    );
  }
}
