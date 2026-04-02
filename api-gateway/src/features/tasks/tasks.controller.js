const proxy = require('../../shared/utils/proxy');

const getByProject = async (req, res, next) => {
  try {
    const { projectId } = req.params;
    const { data } = await proxy.get(`/projects/${projectId}/tasks`, req, {
      limit: req.query.limit,
      lastId: req.query.lastId,
    });
    res.json(data);
  } catch (error) {
    next(error);
  }
};

const getById = async (req, res, next) => {
  try {
    const { projectId, taskId } = req.params;
    const { data } = await proxy.get(`/projects/${projectId}/tasks/${taskId}`, req);
    res.json(data);
  } catch (error) {
    next(error);
  }
};

const create = async (req, res, next) => {
  try {
    const { projectId } = req.params;
    const { data } = await proxy.post(`/projects/${projectId}/tasks`, req);
    res.status(201).json(data);
  } catch (error) {
    next(error);
  }
};

const update = async (req, res, next) => {
  try {
    const { projectId, taskId } = req.params;
    const { data } = await proxy.put(`/projects/${projectId}/tasks/${taskId}`, req);
    res.json(data);
  } catch (error) {
    next(error);
  }
};

const remove = async (req, res, next) => {
  try {
    const { projectId, taskId } = req.params;
    const { data } = await proxy.delete(`/projects/${projectId}/tasks/${taskId}`, req);
    res.json(data);
  } catch (error) {
    next(error);
  }
};

module.exports = { getByProject, getById, create, update, remove };
