import * as React from 'react';
import {list} from 'lib/page';
import * as API from 'lib/api';

interface Props {
  onselect: (option: API.Page) => void;
}

interface State {
  pages: API.Page[];
  selected?: string;
}

export class PagePickerComponent extends React.Component<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = {
      pages: [],
    };
  }

  componentDidMount() {
    list().then((pages) => {
      if (pages.length > 0) {
        this.setState({
          pages: pages,
          selected: pages[0].uuid,
        });
        this.props.onselect(pages[0]);
      }
    });
  }

  view() {
    return (
      <select
        value={this.state.selected}
        onChange={(e: React.ChangeEvent<HTMLSelectElement>) => {
          this.setState({selected: e.target.value});
          for (var i = 0; i < this.state.pages.length; i++) {
            let page = this.state.pages[i];
            if (page.uuid == e.target.value) {
              this.props.onselect(page);
              return;
            }
          }
        }}
      >
        {this.state.pages.map((page) => <option key={page.uuid}>{page.uuid}</option>)}
      </select>
    );
  }
}

interface NewRouteState {
  routePageUUID?: string;
  routePath?: string;
}

export class NewRouteComponent extends React.Component<{}, NewRouteState> {
  constructor(v: any) {
    super(v);
    this.state = {
      // route:
    };
  }

  selectPage = (page: API.Page) => {
    this.setState({routePageUUID: page.uuid});
  };

  view() {
    return (
      <div className="new-route">
        <input
          type="text"
          placeholder="route name"
          onChange={(e) => {
            this.setState({routePath: e.target.value});
          }}
        />
        <PagePickerComponent onselect={this.selectPage} />
      </div>
    );
  }
}
