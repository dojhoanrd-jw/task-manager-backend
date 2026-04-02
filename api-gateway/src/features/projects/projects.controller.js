const proxy = require('../../shared/utils/proxy');

const getByUser = async (req, res, next) => {
  try {
    const { data } = await proxy.get('/projects', req);
    res.json(data);
  } catch (error) {
    next(error);
  }
};

const getById = async (req, res, next) => {
  try {
    const { data } = await proxy.get(`/projects/${req.params.projectId}`, req);
    res.json(data);
  } catch (error) {
    next(error);
  }
};

const create = async (req, res, next) => {
  try {
    const { data } = await proxy.post('/projects', req);
    res.status(201).json(data);
  } catch (error) {
    next(error);
  }
};

const update = async (req, res, next) => {
  try {
    const { data } = await proxy.put(`/projects/${req.params.projectId}`, req);
    res.json(data);
  } catch (error) {
    next(error);
  }
};

const remove = async (req, res, next) => {
  try {
    const { data } = await proxy.delete(`/projects/${req.params.projectId}`, req);
    res.json(data);
  } catch (error) {
    next(error);
  }
};

const addMember = async (req, res, next) => {
  try {
    const { data } = await proxy.post(`/projects/${req.params.projectId}/members`, req);
    res.json(data);
  } catch (error) {
    next(error);
  }
};

const removeMember = async (req, res, next) => {
  try {
    const { projectId, userId } = req.params;
    const { data } = await proxy.delete(`/projects/${projectId}/members/${userId}`, req);
    res.json(data);
  } catch (error) {
    next(error);
  }
};

module.exports = { getByUser, getById, create, update, remove, addMember, removeMember };
