import * as React from 'react';
import * as API from 'lib/api';

interface Props {
  readonly content: API.Content;
  readonly short?: boolean;
}

export default class QuillComponent extends React.Component<Props> {
  quill: any; // todo: figure out how to set this correctly
  quillInput: React.RefObject<HTMLDivElement>;

  constructor(props: Props) {
    super(props);
    this.quillInput = React.createRef();
  }

  get klass(): string {
    let k = 'quill';
    if (this.props.short) {
      return k + 'quill-short';
    }
    return k;
  }

  async initializeQuill(element: HTMLDivElement) {
    if (this.props.content.value) {
      element.innerHTML = this.props.content.value;
    }
    // todo: figure out how to type this correctly
    const Quill: any = await import(/* webpackChunkName: "quill" */ 'quill');
    this.quill = new Quill(`#quill`, {
      placeholder: 'start typing...',
      theme: 'snow',
      modules: {
        toolbar: '#quill-toolbar',
      },
    });
    let updateContent = () => {
      let editor = element.getElementsByClassName('ql-editor')[0];
      this.props.content.value = editor.innerHTML;
    };
    this.quill.on('text-change', updateContent.bind(this));
  }

  componentDidMount() {
    this.initializeQuill(this.quillInput.current);
  }

  render() {
    return (
      <div className={this.klass}>
        <div id="quill-toolbar">
          <div className="ql-formats">
            <select className="ql-header" />
          </div>
          <div className="ql-formats">
            <button className="ql-bold" />
            <button className="ql-italic" />
            <button className="ql-underline" />
            <button className="ql-link" />
          </div>
          <div className="ql-formats">
            <button className="ql-list" value="ordered" />
            <button className="ql-list" value="bullet" />
          </div>
          <div className="ql-formats">
            <button className="ql-clean" />
          </div>
        </div>
        <div id="quill" key="quill" ref={this.quillInput} />
      </div>
    );
  }
}
