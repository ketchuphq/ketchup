import * as React from 'react';
import * as API from 'lib/api';

interface Props {
  key?: string;
  readonly content: API.Content;
  readonly short?: boolean;
}

export default class TextEditorComponent extends React.Component<Props> {
  textInput: React.RefObject<HTMLTextAreaElement>;

  constructor(props: Props) {
    super(props);
    this.textInput = React.createRef();
  }

  componentDidMount() {
    this.textInput.current.value = this.props.content.value;
  }

  get klass(): string {
    let k = '.text';
    if (this.props.short) {
      return k + '.text-short';
    }
    return k;
  }

  render() {
    return (
      <div className={this.klass}>
        <textarea
          onChange={(el: React.ChangeEvent<HTMLTextAreaElement>) => {
            this.props.content.value = (el.target as any).value;
          }}
          ref={this.textInput}
        />
      </div>
    );
  }
}
