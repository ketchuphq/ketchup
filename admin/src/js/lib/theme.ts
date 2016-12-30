import * as API from 'lib/api';

export default class Theme implements API.Theme {
  name: string;
  templates: { [key: string]: API.ThemeTemplate };
  assets: { [key: string]: API.ThemeAsset };

  constructor(config?: API.Theme) {
    this.templates = {};
    this.assets = {};
    if (config) {
      this.name = config.name;
      this.templates = config.templates || {};
      this.assets = config.assets || {};
    }
  }

  getTemplate(name: string): API.ThemeTemplate {
    return this.templates[name];
  }

  static get(name: string): Mithril.Promise<Theme> {
    return m.request({
      method: 'GET',
      url: `/api/v1/themes/${name}`
    })
      .then((data: Theme) => {
        return new Theme(data);
      });
  }

  static list() {
    return m.request({
      method: 'GET',
      url: '/api/v1/themes'
    })
      .then((data: { themes: API.ThemeTemplate[] }) => {
        if (!data.themes) {
          return [];
        }
        return data.themes.map((el) => new Theme(el));
      });
  }
}