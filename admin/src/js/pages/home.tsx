import msx from 'lib/msx';
import { MustAuthController } from 'components/auth';
import Button from 'components/button';

export default class HomePage extends MustAuthController {
  constructor() {
    super();
  }
  static oninit() {
    new HomePage();
  }
  static view() {
    return <div class='home'>
      <header>
        <img src='/admin/images/k.png' />
        <h1>Ketchup!</h1>
      </header>
      <p>Welcome to Ketchup.</p>
      <Button
        class='button--green button--center'
        href='/admin/compose'
      >
        Compose
      </Button>
    </div>;
  }
}
