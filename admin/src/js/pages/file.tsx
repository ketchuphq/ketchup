import Layout from 'components/layout';
import {Loader} from 'components/loading';
import * as API from 'lib/api';
import * as File from 'lib/file';
import * as React from 'react';
import {RouteComponentProps} from 'react-router';
import {Link} from 'react-router-dom';
import Button from 'components/button';
import {del} from 'lib/requests';

interface State {
  file?: API.File;
}

export default class FilePage extends React.Component<RouteComponentProps<{id: string}>, State> {
  constructor(props: any) {
    super(props);
    this.state = {};
  }

  componentDidMount() {
    File.get(this.props.match.params.id).then((file) => {
      this.setState({file});
    });
  }

  deleteImage = () => {
    return del(`/api/v1/files/${this.state.file.uuid}`).then(() => {
      location.assign('/admin/files');
    });
  };

  render() {
    if (!this.state.file) {
      return (
        <Layout className="file">
          <header>
            <h1>
              <Link to="/files">Files</Link> &rsaquo; ...
            </h1>
          </header>
          <Loader show />
        </Layout>
      );
    }

    const file = this.state.file;

    const imageExtensions = ['png', 'jpg', 'jpeg', 'gif'];
    let isImage = imageExtensions.some((v) => file.name.substr(-v.length) == v);

    return (
      <Layout className="file">
        <header>
          <h1>
            <Link to="/files">Files</Link> &rsaquo; <span className="unbold">{file.name}</span>
          </h1>
        </header>
        <div>
          <div className="table">
            <div className="tr tr--center">
              <label>{isImage ? 'Image' : 'File'}</label>
              <span className="input-text">
                {isImage ? (
                  <img src={`${file.url}?x=200x200,fit`} />
                ) : (
                  <a href={file.url}>Download {file.name}</a>
                )}
              </span>
            </div>
            <div className="tr tr--center">
              <label>URL</label>
              <span className="input-text">
                <a href={file.url} target="_blank">
                  {file.url}
                </a>
              </span>
            </div>
            <div className="tr tr--right">
              <Button className="button--red" handler={this.deleteImage}>
                Delete
              </Button>
            </div>
          </div>
        </div>
      </Layout>
    );
  }
}
