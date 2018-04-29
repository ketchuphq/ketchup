import * as React from 'react';

export const Loader: React.SFC<{show: boolean}> = (props) =>
  !props.show ? null : (
    <div className="loader loader-large">
      <div className="loading0" />
      <div className="loading1" />
      <div className="loading2" />
    </div>
  );
