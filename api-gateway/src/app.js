const express = require('express');
const cors = require('cors');
const morgan = require('morgan');
const config = require('./shared/config');

const app = express();

// Global middleware
app.use(cors());
app.use(morgan('dev'));
app.use(express.json());

// Health check
app.get('/health', (req, res) => {
  res.json({ status: 'ok', service: 'api-gateway' });
});

// Routes
const authRoutes = require('./features/auth/auth.routes');
app.use('/api/auth', authRoutes);

// Start server
app.listen(config.port, () => {
  console.log(`API Gateway running on port ${config.port}`);
});

module.exports = app;
