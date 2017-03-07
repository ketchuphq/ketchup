import * as m from 'mithril';
import * as stream from 'mithril/stream';
import * as Toaster from 'components/toaster';
import { AuthController, User } from 'components/auth';

let _: Mithril.Component<{}, LoginPage> = LoginPage;

export default class LoginPage extends AuthController {
  email: Mithril.Stream<string>;
  password: Mithril.Stream<string>;
  showReset: Mithril.Stream<boolean>;

  constructor() {
    super();
    this._userPromise
      .then((user?: User) => {
        if (!!user) {
          m.route.set('/admin');
        }
      });
    this.email = stream('');
    this.password = stream('');
    this.showReset = stream(false);
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

  static oninit(v: Mithril.Vnode<{}, LoginPage>) {
    v.state = new LoginPage();
  }
  static view(v: Mithril.Vnode<{}, LoginPage>) {
    let ctrl = v.state;
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
