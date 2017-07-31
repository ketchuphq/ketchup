import msx from 'lib/msx';

export let loading = (show: boolean) => !show ? null :
  <div class='loader loader-large'>
    <div class='loading0' />
    <div class='loading1' />
    <div class='loading2' />
  </div>;
