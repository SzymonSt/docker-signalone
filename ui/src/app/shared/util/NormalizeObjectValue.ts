export function NormalizeObjectValue(valueToNormalize: Object): Object {
  Object.keys(valueToNormalize).forEach(key => {
    // @ts-ignore
    if (!valueToNormalize[key]) {
      // @ts-ignore
      delete valueToNormalize[key];
    }
  });
  return valueToNormalize;
}