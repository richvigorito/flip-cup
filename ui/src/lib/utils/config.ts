// src/lib/utils/config.ts

const rawUrl = import.meta.env.VITE_WS_URL || 'localhost:8080';
const cleanUrl = stripProtocol(rawUrl);

console.log('rawurl:', rawUrl);
console.log('cleanurl:', cleanUrl);

export const baseHttpUrl = `http://${cleanUrl}`;
export const baseWsUrl = `ws://${cleanUrl}/ws`;


console.log('baseurl:', baseHttpUrl);
console.log('wsurl:', baseWsUrl);

function stripProtocol(url: string) {
  return url.replace(/^wss?:\/\//, '').replace(/^https?:\/\//, '');
}

