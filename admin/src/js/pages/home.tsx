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
      </header>
      <h2>Welcome to Ketchup.</h2>
      <Button
        class='button--green-2 button--center'
        href='/admin/compose'
      >
        Write a new post &rarr;
      </Button>
      <p>
        <a href='https://ketchuphq.com/docs'>
          Learn more
        </a> &#8599;
      </p>
    </div>;
  }
}
