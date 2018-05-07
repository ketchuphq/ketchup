export const get = (url: string) =>
  fetch(url, {
    method: 'GET',
    credentials: 'same-origin',
  });

export const post = (url: string, data?: any, json = true) =>
  fetch(url, {
    method: 'POST',
    credentials: 'same-origin',
    body: data && json ? JSON.stringify(data) : data,
  });

export const del = (url: string, data?: any) =>
  fetch(url, {
    method: 'DELETE',
    credentials: 'same-origin',
    body: data ? JSON.stringify(data) : null,
  });
