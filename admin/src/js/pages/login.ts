import Layout from 'components/layout';

export default class LoginPage {
  email: Mithril.BasicProperty<string>;
  password: Mithril.BasicProperty<string>;

  constructor() {
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
    });
  }

  static controller = LoginPage;
  static view(ctrl: LoginPage) {
    return Layout(
      m('.login',
        m('h1', 'Login'),
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
        m('.button.button--green', {
          onclick: () => ctrl.login()
        }, 'Log In')
      )
    );
  }
}
