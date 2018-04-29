import {Editor, ShortTextEditorComponent} from 'components/content';
import EditorComponent from 'components/editors/editor';
import CodeMirrorComponent from 'components/editors/markdown';
import TextEditorComponent from 'components/editors/text';
import {shallow} from 'enzyme';
import * as API from 'lib/api';
import * as React from 'react';

describe('<Editor>', function() {
  function testEditor(content: API.Content, editorComponent: any) {
    const editor = shallow(<Editor content={content} hideLabel={false} />, {
      disableLifecycleMethods: true,
    });
    expect(editor).toMatchSnapshot();
    expect(editor.containsMatchingElement(editorComponent)).toBe(true);
    expect(editor.find(editorComponent).dive()).toMatchSnapshot();
  }

  it('text:markdown -> LongMarkdownEditor', function() {
    testEditor(
      {
        key: 'akey',
        value: '*hello world*',
        text: {
          title: 'thetitle',
          type: 'markdown',
        },
      },
      CodeMirrorComponent
    );
  });

  it('text:text -> LongTextEditor', function() {
    testEditor(
      {
        key: 'akey',
        value: '*hello world*',
        text: {
          title: 'thetitle',
          type: 'text',
        },
      },
      TextEditorComponent
    );
  });

  it('text:html -> LongHTMLEditor', function() {
    testEditor(
      {
        key: 'akey',
        value: '*hello world*',
        text: {
          title: 'thetitle',
          type: 'html',
        },
      },
      EditorComponent
    );
  });

  it('short:markdown -> ShortMarkdownEditor', function() {
    testEditor(
      {
        key: 'akey',
        value: '*hello world*',
        short: {
          title: 'thetitle',
          type: 'markdown',
        },
      },
      ShortTextEditorComponent
    );
  });
});
