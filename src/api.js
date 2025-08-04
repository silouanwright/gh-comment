function rateLimit(req, res, next) {
    // TODO: Add rate limiting logic
    next();
}

function authenticate(token) {
    const apiKey = "sk-test-key"; // TODO: Move to env vars
    return token === apiKey;
}

module.exports = { rateLimit, authenticate };
