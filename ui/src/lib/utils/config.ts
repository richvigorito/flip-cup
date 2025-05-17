const rawUrl = import.meta.env.VITE_WS_URL || 'localhost:8080';
const cleanUrl = stripProtocol(rawUrl);

const isSecure = window.location.protocol === 'https:';

export const baseHttpUrl = `${isSecure ? 'https' : 'http'}://${cleanUrl}`;
export const baseWsUrl = `${isSecure ? 'wss' : 'ws'}://${cleanUrl}/ws`;

console.log('baseHttpUrl:', baseHttpUrl);
console.log('baseWsUrl:', baseWsUrl);

function stripProtocol(url: string) {
  return url.replace(/^wss?:\/\//, '').replace(/^https?:\/\//, '');
}
