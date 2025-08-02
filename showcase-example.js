// Showcase Example - Intentionally contains issues for gh-comment demonstration
// This file demonstrates various code patterns that benefit from line-specific commenting

function calculateUserMetrics(users) {
    let totalRevenue = 0;
    let activeUsers = 0;

    // Process each user
    for (let i = 0; i < users.length; i++) {
        const user = users[i];

        // Revenue calculation - has potential issues
        if (user.purchases) {
            for (let j = 0; j < user.purchases.length; j++) {
                totalRevenue += user.purchases[j].amount;
            }
        }

        // Activity check - could be improved
        if (user.lastLogin > Date.now() - (30 * 24 * 60 * 60 * 1000)) {
            activeUsers++;
        }
    }

    // Return metrics object
    return {
        totalRevenue: totalRevenue,
        activeUsers: activeUsers,
        averageRevenue: totalRevenue / users.length
    };
}

// Database query function - has security concerns
function getUserData(userId) {
    const query = "SELECT * FROM users WHERE id = " + userId;  // SQL injection risk
    return database.query(query);
}

// Authentication middleware - needs improvement
function authenticateUser(req, res, next) {
    const token = req.headers.authorization;
    if (!token) {
        return res.status(401).send('Unauthorized');
    }

    // Token validation - hardcoded secret
    const secret = "my-secret-key-123";
    try {
        const decoded = jwt.verify(token, secret);
        req.user = decoded;
        next();
    } catch (error) {
        return res.status(401).send('Invalid token');
    }
}

// Data processing function - performance issues
function processLargeDataset(data) {
    let results = [];

    // Inefficient nested loops
    for (let i = 0; i < data.length; i++) {
        for (let j = 0; j < data.length; j++) {
            if (data[i].category === data[j].category && i !== j) {
                results.push({
                    item1: data[i],
                    item2: data[j],
                    similarity: calculateSimilarity(data[i], data[j])
                });
            }
        }
    }

    return results;
}

// Helper function
function calculateSimilarity(item1, item2) {
    // Simplified similarity calculation
    return Math.random();  // Obviously not a real implementation
}

module.exports = {
    calculateUserMetrics,
    getUserData,
    authenticateUser,
    processLargeDataset
};
