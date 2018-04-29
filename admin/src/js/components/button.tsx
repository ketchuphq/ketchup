import * as React from 'react';
import {Link} from 'react-router-dom';
import {Loader} from 'components/loading';

interface Props {
  handler?: () => Promise<any>;
  onClick?: () => any;

  className?: string;
  id?: string;
  href?: string;
}

interface State {
  loading: boolean;
}

export default class Button extends React.Component<Props, State> {
  handler: () => void;

  constructor(props: Props) {
    super(props, null);
    this.state = {
      loading: false,
    };
    this.handler = () => {
      if (!this.props.handler) {
        return;
      }
      this.setState({loading: true});
      this.props.handler().then(
        () => this.setState({loading: false}),
        () => {
          this.setState({loading: false});
          return Promise.reject(null);
        }
      );
    };
  }

  render() {
    let className = `button ${this.props.className}`;
    if (this.state.loading) {
      className += ' button--loading';
    }
    if (!this.props.href) {
      return (
        <a className={className} id={this.props.id} onClick={this.handler}>
          <div className="button__inner">{this.props.children}</div>
          {this.state.loading && <Loader show small />}
        </a>
      );
    }
    return (
      <Link className={className} to={this.props.href} id={this.props.id} onClick={this.handler}>
        <div className="button__inner">{this.props.children}</div>
        {this.state.loading && <Loader show small />}
      </Link>
    );
  }
}
