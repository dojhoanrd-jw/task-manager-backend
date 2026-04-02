const AppError = require('../errors/app-error');

// Global error handling middleware
const errorMiddleware = (err, req, res, next) => {
  // Handle custom AppError
  if (err instanceof AppError) {
    return res.status(err.statusCode).json({ error: err.message });
  }

  // Handle Axios errors from the Go Task Service
  if (err.response) {
    const status = err.response.status;
    const message = err.response.data?.error || 'service error';
    return res.status(status).json({ error: message });
  }

  // Handle connection errors to the Go Task Service
  if (err.code === 'ECONNREFUSED') {
    return res.status(503).json({ error: 'task service is unavailable' });
  }

  // Handle unexpected errors
  console.error('Unexpected error:', err.message);
  res.status(500).json({ error: 'internal server error' });
};

module.exports = errorMiddleware;
