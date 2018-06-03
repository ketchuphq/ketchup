import * as API from 'lib/api';
import * as React from 'react';
import GenericStore, {Data} from 'lib/store';

interface EditorProps {
  index: number;
  routesStore: GenericStore<Data<API.Route[]>>;
}

interface EditorState {
  path: string;
}

class RouteEditor extends React.Component<EditorProps, EditorState> {
  readonly originalPath: string;
  constructor(props: EditorProps) {
    super(props);
    this.state = {
      path: this.getRoute().path || '',
    };
    this.originalPath = this.state.path;
  }

  getRoute = () => this.props.routesStore.obj.current[this.props.index];

  componentDidMount() {
    this.props.routesStore.subscribe(this.originalPath, (data) => {
      this.setState({
        path: data.current[this.props.index].path || '',
      });
    });
  }

  componentWillUnmount() {
    this.props.routesStore.unsubscribe(this.originalPath);
  }

  render() {
    return (
      <div>
        <input
          type="text"
          placeholder="/path/to/page"
          value={this.state.path}
          onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
            this.props.routesStore.update((data) => {
              data.current[this.props.index].path = e.target.value;
            });
          }}
        />
      </div>
    );
  }
}

interface Props {
  routesStore: GenericStore<Data<API.Route[]>>;
}

interface State {
  routes: API.Route[];
}

export default class PageEditRoutesComponent extends React.Component<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = {routes: []};
    this.props.routesStore.subscribe('edit-route', (data) => {
      this.setState({routes: data.current});
    });
  }

  componentWillMount() {
    this.setState({
      routes: this.props.routesStore.obj.current || [],
    });
  }

  componentWillUnmount() {
    this.props.routesStore.unsubscribe('edit-route');
  }

  render() {
    return (
      <div className="edit-route control">
        <div className="label">Permalink</div>
        {this.state.routes.map((r, i) => (
          <RouteEditor key={r.uuid || 'new'} routesStore={this.props.routesStore} index={i} />
        ))}
      </div>
    );
  }
}
