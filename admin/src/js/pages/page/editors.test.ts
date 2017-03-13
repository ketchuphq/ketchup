import { assert } from 'chai';
import * as m from 'mithril';
import * as mq from 'mithril-query';
import * as API from 'lib/api';
import PageEditorsComponent from './editors';
import Page from 'lib/page';

describe('EditorsTest', function() {
  describe('markdown editor', function() {
    it('should propagate deeply set value', function() {
      let attrs = {
        page: new Page(),
        template: {
          placeholders: [
            {
              key: 'content',
              text: {
                type: 'html',
              }
            }
          ]
        } as API.ThemeTemplate
      };
      let out = mq(PageEditorsComponent, attrs);
      let editor = out.find('.editor')[0];
      let enteredValue = 'deeply set value';
      editor.children[0].state.content.value =  enteredValue;
      out.redraw();

      assert.equal(
        out.vnode.state.contentMap['content'].value,
        enteredValue
      );
    });
  })
});
