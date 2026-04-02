const { Router } = require('express');
const authMiddleware = require('../../shared/middleware/auth.middleware');
const validate = require('../../shared/middleware/validate.middleware');
const validateParams = require('../../shared/middleware/params.middleware');
const { createTaskSchema, updateTaskSchema } = require('./tasks.validation');
const { getByProject, getById, create, update, remove } = require('./tasks.controller');

const router = Router();

router.use(authMiddleware);

router.get('/projects/:projectId/tasks', validateParams('projectId'), getByProject);
router.get('/projects/:projectId/tasks/:taskId', validateParams('projectId', 'taskId'), getById);
router.post('/projects/:projectId/tasks', validateParams('projectId'), validate(createTaskSchema), create);
router.put('/projects/:projectId/tasks/:taskId', validateParams('projectId', 'taskId'), validate(updateTaskSchema), update);
router.delete('/projects/:projectId/tasks/:taskId', validateParams('projectId', 'taskId'), remove);

module.exports = router;
