const { Router } = require('express');
const authMiddleware = require('../../shared/middleware/auth.middleware');
const validate = require('../../shared/middleware/validate.middleware');
const validateParams = require('../../shared/middleware/params.middleware');
const { createProjectSchema, updateProjectSchema, addMemberSchema } = require('./projects.validation');
const { getByUser, getById, create, update, remove, addMember, removeMember } = require('./projects.controller');

const router = Router();

router.use(authMiddleware);

router.get('/', getByUser);
router.get('/:projectId', validateParams('projectId'), getById);
router.post('/', validate(createProjectSchema), create);
router.put('/:projectId', validateParams('projectId'), validate(updateProjectSchema), update);
router.delete('/:projectId', validateParams('projectId'), remove);
router.post('/:projectId/members', validateParams('projectId'), validate(addMemberSchema), addMember);
router.delete('/:projectId/members/:userId', validateParams('projectId', 'userId'), removeMember);

module.exports = router;
