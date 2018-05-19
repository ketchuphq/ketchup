import * as React from 'react';
import {Link} from 'react-router-dom';

export const Table: React.SFC<{loading?: boolean}> = (props) => (
  <div className={`table ${props.loading ? 'table-loading' : ''}`}>{props.children}</div>
);

export const Row: React.SFC<{center?: boolean}> = (props) => {
  let k = 'tr';
  if (props.center) {
    k += ' tr--center';
  }
  return <div className={k}>{props.children}</div>;
};

export const LinkRow: React.SFC<{
  href: string;
  link?: boolean;
}> = (props) => {
  if (props.link) {
    return (
      <Link className="tr" to={props.href}>
        {props.children}
      </Link>
    );
  }
  return (
    <a className="tr" href={props.href}>
      {' '}
      {props.children}{' '}
    </a>
  );
};
