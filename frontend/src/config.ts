// config.ts
const protocol = window.location.protocol;
const hostname = window.location.hostname;
const port = window.location.port || (protocol === 'https:' ? '443' : '80');

export const endpoint_base = `${protocol}//${hostname}:${port}`;
export const protocol_base = `${protocol}`;
// Usage
console.log(endpoint_base);  // Will log something like 'https://example.com:443'
