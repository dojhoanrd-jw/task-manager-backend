const axios = require('axios');
const config = require('../config');

// Axios instance configured to communicate with the Go Task Service
const httpClient = axios.create({
  baseURL: config.taskServiceUrl,
  timeout: 10000,
  headers: {
    'Content-Type': 'application/json',
  },
});

// Interceptor to forward request ID to the Go service
httpClient.interceptors.request.use((reqConfig) => {
  if (reqConfig.requestId) {
    reqConfig.headers['X-Request-ID'] = reqConfig.requestId;
  }
  return reqConfig;
});

module.exports = httpClient;
