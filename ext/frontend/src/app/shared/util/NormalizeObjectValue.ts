export function NormalizeObjectValue(valueToNormalize: Object, dateKeys?: string[]): Object {
  if (dateKeys) {
    dateKeys.forEach((dateKey => {
      // @ts-ignore
      if (valueToNormalize[dateKey]) {
        // @ts-ignore
        valueToNormalize[dateKey] = new Date(valueToNormalize[dateKey]).toISOString();
      }
    }))
  }
  Object.keys(valueToNormalize).forEach(key => {
    // @ts-ignore
    if (!valueToNormalize[key] && !(valueToNormalize[key] === 0 || valueToNormalize[key] === false)) {
      // @ts-ignore
      delete valueToNormalize[key];
    }
  });
  return valueToNormalize;
}