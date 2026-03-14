const rawUrl = import.meta.env.VITE_WS_URL || (typeof window !== 'undefined' ? window.location.host : 'localhost:8080');
const cleanUrl = stripProtocol(rawUrl);

const isSecure = typeof window !== 'undefined' ? window.location.protocol === 'https:' : false;

export const baseHttpUrl = `${isSecure ? 'https' : 'http'}://${cleanUrl}`;
export const baseWsUrl = `${isSecure ? 'wss' : 'ws'}://${cleanUrl}/ws`;

function stripProtocol(url: string) {
  if (!url) return 'localhost:8080';
  return url.replace(/^wss?:\/\//, '').replace(/^https?:\/\//, '');
}
