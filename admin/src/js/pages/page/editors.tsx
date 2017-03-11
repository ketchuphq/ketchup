import msx from 'lib/msx';
import * as API from 'lib/api';
import { defaultContent, default as Page } from 'lib/page';
import { renderEditor } from 'components/content';

const mainKey = 'content';

interface PageEditorsAttrs {
  page: Page;
  template: API.ThemeTemplate;
}

export default class PageEditorsComponent {
  private readonly _attrs: PageEditorsAttrs;
  constructor() { }
  static view({ attrs: { page, template } }: Mithril.Vnode<PageEditorsAttrs, {}>) {
    if (!page || !template) {
      return;
    }

    let contentMap: { [key: string]: API.Content } = {};
    let placeholderContents: API.Content[] = [];
    let pageContents: API.Content[] = [];

    // 1. get template placeholders
    if (template && template.placeholders) {
      template.placeholders.forEach((p) => {
        let content: API.Content = API.Content.copy(p, {});
        contentMap[p.key] = content;
        placeholderContents.push(content);
      });
    }

    // 2. get page contents
    (page.contents || []).forEach((c) => {
      // update the oneof type of the content from the placeholder
      // todo: exhaustively map fields; convert markdown <> css more elegantly
      // this block updates the existing key
      if (contentMap[c.key]) {
        contentMap[c.key].uuid = c.uuid;
        contentMap[c.key].timestamps = c.timestamps;
        contentMap[c.key].value = c.value;
      } else {
        contentMap[c.key] = c;
        pageContents.push(c);
      }
    });

    page.contents = pageContents.concat(placeholderContents);

    // 3. main content
    let mainContent;
    if (template.hideContent !== false) {
      if (!contentMap[mainKey]) {
        contentMap[mainKey] = API.Content.copy(defaultContent, {});
        page.contents.push(contentMap[mainKey]);
      }
      mainContent = renderEditor(page, contentMap[mainKey], true);
    }

    return <div>
      {placeholderContents
        .filter((c) => c.key != mainKey)
        .map((p) => renderEditor(page, p, false))}
      {mainContent}
    </div>;
  }
}
