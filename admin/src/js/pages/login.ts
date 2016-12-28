import * as Toaster from 'components/toaster';
import { AuthController, User } from 'components/auth';

export default class LoginPage extends AuthController {
  email: Mithril.BasicProperty<string>;
  password: Mithril.BasicProperty<string>;

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
      m('h1', 'Login'),
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
          m('button.button.button--green', 'Log In')
        ])
    );
  }
}
