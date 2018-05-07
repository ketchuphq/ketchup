import * as API from 'lib/api';
import {get as getRequest, post} from 'lib/requests';
import * as Toaster from 'components/toaster';
import {FileResponse} from './api';

export function create(data: any): Promise<API.File> {
  return post('/api/v1/files', data, false)
    .then((res) => {
      if (res.status >= 500) {
        Toaster.error('Error uploading file');
        throw 'error uploading file';
      }
      return res.json();
    })
    .then((res: API.FileResponse) => res.file);
}
export function get(uuid: string): Promise<API.File> {
  return getRequest(`/api/v1/files/${uuid}`)
    .then((res) => res.json())
    .then((data: API.FileResponse) => {
      return data.file;
    });
}

export function list(): Promise<API.File[]> {
  return getRequest(`/api/v1/files`)
    .then((res) => res.json())
    .then((data: API.ListFilesResponse) => {
      return data.files || [];
    });
}
