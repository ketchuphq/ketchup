import * as m from 'mithril';
import * as API from 'lib/api';
import { default as Route } from 'lib/route';
import * as dateFormat from 'date-fns/format';

const dateHumanFormat = 'MMM Do, h:mma';

let defaultPage: API.Page = {
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
  timestamps: {
    createdAt: null,
    updatedAt: null
  },
  publishedAt: null
};

export default class Page implements API.Page {
  uuid: string;
  name: string;
  theme: string;
  template: string;
  contents: API.Content[];
  timestamps: API.Timestamp;
  publishedAt: string;

  constructor(config?: API.Page) {
    config = config || defaultPage;
    this.uuid = config.uuid;
    this.name = config.name;
    this.theme = config.theme;
    this.template = config.template;
    this.contents = config.contents;
    this.timestamps = config.timestamps;
    this.publishedAt = config.publishedAt;
  }

  save(): Mithril.Promise<API.Page> {
    return m.request({
      method: 'POST',
      url: `/api/v1/pages`,
      data: this
    });
  }

  publish(): Mithril.Promise<API.Page> {
    return m.request({
      method: 'POST',
      url: `/api/v1/pages/${this.uuid}/publish`,
    }).then((page: API.Page) => {
      this.publishedAt = page.publishedAt;
      return page;
    });
  }

  getRoutes() {
    return m.request({
      method: 'GET',
      url: `/api/v1/pages/${this.uuid}/routes`,
    })
      .then((res: { routes: API.Route[] }) =>
        res.routes.map((r) =>
          new Route(r)));
  }

  saveRoutes(routes: API.Route[]) {
    return m.request({
      method: 'POST',
      url: `/api/v1/pages/${this.uuid}/routes`,
      data: {
        routes: routes
      }
    });
  }

  get isPublished() {
    return this.publishedAt != null;
  }

  get formattedCreatedAt() {
    if (!this.timestamps) {
      return '';
    }
    let t = new Date(parseInt(this.timestamps.createdAt) * 1000);
    return dateFormat(t, dateHumanFormat);
  }

  get formattedUpdatedAt() {
    if (!this.timestamps) {
      return '';
    }
    let t = new Date(parseInt(this.timestamps.updatedAt) * 1000);
    return dateFormat(t, dateHumanFormat);
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
      .then((data: { pages: API.Page[] }) => {
        if (!data.pages) {
          return [];
        }
        return data.pages.map((el) => new Page(el));
      });
  }
}