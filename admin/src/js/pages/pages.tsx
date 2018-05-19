import Layout from 'components/layout';
import {Table} from 'components/table';
import * as API from 'lib/api';
import * as Page from 'lib/page';
import * as React from 'react';
import {Link} from 'react-router-dom';

interface State {
  pages: API.Page[];
  viewOption: API.ListPageRequest_ListPageFilter;
  loading: boolean;
}

export default class PagesPage extends React.Component<{}, State> {
  constructor(props: any) {
    super(props);
    this.state = {
      pages: [],
      loading: true,
      viewOption: 'all',
    };
  }

  componentDidMount() {
    this.fetch(this.state.viewOption);
  }

  fetch = (val: API.ListPageRequest_ListPageFilter) => {
    return Page.list(val).then((pages) => {
      this.setState({
        pages,
        viewOption: val,
        loading: false,
      });
    });
  };

  render() {
    let tab = (v: API.ListPageRequest_ListPageFilter, desc?: string) => {
      let classes = 'tab-el';
      if (this.state.viewOption == v) {
        classes += ' tab-selected';
      }
      return (
        <span className={classes} onClick={() => this.fetch(v)}>
          {desc || v}
        </span>
      );
    };
    return (
      <Layout className="pages">
        <header>
          <Link className="button button--green button--center" to="/compose">
            Compose
          </Link>
          <h1>Pages</h1>
        </header>
        <h2 className="tabs">
          {tab('all')}
          <span className="tab-divider">|</span>
          {tab('draft')}
          <span className="tab-divider">|</span>
          {tab('published')}
        </h2>
        <Table loading={this.state.loading}>
          {this.state.pages.map((page) => {
            let status = null;
            let klass = '';
            if (Page.isPublished(page)) {
              status = <div className="label small">published</div>;
            } else {
              status = <div className="label label--gray small">draft</div>;
              klass = 'page--draft';
            }

            return (
              <Link key={page.uuid} className="tr tr--center" to={`/pages/${page.uuid}`}>
                <div className={`tr__expand ${klass}`}>{page.title || 'untitled'}</div>
                {status}
                <div className="page--date">{`${Page.formattedUpdatedAt(page)}`}</div>
              </Link>
            );
          })}
        </Table>
      </Layout>
    );
  }
}
