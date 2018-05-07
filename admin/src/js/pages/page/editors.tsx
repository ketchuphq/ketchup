import {Editor} from 'components/content';
import * as API from 'lib/api';
import {FileUpload} from 'pages/page/files';
import * as React from 'react';

const mainKey = 'content';

interface Props {
  contents: API.Content[];
}

interface State {
  currentRef: string;
}

export default class PageEditorsComponent extends React.PureComponent<Props, State> {
  editorRef: {[key: string]: React.RefObject<Editor>};

  constructor(props: Props) {
    super(props);
    this.editorRef = {};
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
      <div>
        {contents.filter((content) => content.key != mainKey).map((content) => (
          <div
            key={content.key}
            onFocus={() => {
              this.setState({currentRef: content.key});
            }}
          >
            <Editor ref={this.getRef(content.key)} content={content} hideLabel={false} />
          </div>
        ))}
        {contents.filter((content) => content.key == mainKey).map((content) => (
          <div
            key={content.key}
            onFocus={() => {
              this.setState({currentRef: content.key});
            }}
          >
            <Editor ref={this.getRef(content.key)} content={content} hideLabel={false} />
          </div>
        ))}
        <FileUpload
          onDrop={(file) => {
            let ref = this.getRef(this.state.currentRef);
            ref.current.insertFile(file.name, file.url);
          }}
        />
      </div>
    );
  }
}
