import * as m from 'mithril';

import Layout from 'components/layout';
import HomePage from 'pages/home';
// import RoutesPage from 'pages/routes';
import PagesPage from 'pages/pages';
// import PagePage from 'pages/page';
import LoginPage from 'pages/login';
import ThemePage from 'pages/theme';
import ThemesPage from 'pages/themes';
import TemplatePage from 'pages/template';
import SettingsPage from 'pages/settings';
// import InstallThemePage from 'pages/install-theme';

import * as WebFont from 'webfontloader';


export let routes: Mithril.RouteDefs = {
  '/admin': Layout(HomePage),
  // '/admin/routes': RoutesPage,
  '/admin/pages': Layout(PagesPage),

  // '/admin/pages/:id': PagePage,
  // '/admin/compose': PagePage,
  '/admin/themes': Layout(ThemesPage),
  // '/admin/themes/install': InstallThemePage,
  '/admin/themes/:name': Layout(ThemePage),
  '/admin/themes/:name/templates/:template': Layout(TemplatePage),
  '/admin/login': LoginPage,
  '/admin/settings': Layout(SettingsPage)
};

document.addEventListener('DOMContentLoaded', () => {
  WebFont.load({
    google: { families: ['Permanent Marker'] }
  });
  let root = document.getElementById('app');
  m.route.prefix('');
  m.route(root, '/admin', routes);
});
