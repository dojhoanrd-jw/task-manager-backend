const { Router } = require('express');
const validate = require('../../shared/middleware/validate.middleware');
const { registerSchema, loginSchema } = require('./auth.validation');
const { register, login } = require('./auth.controller');

const router = Router();

router.post('/register', validate(registerSchema), register);
router.post('/login', validate(loginSchema), login);

module.exports = router;
