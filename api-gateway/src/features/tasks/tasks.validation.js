const Joi = require('joi');

const createTaskSchema = Joi.object({
  title: Joi.string().min(1).max(200).required(),
  description: Joi.string().max(1000).allow(''),
  assignedTo: Joi.string().allow(''),
});

const updateTaskSchema = Joi.object({
  title: Joi.string().min(1).max(200),
  description: Joi.string().max(1000).allow(''),
  completed: Joi.boolean(),
  assignedTo: Joi.string().allow(''),
}).min(1);

module.exports = { createTaskSchema, updateTaskSchema };
