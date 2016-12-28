import Layout from 'components/layout';

export default class HomePage {
  constructor() { }
  static controller = HomePage;
  static view(_: HomePage) {
    return Layout(m('h1', 'ketchup'));
  }
}
