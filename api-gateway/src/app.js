const express = require('express');
const cors = require('cors');
const morgan = require('morgan');
const config = require('./shared/config');

const app = express();

// Global middleware
const requestIdMiddleware = require('./shared/middleware/request-id.middleware');
app.use(cors());
app.use(requestIdMiddleware);
app.use(morgan(':method :url :status :response-time ms - :req[x-request-id]'));
app.use(express.json());

// Health check
app.get('/health', (req, res) => {
  res.json({ status: 'ok', service: 'api-gateway' });
});

// Routes
const authRoutes = require('./features/auth/auth.routes');
const taskRoutes = require('./features/tasks/tasks.routes');
const projectRoutes = require('./features/projects/projects.routes');
const userRoutes = require('./features/users/users.routes');

app.use('/api/auth', authRoutes);
app.use('/api', taskRoutes);
app.use('/api/projects', projectRoutes);
app.use('/api/users', userRoutes);

// Error handling
const errorMiddleware = require('./shared/middleware/error.middleware');
app.use(errorMiddleware);

// Start server
app.listen(config.port, () => {
  console.log(`API Gateway running on port ${config.port}`);
});

module.exports = app;
