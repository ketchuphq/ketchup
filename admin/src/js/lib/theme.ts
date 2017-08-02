import * as API from 'lib/api';
import * as m from 'mithril';

export default class Theme extends API.Theme {
  ref: string;

  constructor(config?: API.Theme, ref?: string) {
    super();
    API.Theme.copy(config, this);
    this.templates = this.templates || {};
    this.assets = this.assets || {};
    this.ref = ref;
  }

  getTemplate(name: string): API.ThemeTemplate {
    return this.templates[name];
  }

  static get(name: string): Promise<Theme> {
    return m.request({
      method: 'GET',
      url: `/api/v1/themes/${name}`
    })
      .then((data: API.GetThemeResponse) => {
        return new Theme(data.theme, data.ref);
      });
  }

  static getFullTemplate(name: string, template: string): Promise<API.ThemeTemplate> {
    return m.request({
      method: 'GET',
      url: `/api/v1/themes/${name}/templates/${template}`
    });
  }

  static getAll(): Promise<API.Registry> {
    return m.request({
      method: 'GET',
      url: '/api/v1/theme-registry'
    });
  }

  static install(p: API.Package): Promise<API.Registry> {
    return m.request({
      method: 'POST',
      url: '/api/v1/theme-install',
      background: true,
      data: {
        package: p.name // different id?
      }
    });
  }

  static list(): Promise<Theme[]> {
    return m.request({
      method: 'GET',
      url: '/api/v1/themes'
    })
      .then((data: { themes: API.ThemeTemplate[] }) => {
        if (!data.themes) {
          return [];
        }
        for (let i = 0; i < data.themes.length; i++) {
          if (data.themes[i].name == 'none') {
            let none = data.themes.splice(i, 1)[0];
            data.themes.push(none);
            break;
          }
        }

        return data.themes.map((el) => new Theme(el));
      });
  }
}
