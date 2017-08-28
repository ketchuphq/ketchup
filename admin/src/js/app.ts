import * as m from 'mithril';

import Layout from 'components/layout';
import HomePage from 'pages/home';
// import RoutesPage from 'pages/routes';
import PagesPage from 'pages/pages';
import PagePage from 'pages/page';
import LoginPage from 'pages/login';
import ThemePage from 'pages/theme';
import ThemesPage from 'pages/themes';
import TemplatePage from 'pages/template';
import SettingsPage from 'pages/settings';
import DataPage from 'pages/data';
import InstallThemePage from 'pages/install-theme';

import * as WebFont from 'webfontloader';


export let routes: m.RouteDefs = {
  '/admin': Layout(HomePage),
  // '/admin/routes': RoutesPage,
  '/admin/pages': Layout(PagesPage),
  '/admin/pages/:id': Layout(PagePage),
  '/admin/compose': Layout(PagePage),
  '/admin/themes': Layout(ThemesPage),
  '/admin/themes/:name': Layout(ThemePage),
  '/admin/themes/:name/templates/:template...': Layout(TemplatePage),
  '/admin/themes-install': Layout(InstallThemePage),
  '/admin/login': LoginPage,
  '/admin/settings': Layout(SettingsPage),
  '/admin/data': Layout(DataPage)
};

document.addEventListener('DOMContentLoaded', () => {
  WebFont.load({
    google: { families: ['Permanent Marker'] }
  });
  let root = document.getElementById('app');
  m.route.prefix('');
  m.route(root, '/admin', routes);
});
