const rateLimit = require('express-rate-limit');

// General API rate limiter: 100 requests per minute per IP
const apiLimiter = rateLimit({
  windowMs: 60 * 1000,
  max: 100,
  standardHeaders: true,
  legacyHeaders: false,
  message: { error: 'too many requests, please try again later' },
});

// Auth rate limiter: 10 requests per minute per IP (stricter for login/register)
const authLimiter = rateLimit({
  windowMs: 60 * 1000,
  max: 10,
  standardHeaders: true,
  legacyHeaders: false,
  message: { error: 'too many authentication attempts, please try again later' },
});

module.exports = { apiLimiter, authLimiter };
