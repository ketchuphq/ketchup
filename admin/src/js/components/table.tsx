import * as m from 'mithril';
import msx from 'lib/msx';
import { BaseComponent } from 'components/auth';

export class Table extends BaseComponent {
  view(v: m.CVnode<any>) {
    return <div class='table'>{v.children}</div>;
  }
}

export class Row extends BaseComponent {
  view(v: m.CVnode<any>) {
    return <div class='tr'>{v.children}</div>;
  }
}

interface LinkRowProps {
  href: string;
  link?: boolean;
}

export class LinkRow extends BaseComponent<LinkRowProps> {
  view(v: m.CVnode<LinkRowProps>) {
    let oncreate = this.props.link ? m.route.link : () => {};
    return (
      <a class='tr' href={this.props.href} oncreate={oncreate}>
        {v.children}
      </a>
    );
  }
}
