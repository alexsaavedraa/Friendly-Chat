// config.ts
const protocol = window.location.protocol;
const hostname = window.location.hostname;
const port = window.location.port || (protocol === 's:' ? '443' : '80');

export const endpoint_base = `${hostname}:${port}`;
export const protocol_base = `${protocol}`;
// Usage
console.log(endpoint_base);  // Will log something like 's://example.com:443'
console.log(protocol_base); 