// Restricts access based on user roles
// Note: Role validation is handled by the Go Task Service
// This middleware adds an extra layer of protection at the gateway level
const requireRole = (...allowedRoles) => {
  return (req, res, next) => {
    const userRole = req.headers['x-user-role'];

    if (!userRole) {
      return res.status(403).json({ error: 'access denied' });
    }

    if (!allowedRoles.includes(userRole)) {
      return res.status(403).json({ error: 'insufficient permissions' });
    }

    next();
  };
};

module.exports = requireRole;
