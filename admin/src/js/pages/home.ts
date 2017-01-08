import Layout from 'components/layout';
import { MustAuthController } from 'components/auth';

export default class HomePage extends MustAuthController {
  constructor() {
    super();
  }
  static controller = HomePage;
  static view(_: HomePage) {
    return Layout(m('h1', 'ketchup'));
  }
}
