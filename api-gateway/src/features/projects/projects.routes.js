const { Router } = require('express');
const authMiddleware = require('../../shared/middleware/auth.middleware');
const validate = require('../../shared/middleware/validate.middleware');
const { createProjectSchema, updateProjectSchema, addMemberSchema } = require('./projects.validation');
const { getByUser, getById, create, update, remove, addMember, removeMember } = require('./projects.controller');

const router = Router();

router.use(authMiddleware);

router.get('/', getByUser);
router.get('/:projectId', getById);
router.post('/', validate(createProjectSchema), create);
router.put('/:projectId', validate(updateProjectSchema), update);
router.delete('/:projectId', remove);
router.post('/:projectId/members', validate(addMemberSchema), addMember);
router.delete('/:projectId/members/:userId', removeMember);

module.exports = router;
