const httpClient = require('../../shared/utils/http-client');

// Register a new user
const register = async (req, res, next) => {
  try {
    const { data } = await httpClient.post('/auth/register', req.body);
    res.status(201).json(data);
  } catch (error) {
    next(error);
  }
};

// Login user
const login = async (req, res, next) => {
  try {
    const { data } = await httpClient.post('/auth/login', req.body);
    res.status(200).json(data);
  } catch (error) {
    next(error);
  }
};

module.exports = { register, login };
