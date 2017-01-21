import * as m from 'mithril';
import HomePage from 'pages/home';
import RoutesPage from 'pages/routes';
import PagesPage from 'pages/pages';
import PagePage from 'pages/page';
import LoginPage from 'pages/login';
import ThemePage from 'pages/theme';
import ThemesPage from 'pages/themes';
import TemplatePage from 'pages/template';
import SettingsPage from 'pages/settings';
import InstallThemePage from 'pages/install-theme';

import * as WebFont from 'webfontloader';

export let routes: Mithril.Routes = {
  '/admin': HomePage,
  '/admin/routes': RoutesPage,
  '/admin/pages': PagesPage,

  '/admin/pages/:id': PagePage,
  '/admin/compose': PagePage,
  '/admin/themes': ThemesPage,
  '/admin/themes/install': InstallThemePage,
  '/admin/themes/:name': ThemePage,
  '/admin/themes/:name/templates/:template': TemplatePage,
  '/admin/login': LoginPage,
  '/admin/settings': SettingsPage
};

document.addEventListener('DOMContentLoaded', () => {
  WebFont.load({
    google: { families: ['Permanent Marker'] }
  });
  let root = document.getElementById('app');
  m.route.mode = 'pathname';
  m.route(root, '/admin', routes);

  // root.addEventListener('dragenter', (ev) => {
  //   ev.preventDefault()
  // })

  // root.addEventListener('drop', (ev) => {
  //   ev.preventDefault()
  //   console.log(ev.dataTransfer.files)
  // })

});
