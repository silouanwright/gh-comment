// Integration test file with intentional issues for commenting
function processUserData(users) {
    // Missing input validation
    let results = [];

    // SQL injection vulnerability
    const query = "SELECT * FROM users WHERE id = " + users[0].id;

    // Hardcoded secrets
    const apiKey = "sk-1234567890abcdef";

    // Performance issue - nested loops
    for (let i = 0; i < users.length; i++) {
        for (let j = 0; j < users.length; j++) {
            if (users[i].status === users[j].status) {
                results.push(users[i]);
            }
        }
    }

    return results;
}

module.exports = { processUserData };
