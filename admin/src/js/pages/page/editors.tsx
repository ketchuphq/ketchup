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
  contentMap: { [key: string]: API.Content };

  constructor() {
    this.contentMap = {};
  }

  static oninit(v: Mithril.Vnode<PageEditorsAttrs, PageEditorsComponent>) {
    v.state = new PageEditorsComponent();
  }

  static view({ attrs: { page, template }, state }: Mithril.Vnode<PageEditorsAttrs, PageEditorsComponent>) {
    let contentMap = state.contentMap;
    if (!page || !template) {
      return;
    }

    let placeholderContents: API.Content[] = [];
    let pageContents: API.Content[] = [];

    // 1. get template placeholders
    if (template && template.placeholders) {
      template.placeholders.forEach((p) => {
        if (contentMap[p.key]) {
          Object.keys(p).forEach((k: keyof API.ThemePlaceholder) => {
            contentMap[p.key][k] = p[k];
          });
        } else {
          contentMap[p.key] = API.Content.copy(p, {});
        }
        contentMap[p.key] = API.Content.copy(p, {});
        placeholderContents.push(contentMap[p.key]);
      });
    }

    // 2. get page contents
    (page.contents || []).forEach((content) => {
      // update the oneof type of the content from the placeholder
      // todo: exhaustively map fields; convert markdown <> css more elegantly
      // this block updates the existing key
      if (contentMap[content.key]) {
        contentMap[content.key].uuid = content.uuid;
        contentMap[content.key].timestamps = content.timestamps;
        contentMap[content.key].value = content.value;
      } else {
        contentMap[content.key] = content;
        pageContents.push(content);
      }
    });

    page.contents = pageContents.concat(placeholderContents);

    // 3. main content
    let mainContent;
    if (template.hideContent !== false) {
      if (!contentMap[mainKey]) {
        contentMap[mainKey] = API.Content.copy(defaultContent, {});
      }
      if (page.contents.every((c) => c.key != mainKey)) {
        page.contents.push(contentMap[mainKey]);
      }
      mainContent = renderEditor(contentMap[mainKey], true);
    }

    return <div>
      {placeholderContents
        .filter((content) => content.key != mainKey)
        .map((content) => renderEditor(content, false))}
      {mainContent}
    </div>;
  }
}
