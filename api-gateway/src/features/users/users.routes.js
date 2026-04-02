const { Router } = require('express');
const authMiddleware = require('../../shared/middleware/auth.middleware');
const validate = require('../../shared/middleware/validate.middleware');
const validateParams = require('../../shared/middleware/params.middleware');
const { updateRoleSchema } = require('./users.validation');
const { getAll, updateRole } = require('./users.controller');

const router = Router();

router.use(authMiddleware);

router.get('/', getAll);
router.put('/:userId/role', validateParams('userId'), validate(updateRoleSchema), updateRole);

module.exports = router;
