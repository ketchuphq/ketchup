import * as API from 'lib/api';
import {post} from 'lib/requests';

export default class Route implements API.Route {
  uuid: string;
  path: string;
  file?: string;
  pageUuid?: string;

  constructor(config?: API.Route) {
    if (config) {
      this.uuid = config.uuid;
      this.path = config.path;
      this.file = config.file;
      this.pageUuid = config.pageUuid;
    }
  }

  save(): Promise<any> {
    return post('/api/v1/routes', this);
  }

  static format(s: string) {
    if (s == null) {
      return '';
    }
    s = s
      .toLowerCase()
      .replace(/[^a-zA-Z0-9\/]+/gi, '-')
      .replace(/^-+/, '')
      .replace(/-+$/, '')
      .replace(/\/\/+/, '/');
    if (s.length > 0 && s[0] != '/') {
      s = '/' + s;
    }
    return s;
  }

  static list(): Promise<Route[]> {
    return fetch('/api/v1/routes', {
      method: 'GET',
      credentials: 'same-origin',
    })
      .then((res) => res.json())
      .then((data: {routes: API.Route[]}) => {
        if (!data.routes) {
          return [];
        }
        return data.routes.map((el) => new Route(el));
      });
  }

  static getRoutes(page: API.Page): Promise<Route[]> {
    return fetch(`/api/v1/pages/${page.uuid}/routes`)
      .then((res) => res.json())
      .then((res: {routes: API.Route[]}) => {
        return res.routes.map((r) => new Route(r));
      });
  }

  static saveRoutes(page: API.Page, routes: API.Route[]) {
    return post(`/api/v1/pages/${page.uuid}/routes`, {routes});
  }

  static saveList(routes: Route[], pageUUID: string) {
    let chain: Promise<void> = null;
    routes.map((r) => {
      r.pageUuid = pageUUID;
      if (!chain) {
        chain = r.save();
      } else {
        chain.then(() => r.save());
      }
    });
  }
}
