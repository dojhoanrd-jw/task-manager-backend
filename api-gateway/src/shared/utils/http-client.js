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

module.exports = httpClient;
