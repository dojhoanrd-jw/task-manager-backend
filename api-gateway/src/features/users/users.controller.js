const httpClient = require('../../shared/utils/http-client');

// Get all users (admin only)
const getAll = async (req, res, next) => {
  try {
    const { data } = await httpClient.get('/users', {
      headers: { Authorization: req.token },
    });

    res.json(data);
  } catch (error) {
    next(error);
  }
};

// Update user role (admin only)
const updateRole = async (req, res, next) => {
  try {
    const { userId } = req.params;

    const { data } = await httpClient.put(`/users/${userId}/role`, req.body, {
      headers: { Authorization: req.token },
    });

    res.json(data);
  } catch (error) {
    next(error);
  }
};

module.exports = { getAll, updateRole };
