import * as m from 'mithril';
import * as API from 'lib/api';
import { default as Route } from 'lib/route';
import * as dateFormat from 'date-fns/format';

const dateHumanFormat = 'MMM Do, h:mma';

const defaultContent: API.Content = {
  uuid: null,
  text: { type: 'html' },
  key: 'content',
  value: ''
};

const defaultPage: API.Page = {
  uuid: null,
  name: null,
  theme: 'none',
  template: 'default',
  contents: [defaultContent],
  timestamps: {
    createdAt: null,
    updatedAt: null
  },
  publishedAt: null
};

export default class Page extends API.Page {
  routes: API.Route[];

  constructor(config?: API.Page) {
    super();
    config = config || defaultPage;
    this.uuid = config.uuid;
    this.name = config.name;
    this.theme = config.theme;
    this.template = config.template;
    this.contents = config.contents || [defaultContent];
    this.timestamps = config.timestamps;
    this.publishedAt = config.publishedAt;
    this.routes = [];
  }

  get defaultRoute() {
    if (this.routes.length > 0) {
      return this.routes[0].path;
    }
  }

  updateContent(c: API.Content) {
    for (var i = 0; i < this.contents.length; i++) {
      var element = this.contents[i];
      if (element.key == c.key) {
        API.Content.copy(c, element);
        return;
      }
    }
    this.contents.push(API.Content.copy(c));
  }

  save(): Mithril.Promise<API.Page> {
    return m.request({
      method: 'POST',
      url: `/api/v1/pages`,
      data: API.Page.copy(this)
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
      .then((res: { routes: API.Route[] }) => {
        this.routes = res.routes.map((r) => new Route(r));
        return this.routes;
      });
  }

  saveRoutes() {
    return m.request({
      method: 'POST',
      url: `/api/v1/pages/${this.uuid}/routes`,
      data: {
        routes: this.routes
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