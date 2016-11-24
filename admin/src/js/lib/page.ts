import * as m from 'mithril';
import { BaseRoute, default as Route } from 'lib/route';

interface BasePage {
  uuid: string;
  name: string;
  theme: string;
  template: string;
  contents: Content[];
}

export interface Content {
  uuid?: string;
  contentType: 'html' | 'markdown';
  key: string;
  value: string;
}

let defaultPage: BasePage = {
  uuid: null,
  name: null,
  theme: 'basic',
  template: 'index.html',
  contents: [{
    uuid: null,
    contentType: 'html',
    key: 'content',
    value: ''
  }],
};

export default class Page implements BasePage {
  uuid: string;
  name: string;
  theme: string;
  template: string;
  contents: Content[];

  constructor(config?: BasePage) {
    config = config || defaultPage;
    this.uuid = config.uuid;
    this.name = config.name;
    this.theme = config.theme;
    this.template = config.template;
    this.contents = config.contents;
  }

  save(): Mithril.Promise<Page> {
    return m.request({
      method: 'POST',
      url: `/api/v1/pages`,
      data: this
    });
  }

  getRoutes() {
    return m.request({
      method: 'GET',
      url: `/api/v1/pages/${this.uuid}/routes`,
    })
      .then((res: { routes: BaseRoute[] }) =>
        res.routes.map((r) =>
          new Route(r)));
  }

  static get(uuid: string): Mithril.Promise<Page> {
    return m.request({
      method: 'GET',
      url: `/api/v1/pages/${uuid}`
    })
      .then((data: Page) => {
        return new Page(data);
      });
  }

  static list(): Mithril.Promise<Page[]> {
    return m.request({
      method: 'GET',
      url: '/api/v1/pages'
    })
      .then((data: { pages: BasePage[] }) => {
        if (!data.pages) {
          return [];
        }
        return data.pages.map((el) => new Page(el));
      });
  }
}