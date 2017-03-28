import { assert } from 'chai';
import * as mq from 'mithril-query';
import * as API from 'lib/api';
import { renderEditor } from 'components/content';
import CodeMirrorComponent from 'components/editors/markdown';
import TextEditorComponent from 'components/editors/text';
import EditorComponent from 'components/editors/editor';

describe('renderEditor', function() {
  function testEditor(content: API.Content, editorComponent: any) {
    let out = mq(renderEditor(content, false));
    assert.containSubset(out.rootNode, {
      tag: 'div',
      attrs: { className: 'controls' },
      children: [{
        attrs: { className: 'control control-full', },
        children: [{
          attrs: { className: 'label' },
          text: content.key
        }, {
          attrs: { content }
        }]
      }]
    });

    let editor = out.find('.control-full')[0].children[1];
    assert.equal(editor.tag, editorComponent);
    assert.strictEqual(editor.attrs.content, content);

    out = mq(renderEditor(content, true));
    assert.equal(out.find('.label').length, 0);
  }

  it('text:markdown -> LongMarkdownEditor', function() {
    testEditor({
      key: 'akey',
      value: '*hello world*',
      text: {
        title: 'thetitle',
        type: 'markdown',
      }
    }, CodeMirrorComponent);
  });

  it('text:text -> LongTextEditor', function() {
    testEditor({
      key: 'akey',
      value: '*hello world*',
      text: {
        title: 'thetitle',
        type: 'text',
      }
    }, TextEditorComponent);
  });

  it('text:html -> LongHTMLEditor', function() {
    testEditor({
      key: 'akey',
      value: '*hello world*',
      text: {
        title: 'thetitle',
        type: 'html',
      }
    }, EditorComponent);
  });
});
