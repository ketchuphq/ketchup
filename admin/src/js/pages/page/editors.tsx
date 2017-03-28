import msx from 'lib/msx';
import * as API from 'lib/api';
import { renderEditor } from 'components/content';

const mainKey = 'content';

interface PageEditorsAttrs {
  contents: API.Content[];
}

export default class PageEditorsComponent {
  private readonly _attrs: PageEditorsAttrs;

  constructor() {
  }

  static view({ attrs: { contents } }: Mithril.Vnode<PageEditorsAttrs, {}>) {
    return <div>
      {contents
        .filter((content) => content.key != mainKey)
        .map((content) => renderEditor(content, false))}
      {contents
        .filter((content) => content.key == mainKey)
        .map((content) => renderEditor(content, false))}
    </div>;
  }
}
