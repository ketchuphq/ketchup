import * as API from 'lib/api';
import Page from 'lib/page';
import QuillComponent from 'components/quill';

const control = (...el: any[]) => m('.controls', m('.control.control-full', el));

export function renderEditor(page: Page, c: API.Content, hideLabel: boolean) {
  switch (true) {
    case c.short != null:
      return shortEditor(page, c, hideLabel);
    case c.text != null:
      return textEditor(page, c, hideLabel);
    default:
      console.log('no editor defined for object:');
      console.log(c);
  }
}

function textEditor(page: Page, c: API.Content, hideLabel: boolean = false) {
  switch (c.text.type) {
    case 'html':
      return control(
        hideLabel ? '' : m('.label', c.key),
        m.component(QuillComponent, c)
      );

    case 'text':
    case 'markdown':
      return control(
        hideLabel ? '' : m('.label', c.key),
        m('textarea', {
          value: c.value ? c.value : '',
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

function shortEditor(page: Page, c: API.Content, hideLabel: boolean = false) {
  switch (c.short.type) {
    case 'html':
      return control(
        hideLabel ? '' : m('.label', c.key),
        m.component(QuillComponent, c, true)
      );

    case 'text':
    case 'markdown':
      return control(
        hideLabel ? '' : m('.label', c.key),
        m('input[type=text]', {
          value: c.value ? c.value : '',
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