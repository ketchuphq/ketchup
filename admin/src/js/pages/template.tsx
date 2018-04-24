import Layout from 'components/layout';
import {Loader} from 'components/loading';
import * as API from 'lib/api';
import Theme from 'lib/theme';
import * as React from 'react';
import {RouteComponentProps} from 'react-router';
import {Link} from 'react-router-dom';

type Props = RouteComponentProps<{name: string; template: string}>;
interface State {
  template?: API.ThemeTemplate;
}

export default class TemplatePage extends React.Component<Props, State> {
  preRef: React.RefObject<HTMLPreElement>;

  constructor(props: any) {
    super(props);
    this.preRef = React.createRef();
    this.state = {};
  }

  componentDidMount() {
    let {name, template} = this.props.match.params;
    Theme.getFullTemplate(name, template).then((template) => {
      this.setState({template});
    });
  }

  componentDidUpdate(_: Props, prevState: State) {
    if (!prevState.template || prevState.template.data != this.state.template.data) {
      this.colorize(this.preRef.current);
    }
  }

  async colorize(el: HTMLElement) {
    const hljsImports = await Promise.all([
      import(/* webpackChunkName: "hljs" */ 'highlight.js'),
      import(/* webpackChunkName: "hljs" */ 'highlight.js/lib/languages/xml'),
      import(/* webpackChunkName: "hljs" */ 'highlight.js/styles/rainbow.css'),
    ]);
    let hljs = hljsImports[0];
    hljs.highlightBlock(el);
  }

  render() {
    let {name, template} = this.props.match.params;
    return (
      <Layout className="template">
        <header>
          <h1>
            <Link to="/themes">Themes</Link> &rsaquo;{' '}
            <Link to={`/themes/${name}`} className="unbold">
              {name}
            </Link>{' '}
            &rsaquo; <span className="unbold">{template}</span>
          </h1>
        </header>
        {!!this.state.template ? (
          <Placeholders
            hideContent={this.state.template.hideContent}
            placeholders={this.state.template.placeholders}
          />
        ) : (
          <Loader show />
        )}

        <h2>Template</h2>
        {!!this.state.template ? (
          <pre ref={this.preRef}>{this.state.template.data}</pre>
        ) : (
          <Loader show />
        )}
      </Layout>
    );
  }
}

let Placeholders: React.SFC<{hideContent: boolean; placeholders: API.ThemePlaceholder[]}> = ({
  hideContent,
  placeholders,
}) => {
  let content = (placeholders || []).map((placeholder) => {
    if (placeholder.key != 'content') {
      return (
        <div key={placeholder.key} className="tr">
          {placeholder.key}
        </div>
      );
    }
  });
  if (!hideContent) {
    content.push(
      <div key="content" className="tr">
        content
      </div>
    );
  }
  if (content.length == 0) {
    return null;
  }
  return (
    <div>
      <h2>Fields</h2>
      <div className="table">{content}</div>
    </div>
  );
};
