import Layout from '../components/layout';

export default class HomePage {
  constructor() {}
  static controller = HomePage;
  static view(ctrl: HomePage) {
    return Layout('hai');
  }
}
