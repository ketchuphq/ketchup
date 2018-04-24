import * as React from 'react';
import * as API from 'lib/api';
import {Editor} from 'components/content';

const mainKey = 'content';

interface Props {
  contents: API.Content[];
}

export default class PageEditorsComponent extends React.PureComponent<Props> {
  render() {
    let contents = this.props.contents;
    return (
      <div>
        {contents
          .filter((content) => content.key != mainKey)
          .map((content) => <Editor key={content.key} content={content} hideLabel={false} />)}
        {contents
          .filter((content) => content.key == mainKey)
          .map((content) => <Editor key={content.key} content={content} hideLabel={false} />)}
      </div>
    );
  }
}
