import Navigation from 'components/navigation';
import Toaster from 'components/toaster';
// import { getUser } from 'components/auth';
import * as React from 'react';

interface Props {
  className: string;
}

export default class Layout extends React.PureComponent<Props> {
  constructor(props: Props) {
    super(props);
    //     getUser().catch(() => {
    //       if (requestedPath.match('^/admin/?$')) {
    //         m.route.set('/admin/login'); // default path is admin
    //       } else {
    //         m.route.set(`/admin/login?next=${requestedPath}`);
    //       }
    //     });
  }

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
