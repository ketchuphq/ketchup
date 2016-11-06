import * as m from 'mithril';
import HomePage from './pages/home';
import RoutesPage from './pages/routes';
import PagesPage from './pages/pages';
import PagePage from './pages/page';

export let routes: Mithril.Routes = {
  '/admin': HomePage,
  '/admin/routes': RoutesPage,
  '/admin/pages': PagesPage,

  '/admin/pages/:id': PagePage,
  '/admin/compose': PagePage
};

document.addEventListener('DOMContentLoaded', () => {
  let root = document.getElementById('app');
  m.route.mode = 'pathname';
  m.route(root, '/admin', routes);
});
