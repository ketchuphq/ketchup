export interface BaseRoute {
  uuid: string;
  path: string;

  file?: string;
  page_uuid?: string;
  delegate?: string;
}

export default class Route implements BaseRoute {
  uuid: string;
  path: string;
  file?: string;
  page_uuid?: string;
  delegate?: string;

  constructor(config?: BaseRoute) {
    if (config) {
      this.uuid = config.uuid;
      this.path = config.path;
      this.file = config.file;
      this.page_uuid = config.page_uuid;
      this.delegate = config.delegate;
    }
  }

  save() {
    return m.request({
      method: 'POST',
      url: '/api/v1/routes',
      data: this as BaseRoute
    });
  }

  static list() {
    return m.request({
      method: 'GET',
      url: '/api/v1/routes'
    })
      .then((data: { routes: BaseRoute[] }) => {
        if (!data.routes) {
          return [];
        }
        return data.routes.map((el) => new Route(el));
      });
  }

  static saveList(routes: Route[], pageUUID: string) {
    let chain: Mithril.Promise<void> = null;
    routes.map((r) => {
      r.page_uuid = pageUUID;
      if (!chain) {
        chain = r.save();
      } else {
        chain.then(() => r.save());
      }
    });
  }
}