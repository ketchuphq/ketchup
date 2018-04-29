import * as React from 'react';
// import { MustAuthController } from 'components/auth';
import Button from 'components/button';
import Layout from 'components/layout';

export default class HomePage extends React.PureComponent {
  render() {
    return (
      <Layout className="home">
        <header>
          <img src="/admin/images/k.png" />
        </header>
        <h2>Welcome to Ketchup.</h2>
        <Button className="button--green-2 button--center" href="/compose">
          Write a new post &rarr;
        </Button>
        <p>
          <a href="https://ketchuphq.com/docs">Learn more</a> &#8599;
        </p>
      </Layout>
    );
  }
}
