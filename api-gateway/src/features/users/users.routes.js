const { Router } = require('express');
const authMiddleware = require('../../shared/middleware/auth.middleware');
const validate = require('../../shared/middleware/validate.middleware');
const { updateRoleSchema } = require('./users.validation');
const { getAll, updateRole } = require('./users.controller');

const router = Router();

router.use(authMiddleware);

// Role enforcement is handled by the Go Task Service (admin only)
router.get('/', getAll);
router.put('/:userId/role', validate(updateRoleSchema), updateRole);

module.exports = router;
