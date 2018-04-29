import * as dateFormat from 'date-fns/format';
import * as API from 'lib/api';
import {serialize} from 'lib/params';
import {get, post} from 'lib/requests';
import GenericStore from 'lib/store';

const dateHumanFormat = 'MMM Do, h:mma';

export const defaultContent: API.Content = {
  uuid: null,
  text: {type: 'html'},
  key: 'content',
  value: '',
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
    updatedAt: null,
  },
  publishedAt: null,
};

export class Store extends GenericStore<API.Page> {
  constructor(obj?: API.Page) {
    super(API.Page.copy, obj);
  }

  get page() {
    return this.obj;
  }

  get(uuid: string): Promise<API.Page> {
    return get(`/api/v1/pages/${uuid}`)
      .then((res) => res.json())
      .then((page: API.Page) => this.set(page));
  }

  save(): Promise<API.Page> {
    return post(`/api/v1/pages`, API.Page.copy(this.page))
      .then((res) => res.json())
      .then((page: API.Page) => {
        this.set(page);
        return page;
      });
  }

  // updateContent updates the page contents by iterating over page contents
  // and template placeholders
  setThemeTemplate(theme: API.Theme, template: API.ThemeTemplate, initial = false) {
    // todo: keep old content temporarily for 'undo'
    let contentMap: {[key: string]: API.Content} = {};
    let placeholderContents: API.Content[] = [];
    let placeholderContentMap: {[key: string]: boolean} = {};
    let page = this.page;
    // cache existing content
    (page.contents || []).forEach((c) => {
      // on initial load copy all fields
      if (c.uuid || c.value || initial) {
        contentMap[c.key] = c;
      }
    });

    // load placeholders from templates
    (template.placeholders || []).forEach((p) => {
      if (contentMap[p.key]) {
        Object.keys(p).forEach((k: keyof API.ThemePlaceholder) => {
          if (p[k]) {
            contentMap[p.key][k] = p[k];
          }
        });
        placeholderContents.push(contentMap[p.key]);
      } else {
        placeholderContents.push(API.Content.copy(p, {}));
      }
      placeholderContentMap[p.key] = true;
    });

    // load placeholders from existing page contents
    let pageContents: API.Content[] = [];
    (page.contents || []).forEach((c) => {
      if (c.key == 'content' && template.hideContent) {
        return;
      }
      if (contentMap[c.key] && !placeholderContentMap[c.key]) {
        pageContents.push(c);
      }
    });

    // add default content editor
    if (!contentMap['content'] && !placeholderContentMap['content'] && !template.hideContent) {
      placeholderContents.push(API.Content.copy(defaultContent, {}));
    }

    // hack; replace with immutable?
    this.update((page) => {
      page.theme = theme.name;
      page.template = template.name;
      page.contents = pageContents.concat(placeholderContents);
    });
  }

  publish(): Promise<API.Page> {
    return post(`/api/v1/pages/${this.page.uuid}/publish`)
      .then((res) => res.json())
      .then((p: API.Page) => {
        return this.update((page) => {
          page.publishedAt = p.publishedAt;
        });
      });
  }

  unpublish(): Promise<API.Page> {
    return post(`/api/v1/pages/${this.page.uuid}/unpublish`)
      .then((res) => res.json())
      .then((_: API.Page) => {
        return this.update((page) => {
          page.publishedAt = null;
        });
      });
  }
}

export function newPage(page?: API.Page): API.Page {
  page = API.Page.copy(page || defaultPage);
  page.contents = [API.Content.copy(defaultContent)];
  return page;
}

export function deletePage(page: API.Page): Promise<Response> {
  return fetch(`/api/v1/pages/${page.uuid}`, {
    method: 'DELETE',
    credentials: 'same-origin',
  });
}

export function isPublished(page: API.Page) {
  return page.publishedAt != null;
}

export function formattedCreatedAt(page: API.Page) {
  if (!page.timestamps) {
    return '';
  }
  let t = new Date(parseInt(page.timestamps.createdAt));
  return dateFormat(t, dateHumanFormat);
}

export function formattedUpdatedAt(page: API.Page) {
  if (!page.timestamps) {
    return '';
  }
  let t = new Date(parseInt(page.timestamps.updatedAt));
  return dateFormat(t, dateHumanFormat);
}

export function list(filter: API.ListPageRequest_ListPageFilter = 'all'): Promise<API.Page[]> {
  let q = serialize<API.ListPageRequest>({
    options: {filter: filter},
  });
  return get(`/api/v1/pages?${q}`)
    .then((res) => res.json())
    .then((data: API.ListPageResponse) => {
      if (!data.pages) {
        return [];
      }
      return data.pages.map((el) => newPage(el));
    });
}
