import Layout from 'components/layout';
import {Loader} from 'components/loading';
import * as API from 'lib/api';
import * as File from 'lib/file';
import * as React from 'react';
import {Link} from 'react-router-dom';
import {FileUpload} from 'pages/page/files';

interface State {
  files: API.File[];
  loading: boolean;
}

export default class FilesPage extends React.Component<{}, State> {
  constructor(props: any) {
    super(props);
    this.state = {
      files: [],
      loading: true,
    };
  }

  componentDidMount() {
    this.fetch();
  }

  fetch = () => {
    return File.list().then((files) => {
      this.setState({files, loading: false});
    });
  };

  render() {
    if (this.state.loading) {
      <Layout className="files">
        <header>
          <h1>Files</h1>
        </header>
        <div className="table">
          <Loader show />
        </div>
      </Layout>;
    }

    return (
      <Layout className="files">
        <header>
          <h1>Files</h1>
        </header>
        <div className="table">
          {!this.state.loading && this.state.files.length == 0 ? (
            <div className="tr">No files uploaded yet. Drag a file in to upload.</div>
          ) : null}
          {this.state.files.map((file) => (
            <Link key={file.uuid} className="tr tr--center" to={`/files/${file.uuid}`}>
              <div className="tr__expand">{file.name}</div>
            </Link>
          ))}
          <FileUpload
            onDrop={(file) => {
              this.setState((state) => {
                let files = state.files.slice();
                files.push(file);
                return {files};
              });
            }}
          />
        </div>
      </Layout>
    );
  }
}
