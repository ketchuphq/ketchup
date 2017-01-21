import * as Toaster from 'components/toaster';
import { AuthController, User } from 'components/auth';

export default class LoginPage extends AuthController {
  email: Mithril.BasicProperty<string>;
  password: Mithril.BasicProperty<string>;
  showReset: Mithril.BasicProperty<boolean>;

  constructor() {
    super();
    this._userPromise
      .then((user?: User) => {
        if (!!user) {
          m.route('/admin');
        }
      });
    this.email = m.prop('');
    this.password = m.prop('');
    this.showReset = m.prop(false);
  }

  login() {
    let data = new FormData();
    data.append('email', this.email());
    data.append('password', this.password());
    m.request({
      method: 'POST',
      url: '/api/v1/login',
      serialize: (x) => x,
      data: data
    }).then(() => {
      location.reload();
    });
  }

  static controller = LoginPage;
  static view(ctrl: LoginPage) {
    return m('.login',
      Toaster.render(),
      m('.login-logo', 'Ketchup'),
      m('.login-box',
        m('h1', 'Log In'),
        m('form', {
          onsubmit: (event: Event) => {
            event.preventDefault();
            ctrl.login();
          }
        }, [
            m('div',
              m('input[type=text]', {
                placeholder: 'email',
                onchange: m.withAttr('value', ctrl.email)
              }),
            ),
            m('div',
              m('input[type=password]', {
                placeholder: 'password',
                onchange: m.withAttr('value', ctrl.password)
              })
            ),
            m('button.button.button--green', 'Log In'),
            m('.button.small', {
              onclick: () => ctrl.showReset(true)
            }, 'Forgot your password?')
          ]),
        !ctrl.showReset() ? '' : m('.reset',
          m('p', 'You can reset your password in the command-line using the following command:'),
          m('pre', './ketchup users:password me@gmail.com')
        )
      )
    );
  }
}
