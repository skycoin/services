export const getParameterByName = (name, url) => {
  const u = url || window.location.href;
  const n = name.replace(/[\[\]]/g, '\\$&');
  const regex = new RegExp(`[?&]${n}(=([^&#]*)|&|#|$)`);
  const results = regex.exec(u);
  if (!results) return null;
  if (!results[2]) return null;
  return decodeURIComponent(results[2].replace(/\+/g, ' '));
};
