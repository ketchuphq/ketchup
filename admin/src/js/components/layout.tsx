import Navigation from 'components/navigation';
import * as React from 'react';

interface Props {
  className: string;
}

export default class Layout extends React.PureComponent<Props> {
  render() {
    return (
      <div className="container">
        <Navigation />
        <div className={`container__body ${this.props.className}`}>{this.props.children}</div>
      </div>
    );
  }
}
