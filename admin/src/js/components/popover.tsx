import * as React from 'react';

interface Props {
  visible: boolean;
}

export default class Popover extends React.PureComponent<Props> {
  render() {
    let klass = 'popover-outer';
    if (!this.props.visible) {
      klass += ' popover-outer-hidden';
    }
    return (
      <div className={klass}>
        <div className="popover">{this.props.children}</div>
      </div>
    );
  }
}
