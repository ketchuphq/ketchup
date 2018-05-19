import Button from 'components/button';
import {Editor} from 'components/content';
import Layout from 'components/layout';
import {Table} from 'components/table';
import Data from 'lib/data';
import * as React from 'react';

let defaultData = ['title'];

interface State {
  data: Data[];
  loading: boolean;
}

export default class DataPage extends React.Component<{}, State> {
  constructor(props: any) {
    super(props);
    this.state = {
      data: [],
      loading: true,
    };
  }

  async componentDidMount() {
    let data = await Data.list();
    let headMap: {[key: string]: Data} = {};
    let head: Data[] = [];
    let tail: Data[] = [];
    data.map((d) => {
      if (defaultData.indexOf(d.key) > -1) {
        headMap[d.key] = d;
      } else {
        tail.push(d);
      }
    });
    defaultData.map((k) => {
      if (k in headMap) {
        head.push(headMap[k]);
      } else {
        head.push({
          key: k,
          short: {type: 'text'},
        });
      }
    });

    this.setState({
      data: head.concat(tail),
      loading: false,
    });
  }

  render() {
    return (
      <Layout className="data">
        <header>
          <h1>Data</h1>
        </header>
        <Table loading={this.state.loading}>
          {this.state.data.map((data) => (
            <div key={data.key} className="tr tr--center">
              <label>{data.key}</label>
              <Editor content={data} hideLabel />
            </div>
          ))}
          <div className="tr tr--right">
            <Button className="button--green" handler={() => Data.saveList(this.state.data)}>
              Save
            </Button>
          </div>
        </Table>
      </Layout>
    );
  }
}
