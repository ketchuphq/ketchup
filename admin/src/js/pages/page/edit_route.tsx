import * as API from 'lib/api';
import Route from 'lib/route';
import * as React from 'react';

interface Props {
  page: API.Page;
  routes: API.Route[];
}

interface State {
  dirty: boolean;
}

export default class EditRoutesComponent extends React.Component<Props, State> {
  constructor(props: Props) {
    super(props);
    let {routes} = props;
    let dirty = true;
    if (routes.length == 0) {
      routes.push({});
      dirty = false;
    } else if (!routes[0].path) {
      dirty = false;
    }
    this.state = {dirty};
  }

  infer() {
    let {page, routes} = this.props;
    if (this.state.dirty) {
      return;
    }
    if (routes.length < 1) {
      return;
    }
    if (!!routes[0].path) {
      return;
    }
    if (!page.title || page.title.trim() == '') {
      return;
    }
    routes[0].path = Route.format(page.title);
  }

  routeEditor = (route: API.Route, i: number) => {
    this.infer();
    let key = route.uuid || route.path || 'new';
    return (
      <div key={key}>
        <input
          type="text"
          placeholder="/path/to/page"
          value={Route.format(route.path)}
          onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
            this.setState({dirty: true});
            route.path = Route.format(e.target.value);
          }}
        />
        {i > 0 ? <a onClick={() => this.props.routes.splice(i, 1)}>&times;</a> : null}
      </div>
    );
  };

  render() {
    return (
      <div className="edit-route control">
        <div className="label">Permalink</div>
        {this.props.routes.map(this.routeEditor)}
      </div>
    );
  }
}
