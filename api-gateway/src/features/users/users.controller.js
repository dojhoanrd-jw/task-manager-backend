const proxy = require('../../shared/utils/proxy');

const getAll = async (req, res, next) => {
  try {
    const { data } = await proxy.get('/users', req);
    res.json(data);
  } catch (error) {
    next(error);
  }
};

const updateRole = async (req, res, next) => {
  try {
    const { data } = await proxy.put(`/users/${req.params.userId}/role`, req);
    res.json(data);
  } catch (error) {
    next(error);
  }
};

module.exports = { getAll, updateRole };
