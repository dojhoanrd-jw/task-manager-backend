const express = require('express');
const cors = require('cors');
const helmet = require('helmet');
const morgan = require('morgan');
const config = require('./shared/config');

const app = express();

// Global middleware
const requestIdMiddleware = require('./shared/middleware/request-id.middleware');
app.use(helmet());
app.use(cors());
app.use(requestIdMiddleware);
app.use(morgan(':method :url :status :response-time ms - :req[x-request-id]'));
app.use(express.json({ limit: '1mb' }));

// Health check
app.get('/health', (req, res) => {
  res.json({ status: 'ok', service: 'api-gateway' });
});

// Rate limiting
const { apiLimiter, authLimiter } = require('./shared/middleware/rate-limit.middleware');

// Routes
const authRoutes = require('./features/auth/auth.routes');
const taskRoutes = require('./features/tasks/tasks.routes');
const projectRoutes = require('./features/projects/projects.routes');
const userRoutes = require('./features/users/users.routes');

app.use('/api/auth', authLimiter, authRoutes);
app.use('/api', apiLimiter, taskRoutes);
app.use('/api/projects', apiLimiter, projectRoutes);
app.use('/api/users', apiLimiter, userRoutes);

// Error handling
const errorMiddleware = require('./shared/middleware/error.middleware');
app.use(errorMiddleware);

// Start server
app.listen(config.port, () => {
  console.log(`API Gateway running on port ${config.port}`);
});

module.exports = app;
