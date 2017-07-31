import * as m from 'mithril';
import msx from 'lib/msx';
import * as API from 'lib/api';
import { renderEditor } from 'components/content';
import { BaseComponent } from '../../components/auth';

const mainKey = 'content';

interface PageEditorsAttrs {
  contents: API.Content[];
}

export default class PageEditorsComponent extends BaseComponent<PageEditorsAttrs> {
  view(v: m.Vnode<PageEditorsAttrs, {}>) {
    let contents = v.attrs.contents;
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
