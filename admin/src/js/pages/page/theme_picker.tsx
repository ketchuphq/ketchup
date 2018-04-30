import * as React from 'react';
import {Loader} from 'components/loading';
import Theme from 'lib/theme';
import * as Page from 'lib/page';

interface Props {
  store: Page.Store;
}

interface State {
  themes: Theme[];
  ready: boolean;
  templates: string[];

  selectedTheme: string;
  selectedTemplate: string;
}

export default class ThemePickerComponent extends React.Component<Props, State> {
  themeSelectRef: React.RefObject<HTMLSelectElement>;
  templateSelectRef: React.RefObject<HTMLSelectElement>;

  constructor(props: Props) {
    super(props);
    const store = props.store;
    this.state = {
      themes: [],
      templates: [],
      ready: false,
      selectedTheme: store.page.theme,
      selectedTemplate: store.page.template,
    };

    this.themeSelectRef = React.createRef();
    this.templateSelectRef = React.createRef();
  }

  componentDidMount() {
    const store = this.props.store;
    // get all the themes
    Theme.list()
      .then((themes) => {
        this.setState({themes});
        return this.selectTheme(store.page.theme, store.page.template);
      })
      .then(() => {
        this.templateSelectRef.current.value = this.state.selectedTemplate;
        this.themeSelectRef.current.value = this.state.selectedTheme;
      })
      .then(() => this.setState({ready: true}), () => this.setState({ready: true}));

    // subscribe to page changes
    store.subscribe('theme-picker', (page) => {
      this.setState({
        selectedTheme: page.theme,
        selectedTemplate: page.template,
      });
    });
  }

  componentWillUnmount() {
    this.props.store.unsubscribe('theme-picker');
  }

  selectTheme(name: string, template?: string) {
    return Theme.get(name).then(({theme}) => {
      let templates = Object.keys(theme.templates).sort();
      this.setState({templates});
      this.props.store.setThemeTemplate(theme, Theme.getTemplate(theme, template || templates[0]));
    });
  }

  selectTemplate(template: string) {
    this.state.themes.map((theme) => {
      if (theme.name == this.state.selectedTheme) {
        this.props.store.setThemeTemplate(theme, Theme.getTemplate(theme, template));
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
            onChange={(e: React.ChangeEvent<HTMLSelectElement>) => {
              this.selectTemplate(e.target.value);
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
