import { Editor } from 'components/content';
import * as API from 'lib/api';
import { FileUpload } from 'pages/page/files';
import { ContentPreview } from 'pages/page/preview';
import * as React from 'react';

const mainKey = 'content';

interface Props {
  showPreview: boolean;
  contents: API.Content[];
}

interface State {
  currentContent?: API.Content;
  previewContentRequest?: string;
}

export default class PageEditorsComponent extends React.Component<Props, State> {
  editorRef: {[key: string]: React.RefObject<Editor>};
  update: any;

  constructor(props: Props) {
    super(props);
    this.editorRef = {};
    this.state = {};
  }

  refreshContent = (content: API.Content) => {
    if (content && this.props.showPreview) {
      this.setState({
        previewContentRequest: JSON.stringify({
          content: content,
        } as API.PreviewContentRequest),
      });
    }
  };

  componentDidMount() {
    if (this.props.contents.length > 0) {
      this.setState({
        currentContent: this.props.contents[0],
        previewContentRequest: this.props.showPreview
          ? JSON.stringify({
              content: this.props.contents[0],
            } as API.PreviewContentRequest)
          : '',
      });
    }
    this.update = setInterval(() => this.refreshContent(this.state.currentContent), 1000);
  }

  componentWillUnmount() {
    clearInterval(this.update);
  }

  getRef = (key: string) => {
    if (!this.editorRef[key]) {
      this.editorRef[key] = React.createRef();
    }
    return this.editorRef[key];
  };

  render() {
    let contents = this.props.contents;
    return (
      <div className="page-editors">
        <div className="preview-left">
          {contents.filter((content) => content.key != mainKey).map((content) => (
            <div
              key={content.key}
              onFocus={() => {
                this.setState({currentContent: content});
                this.refreshContent(content);
              }}
            >
              <Editor ref={this.getRef(content.key)} content={content} hideLabel={false} />
            </div>
          ))}
          {contents.filter((content) => content.key == mainKey).map((content) => (
            <div
              key={content.key}
              onFocus={() => {
                this.setState({currentContent: content});
                this.refreshContent(content);
              }}
            >
              <Editor ref={this.getRef(content.key)} content={content} hideLabel={false} />
            </div>
          ))}
          <FileUpload
            onDrop={(file) => {
              let ref = this.getRef(this.state.currentContent.key);
              ref.current.insertFile(file.name, file.url);
            }}
          />
        </div>
        {this.props.showPreview ? (
          <div className="preview-right">
            <div className="preview-right__inner">
              {this.state.currentContent ? (
                <p className="preview-note">Note: Preview does not reflect final style.</p>
              ) : null}
              {this.state.currentContent ? (
                <ContentPreview
                  contentKey={this.state.currentContent.key}
                  content={this.state.previewContentRequest}
                />
              ) : (
                <p key="1" className="preview-note">
                  No preview
                </p>
              )}
            </div>
          </div>
        ) : null}
      </div>
    );
  }
}
