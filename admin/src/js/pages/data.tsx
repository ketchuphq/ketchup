import msx from 'lib/msx';
import Data from 'lib/data';
import { MustAuthController } from 'components/auth';
import { loading } from 'components/loading';
import { renderEditor } from 'components/content';
import Button from 'components/button';


let defaultData = [
  'title',
]

export default class DataPage extends MustAuthController {
  data: Data[];
  loading: boolean;

  constructor() {
    super();
    this.data = [];
    this.loading = true;
    Data.list()
      .then((data) => {
        let headMap: { [key: string]: Data} = {}
        let head: Data[] = []
        let tail: Data[] = []
        data.map((d) => {
          if (defaultData.indexOf(d.key) > -1) {
            headMap[d.key] = d
          } else {
            tail.push(d)
          }
        })
        defaultData.map((k) => {
          if (k in headMap) {
            head.push(headMap[k])
          } else {
            head.push({
              key: k,
              short: { type: 'text' }
            })
          }
        })

        this.data = head.concat(tail)
        this.loading = false
      })
  }

  static oninit(v: Mithril.Vnode<{}, DataPage>) {
    v.state = new DataPage();
  };

  static view(v: Mithril.Vnode<{}, DataPage>) {
    let ctrl = v.state;
    return <div class='data'>
      <header>
        <h1>Data</h1>
      </header>
      <div class='table'>
        {loading(ctrl.loading)}
        {ctrl.data.map((data) =>
          <div class='tr tr--center'>
            <label>{data.key}</label>
            <div>{renderEditor(data, true)}</div>
          </div>
        )}
        <div class='tr tr--right'>
          <Button
            class='button--green'
            handler={() => Data.saveList(ctrl.data) }>
            Save
          </Button>
        </div>
      </div>
    </div>;
  }
}

let _: Mithril.Component<{}, DataPage> = DataPage;
