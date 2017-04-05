export let serialize = <T>(obj: T): string => {
  let data: { [key: string]: any } = {};
  flatten(obj, data);
  return Object.keys(data)
    .map((k) => `${k}=${data[k]}`)
    .join('&');
};

// flatten object by joining nested keys with '.'
let flatten = (obj: any, output: { [key: string]: any }, prefix = '') => {
  if (prefix != '') {
    prefix += '.';
  }

  Object.keys(obj).forEach((k) => {
    let v = obj[k];
    if (typeof v !== 'object') {
      output[prefix + k] = v;
    } else {
      flatten(v, output, prefix + k);
    }
  });
};
