export interface BaseTemplate {
  name: string;
  engine: string;
}

interface BaseAsset {
  name: string;
}

export interface BaseTheme {
  name: string;
  templates: { [key: string]: BaseTemplate };
  assets: { [key: string]: BaseAsset };
}

export default class Theme implements BaseTheme {
  name: string;
  templates: { [key: string]: BaseTemplate };
  assets: { [key: string]: BaseAsset };

  constructor(config?: BaseTheme) {
    this.templates = {};
    this.assets = {};
    if (config) {
      this.name = config.name;
      this.templates = config.templates || {};
      this.assets = config.assets || {};
    }
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
      .then((data: { themes: BaseTheme[] }) => {
        if (!data.themes) {
          return [];
        }
        return data.themes.map((el) => new Theme(el));
      });
  }
}