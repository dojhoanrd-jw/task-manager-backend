const Joi = require('joi');

const createProjectSchema = Joi.object({
  name: Joi.string().min(1).max(200).required(),
  description: Joi.string().max(1000).allow(''),
});

const updateProjectSchema = Joi.object({
  name: Joi.string().min(1).max(200),
  description: Joi.string().max(1000).allow(''),
}).min(1);

const addMemberSchema = Joi.object({
  userId: Joi.string().required(),
});

module.exports = { createProjectSchema, updateProjectSchema, addMemberSchema };
