import Navigation from 'components/navigation';
import Toaster from 'components/toaster';
// import { getUser } from 'components/auth';
import * as React from 'react';

interface Props {
  className: string;
}

export default class Layout extends React.PureComponent<Props> {
  render() {
    return (
      <div className="container">
        <Navigation />
        <Toaster />
        <div className={`container__body ${this.props.className}`}>{this.props.children}</div>
      </div>
    );
  }
}
