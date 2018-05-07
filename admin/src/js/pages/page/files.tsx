import * as Toaster from 'components/toaster';
import * as API from 'lib/api';
import * as React from 'react';
import * as ReactDOM from 'react-dom';
import * as File from 'lib/file';

interface Props {
  onDrop?: (file: API.File) => void;
}

interface State {
  entered: number;
}

export class FileUpload extends React.PureComponent<Props, State> {
  constructor(props: any) {
    super(props);
    this.state = {
      entered: 0,
    };
  }

  handleDragEvent = (_: DragEvent) => {
    this.setState((state) => ({entered: state.entered + 1}));
  };
  handleDragLeave = (_: DragEvent) => {
    this.setState((state) => ({entered: Math.max(0, state.entered - 1)}));
  };
  handleDragOver = (e: DragEvent) => {
    e.preventDefault();
  };
  handleDrop = (e: React.DragEvent<any> | DragEvent) => {
    e.preventDefault();
    this.setState({entered: 0});
    if (e.dataTransfer) {
      let items = e.dataTransfer.items;
      for (let i = 0; i < items.length; i++) {
        const element = items[i];
        let data = new FormData();
        data.set('file', element.getAsFile());
        File.create(data).then((file) => {
          Toaster.add(`File ${file.name} uploaded.`);
          if (this.props.onDrop) {
            this.props.onDrop(file);
          }
        });
      }
    }
  };

  componentDidMount() {
    document.addEventListener('dragenter', this.handleDragEvent);
    document.addEventListener('dragleave', this.handleDragLeave);
    document.addEventListener('dragover', this.handleDragOver);
    document.addEventListener('drop', this.handleDrop);
  }

  componentWillUnmount() {
    document.removeEventListener('dragenter', this.handleDragEvent);
    document.removeEventListener('dragleave', this.handleDragLeave);
    document.removeEventListener('dragover', this.handleDragOver);
    document.removeEventListener('drop', this.handleDrop);
  }

  render() {
    if (this.state.entered > 0) {
      return ReactDOM.createPortal(
        <div onDrop={this.handleDrop} className="file-upload-overlay animate-fade-in animate-fill">
          <div className="file-upload">drop to upload file</div>
        </div>,
        document.getElementById('react-root')
      );
    }
    return <div />;
  }
}
