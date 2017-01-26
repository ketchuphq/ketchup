import * as API from 'lib/api';
import Theme from 'lib/theme';
import Layout from 'components/layout';
import { MustAuthController } from 'components/auth';

export default class TemplatePage extends MustAuthController {
  template: Mithril.Property<API.ThemeTemplate>;

  constructor() {
    super();
    this.template = m.prop<API.ThemeTemplate>();
    Theme.getFullTemplate(m.route.param('name'), m.route.param('template')).then((t) => {
      this.template(t);
    });
    // todo: catch
  }

  colorize(el: HTMLElement) {
    require.ensure([
      'highlight.js',
      'highlight.js/lib/languages/xml',
      'highlight.js/styles/rainbow.css'
    ], (require) => {
      let hljs: any = require('highlight.js');
      require('highlight.js/lib/languages/xml');
      require('highlight.js/styles/rainbow.css');
      hljs.highlightBlock(el);
    }, 'hljs');
  }

  static controller = TemplatePage;
  static view(ctrl: TemplatePage) {
    let name = m.route.param('name');
    return Layout(m('.template', [
      m('header',
        m('h1', [
          m('a[href=/admin/themes]', { config: m.route }, 'Themes'),
          m.trust(' &rsaquo; '),
          m(`a[href=/admin/themes/${name}]`, { config: m.route }, name),
          m.trust(' &rsaquo; '),
          m.route.param('template')
        ]),
      ),
      m('pre', {
        config: (el: HTMLElement, isInitialized: boolean) => {
          if (!isInitialized) {
            ctrl.colorize(el);
          }
        }
      },
        ctrl.template().data)
    ]));
  }
}