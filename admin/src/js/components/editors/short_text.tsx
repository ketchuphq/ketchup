import * as API from 'lib/api';
import * as React from 'react';

export default class ShortTextEditorComponent extends React.PureComponent<{content: API.Content}> {
  textInput: React.RefObject<HTMLInputElement>;
  render() {
    return (
      <input
        type="text"
        defaultValue={this.props.content.value}
        onChange={(e: React.ChangeEvent<HTMLInputElement>) => {
          this.props.content.value = e.target.value;
        }}
      />
    );
  }
}
