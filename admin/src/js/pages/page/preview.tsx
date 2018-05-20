import * as React from 'react';

interface Props {
  contentKey: string;
  content: string;
}

const formTarget = 'content-preview';

export class ContentPreview extends React.Component<Props> {
  ref: React.RefObject<HTMLDivElement>;
  update: any;

  constructor(props: Props) {
    super(props);
    this.ref = React.createRef();
  }

  makeFrameRequest = (content: string) => {
    const form = document.createElement('form');
    form.action = '/admin/preview';
    form.method = 'post';
    form.target = formTarget;

    const input = document.createElement('input');
    input.type = 'hidden';
    input.name = 'content';
    input.value = content;

    form.appendChild(input);
    this.ref.current.appendChild(form);

    form.submit();
    form.remove();
  };

  updateFrame() {
    clearTimeout(this.update);
    this.update = setTimeout(() => this.makeFrameRequest(this.props.content), 1000);
  }

  componentWillUnmount() {
    clearTimeout(this.update);
  }

  componentDidMount() {
    if (this.props.content) {
      this.makeFrameRequest(this.props.content);
    } else {
      this.updateFrame();
    }
  }

  componentWillReceiveProps(props: Props) {
    if (this.props.contentKey != props.contentKey) {
      this.makeFrameRequest(props.content);
    } else if (this.props.content != props.content) {
      this.updateFrame();
    }
  }

  render() {
    return (
      <div className="content-preview" ref={this.ref}>
        <iframe name={formTarget} sandbox="allow-scripts" />
      </div>
    );
  }
}
