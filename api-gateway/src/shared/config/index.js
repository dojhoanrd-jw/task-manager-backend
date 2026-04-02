const dotenv = require('dotenv');
dotenv.config();

const config = {
  port: process.env.PORT || 8080,
  taskServiceUrl: process.env.TASK_SERVICE_URL || 'http://localhost:8081',
  jwtSecret: process.env.JWT_SECRET || '',
};

if (!config.jwtSecret) {
  console.error('JWT_SECRET is required');
  process.exit(1);
}

module.exports = config;
