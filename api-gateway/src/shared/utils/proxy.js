const httpClient = require('./http-client');

// proxy wraps httpClient calls to reduce boilerplate in controllers
const proxy = {
  get: (path, req, params = {}) =>
    httpClient.get(path, {
      headers: { Authorization: req.token, 'X-Request-ID': req.requestId },
      params,
    }),

  post: (path, req) =>
    httpClient.post(path, req.body, {
      headers: { Authorization: req.token, 'X-Request-ID': req.requestId },
    }),

  put: (path, req) =>
    httpClient.put(path, req.body, {
      headers: { Authorization: req.token, 'X-Request-ID': req.requestId },
    }),

  delete: (path, req) =>
    httpClient.delete(path, {
      headers: { Authorization: req.token, 'X-Request-ID': req.requestId },
    }),
};

module.exports = proxy;
