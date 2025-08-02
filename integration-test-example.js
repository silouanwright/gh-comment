// Integration test file with intentional issues for commenting
function processUserData(users) {
    // Added input validation - security improvement!
    if (!users || !Array.isArray(users) || users.length === 0) {
        throw new Error('Invalid users array provided');
    }
    let results = [];

    // Fixed SQL injection vulnerability with parameterized query
    const query = "SELECT * FROM users WHERE id = ?";
    const dbResult = db.query(query, [users[0].id]);

    // Moved API key to environment variable for security
    const apiKey = process.env.API_KEY || "development-fallback-key";

    // Fixed performance issue - use Set for O(n) instead of O(n²)
    const statusSet = new Set();
    const statusCounts = {};

    // First pass: collect unique statuses
    for (let user of users) {
        if (!statusCounts[user.status]) {
            statusCounts[user.status] = [];
        }
        statusCounts[user.status].push(user);
    }

    // Second pass: add users with duplicate statuses
    for (let status in statusCounts) {
        if (statusCounts[status].length > 1) {
            results.push(...statusCounts[status]);
        }
    }

    return results;
}

module.exports = { processUserData };
