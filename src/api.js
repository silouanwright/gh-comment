const express = require('express');
const rateLimit = require('express-rate-limit');

// TODO: Add comprehensive rate limiting
const limiter = rateLimit({
    windowMs: 15 * 60 * 1000, // 15 minutes
    max: 100, // limit each IP to 100 requests per windowMs
    message: 'Too many requests from this IP'
});

function authenticate(req, res, next) {
    const token = req.headers.authorization;
    // TODO: Implement proper JWT validation
    if (!token) {
        return res.status(401).json({ error: 'No token provided' });
    }
    next();
}

// Middleware with potential security issue for testing
function processUserInput(input) {
    // SECURITY: This should be sanitized
    return input.toLowerCase();
}

module.exports = { limiter, authenticate, processUserInput };
