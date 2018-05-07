import * as React from 'react';
import * as API from 'lib/api';

interface Props {
  key?: string;
  readonly content: API.Content;
  readonly short?: boolean;
}

export default class TextEditorComponent extends React.Component<Props> {
  render() {
    let k = this.props.short ? '.text .text-short' : '.text';
    return (
      <div className={k}>
        <textarea
          onChange={(el: React.ChangeEvent<HTMLTextAreaElement>) => {
            this.props.content.value = (el.target as any).value;
          }}
          defaultValue={this.props.content.value}
        />
      </div>
    );
  }
}
