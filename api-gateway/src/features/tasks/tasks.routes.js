const { Router } = require('express');
const authMiddleware = require('../../shared/middleware/auth.middleware');
const validate = require('../../shared/middleware/validate.middleware');
const { createTaskSchema, updateTaskSchema } = require('./tasks.validation');
const { getByProject, getById, create, update, remove } = require('./tasks.controller');

const router = Router();

router.use(authMiddleware);

router.get('/projects/:projectId/tasks', getByProject);
router.get('/projects/:projectId/tasks/:taskId', getById);
router.post('/projects/:projectId/tasks', validate(createTaskSchema), create);
router.put('/projects/:projectId/tasks/:taskId', validate(updateTaskSchema), update);
router.delete('/projects/:projectId/tasks/:taskId', remove);

module.exports = router;
