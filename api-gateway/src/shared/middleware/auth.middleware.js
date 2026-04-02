const config = require('../config');

// Verifies JWT token by forwarding the request to the Go Task Service
// The Go service handles the actual JWT validation
const authMiddleware = (req, res, next) => {
  const authHeader = req.headers.authorization;

  if (!authHeader) {
    return res.status(401).json({ error: 'authorization header is required' });
  }

  if (!authHeader.startsWith('Bearer ')) {
    return res.status(401).json({ error: 'invalid authorization format' });
  }

  // Forward the authorization header to the Go service
  // The Go service will validate the token and return user data in headers
  req.token = authHeader;
  next();
};

module.exports = authMiddleware;
