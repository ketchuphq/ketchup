import * as m from 'mithril';
import * as API from 'lib/api';
import { default as Route } from 'lib/route';
import * as dateFormat from 'date-fns/format';
import { serialize } from 'lib/params';

const dateHumanFormat = 'MMM Do, h:mma';

export const defaultContent: API.Content = {
  uuid: null,
  text: { type: 'html' },
  key: 'content',
  value: ''
};

const defaultPage: API.Page = {
  uuid: null,
  title: null,
  theme: 'none',
  template: 'html',
  contents: [defaultContent],
  authors: [],
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
    this.title = config.title;
    this.theme = config.theme;
    this.template = config.template;
    this.contents = config.contents || [API.Content.copy(defaultContent)];
    this.timestamps = config.timestamps;
    this.publishedAt = config.publishedAt;
    this.authors = config.authors;
    this.routes = [];
  }

  get defaultRoute() {
    if (this.routes.length > 0) {
      return this.routes[0].path;
    }
  }

  save(): Promise<API.Page> {
    return m.request({
      method: 'POST',
      url: `/api/v1/pages`,
      data: API.Page.copy(this)
    });
  }

  delete(): Promise<API.Page> {
    return m.request({
      method: 'DELETE',
      url: `/api/v1/pages/${this.uuid}`,
    });
  }

  publish(): Promise<API.Page> {
    return m.request({
      method: 'POST',
      url: `/api/v1/pages/${this.uuid}/publish`,
    }).then((page: API.Page) => {
      this.publishedAt = page.publishedAt;
      return page;
    });
  }

  unpublish(): Promise<API.Page> {
    return m.request({
      method: 'POST',
      url: `/api/v1/pages/${this.uuid}/unpublish`,
    }).then((page: API.Page) => {
      this.publishedAt = null;
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
    let t = new Date(parseInt(this.timestamps.createdAt));
    return dateFormat(t, dateHumanFormat);
  }

  get formattedUpdatedAt() {
    if (!this.timestamps) {
      return '';
    }
    let t = new Date(parseInt(this.timestamps.updatedAt));
    return dateFormat(t, dateHumanFormat);
  }

  static get(uuid: string): Promise<Page> {
    return m.request({
      method: 'GET',
      url: `/api/v1/pages/${uuid}`
    })
      .then((data: Page) => {
        return new Page(data);
      });
  }

  static list(filter: API.ListPageRequest_ListPageFilter = 'all'): Promise<Page[]> {
    let q = serialize<API.ListPageRequest>({
      options: { filter: filter }
    });
    return m.request({
      method: 'GET',
      url: `/api/v1/pages?${q}`,
    })
      .then((data: API.ListPageResponse) => {
        if (!data.pages) {
          return [];
        }
        return data.pages.map((el) => new Page(el));
      });
  }
}