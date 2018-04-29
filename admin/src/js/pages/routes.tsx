import * as React from 'react';
import Route from 'lib/route';
import {Link} from 'react-router-dom';

interface State {
  routes: Route[];
}

export default class RoutesPage extends React.Component<{}, State> {
  constructor(v: any) {
    super(v);
    this.state = {routes: []};
  }

  componentDidMount() {
    Route.list().then((data) => {
      this.setState({routes: data});
    });
  }

  render() {
    return (
      <div className="routes">
        <h1>Routes</h1>
        <div className="table">
          {this.state.routes.map((r) => (
            <div key={r.uuid} className="tr">
              <a href={r.path ? r.path : '#'}>{r.path}</a>
              {!r.pageUuid ? (
                ''
              ) : (
                <Link className="list-link" to={`/admin/pages/${r.pageUuid}`}>
                  edit page
                </Link>
              )}
            </div>
          ))}
        </div>
      </div>
    );
  }
}
