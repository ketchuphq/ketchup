/// <reference path="../../typings/browser.d.ts" />

let Component = {
  controller: class ComponentController { },
  view: () => m('div', 'hello world!')
};

export let routes: _mithril.MithrilRoutes = {
  '/': Component
};

document.addEventListener('DOMContentLoaded', () => {
  let root = document.getElementById('app');
  m.route.mode = 'pathname';
  m.route(root, '/', routes);
});