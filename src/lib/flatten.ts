// 👁️ https://www.30secondsofcode.org/js/s/flatten-unflatten-object/

export function flattenObject(
  obj: any,
  delimiter = '.',
  prefix = ''
): any {
  return Object.keys(obj).reduce(
    (acc, k) => {
      const pre = prefix.length ? `${prefix}${delimiter}` : '';
      if (
        typeof obj[k] === 'object' &&
        obj[k] !== null &&
        Object.keys(obj[k]).length > 0
      )
        Object.assign(acc, flattenObject(obj[k], delimiter, pre + k));
      else acc[pre + k] = obj[k];
      return acc;
    },
    <any>{}
  );
}
