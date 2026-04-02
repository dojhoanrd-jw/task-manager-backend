const httpClient = require('../../shared/utils/http-client');

// Get tasks by project (paginated)
const getByProject = async (req, res, next) => {
  try {
    const { projectId } = req.params;
    const { limit, lastId } = req.query;

    const { data } = await httpClient.get(`/projects/${projectId}/tasks`, {
      headers: { Authorization: req.token },
      params: { limit, lastId },
    });

    res.json(data);
  } catch (error) {
    next(error);
  }
};

// Get task by ID
const getById = async (req, res, next) => {
  try {
    const { projectId, taskId } = req.params;

    const { data } = await httpClient.get(`/projects/${projectId}/tasks/${taskId}`, {
      headers: { Authorization: req.token },
    });

    res.json(data);
  } catch (error) {
    next(error);
  }
};

// Create task
const create = async (req, res, next) => {
  try {
    const { projectId } = req.params;

    const { data } = await httpClient.post(`/projects/${projectId}/tasks`, req.body, {
      headers: { Authorization: req.token },
    });

    res.status(201).json(data);
  } catch (error) {
    next(error);
  }
};

// Update task
const update = async (req, res, next) => {
  try {
    const { projectId, taskId } = req.params;

    const { data } = await httpClient.put(`/projects/${projectId}/tasks/${taskId}`, req.body, {
      headers: { Authorization: req.token },
    });

    res.json(data);
  } catch (error) {
    next(error);
  }
};

// Delete task
const remove = async (req, res, next) => {
  try {
    const { projectId, taskId } = req.params;

    const { data } = await httpClient.delete(`/projects/${projectId}/tasks/${taskId}`, {
      headers: { Authorization: req.token },
    });

    res.json(data);
  } catch (error) {
    next(error);
  }
};

module.exports = { getByProject, getById, create, update, remove };
