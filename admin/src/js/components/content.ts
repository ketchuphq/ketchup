import * as m from 'mithril';
import * as API from 'lib/api';
import Page from 'lib/page';
import EditorComponent from 'components/editors/editor';
import QuillComponent from 'components/editors/html';
import TextEditorComponent from 'components/editors/text';
import CodeMirrorComponent from 'components/editors/markdown';

const control = (...el: any[]) => m('.controls', m('.control.control-full', el));

export function renderEditor(page: Page, c: API.Content, hideLabel: boolean) {
  switch (true) {
    case c.short != null:
      return shortEditor(page, c, hideLabel);
    case c.text != null:
      return textEditor(page, c, hideLabel);
    default:
      console.log('warning: no editor defined for object:', c);
      c.text = { type: 'html' }; // set a default
      return textEditor(page, c, hideLabel);
  }
}

// textEditor is an editor for a large block of text
function textEditor(_: Page, c: API.Content, hideLabel: boolean = false) {
  switch (c.text.type) {
    case 'html':
      return control(
        hideLabel ? '' : m('.label', c.key),
        m(EditorComponent, { content: c })
      );

    case 'markdown':
      return control(
        hideLabel ? '' : m('.label', c.key),
        m(CodeMirrorComponent, { content: c }) // todo: assign id?
      );

    case 'text':
      return control(
        hideLabel ? '' : m('.label', c.key),
        m(TextEditorComponent, { content: c }) // todo: assign id?
      );

    default:
      return null;
  }
}

function shortEditor(page: Page, c: API.Content, hideLabel: boolean = false) {
  switch (c.short.type) {
    case 'html':
      return control(
        hideLabel ? '' : m('.label', c.key),
        m(QuillComponent, { content: c })
      );

    case 'text':
    case 'markdown':
      return control(
        hideLabel ? '' : m('.label', c.key),
        m('input[type=text]', {
          oncreate: (v: Mithril.VnodeDOM<any, any>) => {
            (v.dom as HTMLTextAreaElement).value = c.value || '';
          },
          onchange: (e: EventTarget) => {
            c.value = (e as any).target.value;
            page.updateContent(c);
          }
        })
      );

    default:
      return null;
  }
}