export let loading = (show: boolean) => !show ? '' :
  m('.loading', [
    m('div'),
    m('div'),
    m('div'),
    'loading'
  ]);