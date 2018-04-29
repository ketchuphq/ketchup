import * as API from 'lib/api';
import {get, post} from 'lib/requests';

export default class Data extends API.Data {
  constructor(config?: API.Data) {
    super();
    API.Data.copy(config, this);
  }

  static saveList(data: API.Data[]): Promise<API.Data> {
    return post(`/api/v1/data`, {data}).then((res) => res.json());
  }

  static list(): Promise<Data[]> {
    return get('/api/v1/data')
      .then((res) => res.json())
      .then((res: API.ListDataResponse) => {
        return res.data.map((el) => new Data(el));
      });
  }
}
