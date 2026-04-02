const httpClient = require('../../shared/utils/http-client');

// Get projects for the authenticated user
const getByUser = async (req, res, next) => {
  try {
    const { data } = await httpClient.get('/projects', {
      headers: { Authorization: req.token },
    });

    res.json(data);
  } catch (error) {
    next(error);
  }
};

// Get project by ID
const getById = async (req, res, next) => {
  try {
    const { projectId } = req.params;

    const { data } = await httpClient.get(`/projects/${projectId}`, {
      headers: { Authorization: req.token },
    });

    res.json(data);
  } catch (error) {
    next(error);
  }
};

// Create project
const create = async (req, res, next) => {
  try {
    const { data } = await httpClient.post('/projects', req.body, {
      headers: { Authorization: req.token },
    });

    res.status(201).json(data);
  } catch (error) {
    next(error);
  }
};

// Update project
const update = async (req, res, next) => {
  try {
    const { projectId } = req.params;

    const { data } = await httpClient.put(`/projects/${projectId}`, req.body, {
      headers: { Authorization: req.token },
    });

    res.json(data);
  } catch (error) {
    next(error);
  }
};

// Delete project
const remove = async (req, res, next) => {
  try {
    const { projectId } = req.params;

    const { data } = await httpClient.delete(`/projects/${projectId}`, {
      headers: { Authorization: req.token },
    });

    res.json(data);
  } catch (error) {
    next(error);
  }
};

// Add member to project
const addMember = async (req, res, next) => {
  try {
    const { projectId } = req.params;

    const { data } = await httpClient.post(`/projects/${projectId}/members`, req.body, {
      headers: { Authorization: req.token },
    });

    res.json(data);
  } catch (error) {
    next(error);
  }
};

// Remove member from project
const removeMember = async (req, res, next) => {
  try {
    const { projectId, userId } = req.params;

    const { data } = await httpClient.delete(`/projects/${projectId}/members/${userId}`, {
      headers: { Authorization: req.token },
    });

    res.json(data);
  } catch (error) {
    next(error);
  }
};

module.exports = { getByUser, getById, create, update, remove, addMember, removeMember };
