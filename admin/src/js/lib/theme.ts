import * as API from 'lib/api';
import {get, post} from 'lib/requests';

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

  static checkForUpdates(name: string): Promise<API.CheckThemeForUpdateResponse> {
    return get(`/api/v1/themes/${name}/updates`).then((res) => res.json());
  }

  static get(name: string): Promise<Theme> {
    return get(`/api/v1/themes/${name}`)
      .then((res) => res.json())
      .then((data: API.GetThemeResponse) => {
        return new Theme(data.theme, data.ref);
      });
  }

  static getFullTemplate(name: string, template: string): Promise<API.ThemeTemplate> {
    return get(`/api/v1/themes/${name}/templates/${template}`).then((res) => res.json());
  }

  static getAll(): Promise<API.Registry> {
    return get('/api/v1/theme-registry').then((res) => res.json());
  }

  static install(p: API.Package): Promise<API.Registry> {
    let data: API.InstallThemeRequest = {
      name: p.name,
      vcsUrl: p.vcsUrl,
    };

    return post('/api/v1/theme-install', data).then((res) => res.json());
  }

  static list(): Promise<Theme[]> {
    return get('/api/v1/themes')
      .then((res) => res.json())
      .then((data: {themes: API.ThemeTemplate[]}) => {
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
