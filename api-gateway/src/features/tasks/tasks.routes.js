const { Router } = require('express');
const authMiddleware = require('../../shared/middleware/auth.middleware');
const { getByProject, getById, create, update, remove } = require('./tasks.controller');

const router = Router();

router.use(authMiddleware);

router.get('/projects/:projectId/tasks', getByProject);
router.get('/projects/:projectId/tasks/:taskId', getById);
router.post('/projects/:projectId/tasks', create);
router.put('/projects/:projectId/tasks/:taskId', update);
router.delete('/projects/:projectId/tasks/:taskId', remove);

module.exports = router;
