export const get = (url: string) =>
  fetch(url, {
    method: 'GET',
    credentials: 'same-origin',
  });

export const post = (url: string, data?: any) =>
  fetch(url, {
    method: 'POST',
    credentials: 'same-origin',
    body: data ? JSON.stringify(data) : null,
  });
