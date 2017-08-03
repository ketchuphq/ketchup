import * as m from 'mithril';
import * as API from 'lib/api';

export default class Data extends API.Data {
  constructor(config?: API.Data) {
    super();
    API.Data.copy(config, this);
  }

  static saveList(data: API.Data[]): Promise<API.Data> {
    return m.request({
      method: 'POST',
      url: `/api/v1/data`,
      data: { data }
    });
  }

  static list(): Promise<Data[]> {
    return m
      .request({
        method: 'GET',
        url: `/api/v1/data`
      })
      .then((res: API.ListDataResponse) => {
        return res.data.map((el) => new Data(el));
      });
  }
}
