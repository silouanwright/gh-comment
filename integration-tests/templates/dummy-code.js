// Integration Test File - Contains intentional issues for commenting
function calculateTotal(items) {
    let total = 0;
    for (let i = 0; i < items.length; i++) {
        total += items[i].price * items[i].quantity; // Potential null pointer
    }
    return total; // Missing input validation
}

// TODO: Add error handling
// FIXME: Handle empty arrays
const processOrder = (order) => {
    const total = calculateTotal(order.items);
    return { total, tax: total * 0.08 }; // Hardcoded tax rate
};

// Additional test scenarios
const validateInput = (data) => {
    // Missing validation logic
    return data;
};

const formatCurrency = (amount) => {
    return "$" + amount.toFixed(2); // Assumes USD
};

module.exports = { calculateTotal, processOrder, validateInput, formatCurrency };
