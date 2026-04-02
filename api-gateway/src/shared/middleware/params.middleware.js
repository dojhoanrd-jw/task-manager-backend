// Validates that route params are non-empty alphanumeric strings (Firestore document IDs)
const firestoreIdPattern = /^[a-zA-Z0-9]{1,128}$/;

const validateParams = (...paramNames) => {
  return (req, res, next) => {
    for (const param of paramNames) {
      const value = req.params[param];
      if (!value || !firestoreIdPattern.test(value)) {
        return res.status(400).json({ error: `invalid ${param} format` });
      }
    }
    next();
  };
};

module.exports = validateParams;
