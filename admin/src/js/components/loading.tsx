import * as React from 'react';

export const Loader: React.SFC<{show: boolean; small?: boolean}> = (props) =>
  !props.show ? null : (
    <div className={props.small ? 'loader' : 'loader loader-large'}>
      <div className="loading0" />
      <div className="loading1" />
      <div className="loading2" />
    </div>
  );

export const LoadingTable: React.SFC = () => <div className="table table-loading" />;
