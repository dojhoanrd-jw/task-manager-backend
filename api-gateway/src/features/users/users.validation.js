const Joi = require('joi');

const updateRoleSchema = Joi.object({
  role: Joi.string().valid('admin', 'member', 'viewer').required(),
});

module.exports = { updateRoleSchema };
