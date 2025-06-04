export const objDeepCopy = (source: Object) => {
  const result = JSON.parse(JSON.stringify(source || {}))
  return result
}

export function copyToClipboard(val: any) {
  const t = document.createElement("textarea");
  document.body.appendChild(t);
  t.value = val;
  t.select();
  document.execCommand('copy');
  document.body.removeChild(t);
}