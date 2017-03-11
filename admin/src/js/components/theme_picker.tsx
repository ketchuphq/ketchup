import msx from 'lib/msx';
import { loading } from 'components/loading';
import Theme from 'lib/theme';

type ThemePickerCallback = (theme: string, template: string) => void;

interface ThemePickerAttrs {
  theme: string;
  template: string;
  callback: ThemePickerCallback;
}

export default class ThemePickerComponent {
  callback: ThemePickerCallback;
  themes: Theme[];
  ready: boolean;
  templates: string[];
  selectedTheme: string;
  selectedTemplate: string;

  constructor(attrs: ThemePickerAttrs) {
    this.callback = attrs.callback;
    this.themes = [];
    this.templates = [];
    this.ready = false;

    this.selectedTheme = attrs.theme;
    this.selectedTemplate = attrs.template;

    Theme.list()
      .then((themes) => {
        this.themes = themes;
        return this.selectTheme(attrs.theme, true);
      })
      .then(() => this.selectedTemplate = attrs.template)
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

  static oninit(v: Mithril.Vnode<ThemePickerAttrs, ThemePickerComponent>) {
    v.state = new ThemePickerComponent(v.attrs);
  };

  static view(v: Mithril.Vnode<ThemePickerAttrs, ThemePickerComponent>) {
    let ctrl = v.state;
    if (!ctrl.ready) {
      return loading(true);
    }
    return <div class='theme-picker'>
      <div class='control'>
        <div class='label'>Theme</div>
        <select
          oncreate={(v: Mithril.VnodeDOM<any, any>) => {
            (v.dom as HTMLSelectElement).value = ctrl.selectedTheme;
          }}
          onchange={(e: Event) => {
            let target = e.target as HTMLInputElement;
            ctrl.selectTheme(target.value);
          }}
        >
          {ctrl.themes.map((theme: Theme) =>
          <option value={theme.name}>
            {theme.name}
          </option>)}
        </select>
      </div>

      <div class='control'>
        <div class='label'>Template</div>
        <select
          oncreate={(el: HTMLSelectElement, isInitialized: boolean) => {
            if (!isInitialized) {
              el.value = ctrl.selectedTemplate;
            }
          }}
          onchange={(e: Event) => {
            let target = e.target as HTMLInputElement;
            ctrl.selectTemplate(target.value);
          }}
        >
        {ctrl.templates.map((template: string) =>
          <option value={template} key={template}>
            {template}
          </option>)}
        )}
        </select>
      </div>
    </div>;
  }
}