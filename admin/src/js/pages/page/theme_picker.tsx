import * as React from 'react';
import {Loader} from 'components/loading';
import Theme from 'lib/theme';
import * as Page from 'lib/page';
import * as API from 'lib/api';

interface Props {
  store: Page.Store;
}

interface State {
  themes: Theme[];
  ready: boolean;
  templates: string[];

  selectedTheme?: API.Theme;
  selectedTemplate?: string;
}

export default class ThemePickerComponent extends React.Component<Props, State> {
  // themes fetched with List are incomplete, so we fetch themes individually,
  // and cached them here
  fullThemes: {[themeName: string]: API.Theme};

  constructor(props: Props) {
    super(props);
    this.state = {
      themes: [],
      templates: [],
      ready: false,
    };
    this.fullThemes = {};
  }

  componentDidMount() {
    const store = this.props.store;

    // get the selected theme
    let getPromise = this.selectTheme(store.page.theme, store.page.template);

    // list all the themes to populate the dropdown
    let listPromise = Theme.list().then((themes) => {
      this.setState({themes});
    });

    // mark as ready when both get and list are done
    Promise.all([getPromise, listPromise]).then(
      () => this.setState({ready: true}),
      () => this.setState({ready: true})
    );

    // subscribe to page changes
    store.subscribe('theme-picker', (page) => {
      this.selectTheme(page.theme, page.template);
    });
  }

  componentWillUnmount() {
    this.props.store.unsubscribe('theme-picker');
  }

  selectTheme(themeName: string, templateName?: string) {
    const page = this.props.store.page;
    return Promise.resolve(this.fullThemes[themeName])
      .then((theme) => {
        if (theme) {
          return theme;
        }
        return Theme.get(themeName).then(({theme}) => {
          this.fullThemes[theme.name] = theme;
          return theme;
        });
      })
      .then((theme) => {
        let templates = Object.keys(theme.templates).sort();
        let template = Theme.getTemplate(theme, templateName || templates[0]);
        this.setState({
          selectedTheme: theme,
          selectedTemplate: template.name,
          templates,
        });

        // update page if different
        if (page.theme != themeName || page.template != templateName) {
          this.props.store.setThemeTemplate(theme, template);
        }
      });
  }

  render() {
    if (!this.state.ready) {
      return <Loader show={true} />;
    }
    return (
      <div className="theme-picker">
        <div className="control">
          <div className="label">Theme</div>
          <select
            defaultValue={this.state.selectedTheme.name}
            onChange={(e: React.ChangeEvent<HTMLSelectElement>) => {
              this.selectTheme(e.target.value);
            }}
          >
            {this.state.themes.map((theme: Theme) => (
              <option key={theme.name} value={theme.name}>
                {theme.name}
              </option>
            ))}
          </select>
        </div>

        <div className="control">
          <div className="label">Template</div>
          <select
            defaultValue={this.state.selectedTemplate}
            onChange={(e: React.ChangeEvent<HTMLSelectElement>) => {
              this.selectTheme(this.state.selectedTheme.name, e.target.value);
            }}
          >
            {this.state.templates.map((template: string) => (
              <option key={template} value={template}>
                {template}
              </option>
            ))}
            )}
          </select>
        </div>
      </div>
    );
  }
}
