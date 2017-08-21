import * as m from 'mithril';
import msx from 'lib/msx';
import { loading } from 'components/loading';
import Theme from 'lib/theme';
import { BaseComponent } from 'components/auth';

type ThemePickerCallback = (theme: string, template: string) => void;

interface ThemePickerAttrs {
  theme: string;
  template: string;
  callback: ThemePickerCallback;
}

export default class ThemePickerComponent extends BaseComponent<ThemePickerAttrs> {
  callback: ThemePickerCallback;
  themes: Theme[];
  ready: boolean;
  templates: string[];
  selectedTheme: string;
  selectedTemplate: string;

  constructor(v: m.CVnode<ThemePickerAttrs>) {
    super(v);
    this.callback = v.attrs.callback;
    this.themes = [];
    this.templates = [];
    this.ready = false;

    this.selectedTheme = v.attrs.theme;
    this.selectedTemplate = v.attrs.template;

    Theme.list()
      .then((themes) => {
        this.themes = themes;
        return this.selectTheme(v.attrs.theme, true);
      })
      .then(() => this.selectedTemplate = v.attrs.template)
      .then(() => this.ready = true)
      .catch(() => this.ready = true);
  }

  selectTheme(name: string, initial = false) {
    this.selectedTheme = name;
    return Theme.get(this.selectedTheme)
      .then((t) => {
        let templates = Object.keys(t.templates).sort();
        this.templates = templates;
        this.selectedTemplate = templates[0];
        if (!initial) {
          this.callback(this.selectedTheme, this.selectedTemplate);
        }
      });
  }

  selectTemplate(template: string) {
    this.selectedTemplate = template;
    this.callback(this.selectedTheme, this.selectedTemplate);
  }

  view() {
    if (!this.ready) {
      return loading(true);
    }
    return <div class='theme-picker'>
      <div class='control'>
        <div class='label'>Theme</div>
        <select
          oncreate={(v: m.VnodeDOM<any, any>) => {
            (v.dom as HTMLSelectElement).value = this.selectedTheme;
          }}
          onchange={(e: Event) => {
            let target = e.target as HTMLInputElement;
            this.selectTheme(target.value);
          }}
        >
          {this.themes.map((theme: Theme) =>
          <option value={theme.name}>
            {theme.name}
          </option>)}
        </select>
      </div>

      <div class='control'>
        <div class='label'>Template</div>
        <select
          oncreate={(v: m.VnodeDOM<any, any>) => {
            (v.dom as HTMLSelectElement).value = this.selectedTemplate;
          }}
          onchange={(e: Event) => {
            let target = e.target as HTMLInputElement;
            this.selectTemplate(target.value);
          }}
        >
        {this.templates.map((template: string) =>
          <option value={template} key={template}>
            {template}
          </option>)}
        )}
        </select>
      </div>
    </div>;
  }
}
