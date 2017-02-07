import msx from 'lib/msx';
import Navigation from 'components/navigation';
import * as Toaster from 'components/toaster';

let retain = (_: any, __: any, context: Mithril.Context) => context.retain = true;

export default (content: string | Mithril.VirtualElement) =>
  <div class='container' key='container' config={retain}>
    {Toaster.render()}
    <Navigation />
    <div class='container__body'>{content}</div>
  </div>;
