const { Router } = require('express');
const authMiddleware = require('../../shared/middleware/auth.middleware');
const { getByUser, getById, create, update, remove, addMember, removeMember } = require('./projects.controller');

const router = Router();

router.use(authMiddleware);

router.get('/', getByUser);
router.get('/:projectId', getById);
router.post('/', create);
router.put('/:projectId', update);
router.delete('/:projectId', remove);
router.post('/:projectId/members', addMember);
router.delete('/:projectId/members/:userId', removeMember);

module.exports = router;
