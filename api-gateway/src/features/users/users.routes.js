const { Router } = require('express');
const authMiddleware = require('../../shared/middleware/auth.middleware');
const { getAll, updateRole } = require('./users.controller');

const router = Router();

router.use(authMiddleware);

// Role enforcement is handled by the Go Task Service (admin only)
router.get('/', getAll);
router.put('/:userId/role', updateRole);

module.exports = router;
