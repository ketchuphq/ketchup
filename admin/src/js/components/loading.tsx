import msx from 'lib/msx';

export const Loader = (
  <div class='loader loader-large'>
    <div class='loading0' />
    <div class='loading1' />
    <div class='loading2' />
  </div>
);

export const loading = (show: boolean) => (!show ? null : Loader);
