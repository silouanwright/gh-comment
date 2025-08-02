package cmd

import (
	"strconv"
	"testing"
	"unicode/utf8"
)

// FuzzCommentID tests comment ID parsing with random inputs
func FuzzCommentID(f *testing.F) {
	// Seed with known good and bad inputs
	f.Add("123456")
	f.Add("0")
	f.Add("-1")
	f.Add("invalid")
	f.Add("")
	f.Add("999999999999999999999") // Very large number
	f.Add("123abc")
	f.Add("abc123")
	f.Add(" 123 ")
	f.Add("123.456")

	f.Fuzz(func(t *testing.T, input string) {
		// Test that comment ID parsing doesn't panic
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Comment ID parsing panicked with input %q: %v", input, r)
			}
		}()

		// Test strconv.Atoi behavior (used in runReply and runResolve)
		commentID, err := strconv.Atoi(input)

		if err == nil {
			// If parsing succeeded, the ID should be reasonable
			if commentID < 0 {
				// Negative IDs are invalid for GitHub
				t.Logf("Negative comment ID parsed: %d from input %q", commentID, input)
			}
			if commentID == 0 {
				// Zero IDs are typically invalid for GitHub
				t.Logf("Zero comment ID parsed from input %q", input)
			}
		}

		// Test formatValidationError doesn't panic with this input
		validationErr := formatValidationError("comment ID", input, "must be a valid integer")
		if validationErr == nil {
			t.Errorf("formatValidationError returned nil for input %q", input)
		}

		// Ensure error message is valid UTF-8
		errMsg := validationErr.Error()
		if !utf8.ValidString(errMsg) {
			t.Errorf("Error message contains invalid UTF-8 for input %q: %q", input, errMsg)
		}
	})
}

// FuzzReactionValidation tests reaction validation with random inputs
func FuzzReactionValidation(f *testing.F) {
	// Seed with known reactions and variations
	f.Add("+1")
	f.Add("-1")
	f.Add("laugh")
	f.Add("confused")
	f.Add("heart")
	f.Add("hooray")
	f.Add("rocket")
	f.Add("eyes")
	f.Add("invalid")
	f.Add("")
	f.Add("LAUGH")
	f.Add("+2")
	f.Add("👍")
	f.Add("thumbsup")
	f.Add(" +1 ")

	f.Fuzz(func(t *testing.T, reaction string) {
		// Test that reaction validation doesn't panic
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Reaction validation panicked with input %q: %v", reaction, r)
			}
		}()

		// Test validateReaction function
		isValid := validateReaction(reaction)

		// Known valid reactions should always be valid
		validReactions := []string{"+1", "-1", "laugh", "confused", "heart", "hooray", "rocket", "eyes"}
		for _, valid := range validReactions {
			if reaction == valid && !isValid {
				t.Errorf("Valid reaction %q was marked as invalid", reaction)
			}
		}

		// Test formatValidationError with this reaction
		if !isValid {
			validationErr := formatValidationError("reaction", reaction, "must be one of: +1, -1, laugh, confused, heart, hooray, rocket, eyes")
			if validationErr == nil {
				t.Errorf("formatValidationError returned nil for invalid reaction %q", reaction)
			}

			// Ensure error message is valid UTF-8
			errMsg := validationErr.Error()
			if !utf8.ValidString(errMsg) {
				t.Errorf("Error message contains invalid UTF-8 for reaction %q: %q", reaction, errMsg)
			}
		}
	})
}

// FuzzAuthorFilter tests author filtering with random inputs
func FuzzAuthorFilter(f *testing.F) {
	// Seed with typical GitHub usernames and edge cases
	f.Add("octocat")
	f.Add("user-123")
	f.Add("user_name")
	f.Add("User.Name")
	f.Add("")
	f.Add("a")
	f.Add("very-long-username-that-might-cause-issues")
	f.Add("123")
	f.Add("user@domain.com")
	f.Add("user name") // Space in username
	f.Add("用户")        // Unicode characters

	f.Fuzz(func(t *testing.T, authorInput string) {
		// Test that author filtering doesn't panic
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Author filtering panicked with input %q: %v", authorInput, r)
			}
		}()

		// Create test comments with various authors
		testComments := []Comment{
			{ID: 1, Author: "alice", Body: "Test comment 1"},
			{ID: 2, Author: "bob", Body: "Test comment 2"},
			{ID: 3, Author: "charlie", Body: "Test comment 3"},
			{ID: 4, Author: "", Body: "Anonymous comment"},
			{ID: 5, Author: "用户", Body: "Unicode author"},
		}

		// Set global author filter
		originalAuthor := author
		author = authorInput
		defer func() { author = originalAuthor }()

		// Test filterComments function
		filtered := filterComments(testComments)

		// Note: Go returns nil slice when no elements are appended, this is normal behavior
		// We should check len(filtered) instead of nil check

		// If author filter is empty, should return all comments
		if authorInput == "" && len(filtered) != len(testComments) {
			t.Errorf("Empty author filter should return all comments, got %d, expected %d", len(filtered), len(testComments))
		}

		// If author filter is set, should only return matching comments
		if authorInput != "" {
			// Check that all returned comments match the filter using our enhanced matching
			for _, comment := range filtered {
				if !matchesAuthorFilter(comment.Author, authorInput) {
					t.Errorf("Filtered comment has author %q, doesn't match filter %q", comment.Author, authorInput)
				}
			}

			// It's OK if no comments match the filter (empty slice)
			// This is expected behavior for non-matching authors
		}

		// Ensure all returned comments are valid
		for _, comment := range filtered {
			if comment.ID <= 0 {
				t.Errorf("Invalid comment ID %d in filtered results", comment.ID)
			}
			if !utf8.ValidString(comment.Author) {
				t.Errorf("Invalid UTF-8 in author name: %q", comment.Author)
			}
			if !utf8.ValidString(comment.Body) {
				t.Errorf("Invalid UTF-8 in comment body: %q", comment.Body)
			}
		}
	})
}

// FuzzSuggestionExpansion tests suggestion syntax expansion with random inputs
func FuzzSuggestionExpansion(f *testing.F) {
	// Seed with various suggestion formats
	f.Add("```suggestion\nfixed code\n```")
	f.Add("```suggestion\n```")
	f.Add("```\nsuggestion\n```")
	f.Add("```suggestion")
	f.Add("suggestion```")
	f.Add("```suggestion\nline1\nline2\n```")
	f.Add("Before\n```suggestion\nfix\n```\nAfter")
	f.Add("```suggestion\n\n```")   // Empty suggestion
	f.Add("```suggestion\n\t\n```") // Whitespace only

	f.Fuzz(func(t *testing.T, input string) {
		// Test that suggestion expansion doesn't panic
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Suggestion expansion panicked with input %q: %v", input, r)
			}
		}()

		// Test expandSuggestions function
		result := expandSuggestions(input)

		// Ensure result is valid UTF-8
		if !utf8.ValidString(result) {
			t.Errorf("expandSuggestions returned invalid UTF-8 for input %q", input)
		}

		// Result should never be nil
		if result == "" && input != "" {
			t.Errorf("expandSuggestions returned empty string for non-empty input %q", input)
		}

		// If input has no suggestions, result should be identical
		if !containsSuggestion(input) && result != input {
			t.Errorf("expandSuggestions modified input without suggestions: %q -> %q", input, result)
		}
	})
}

// containsSuggestion checks if input contains suggestion syntax
func containsSuggestion(input string) bool {
	return len(input) > 0 && (
	// Look for suggestion markers
	containsSubstring(input, "```suggestion") ||
		containsSubstring(input, "suggestion```"))
}

// containsSubstring is a simple substring check
func containsSubstring(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr ||
			(len(s) > len(substr) &&
				(s[:len(substr)] == substr ||
					s[len(s)-len(substr):] == substr ||
					containsAt(s, substr))))
}

// containsAt checks if substr exists anywhere in s
func containsAt(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

// FuzzPRNumber tests PR number parsing with random inputs
func FuzzPRNumber(f *testing.F) {
	// Seed with various PR number formats
	f.Add("1")
	f.Add("123")
	f.Add("999999")
	f.Add("0")
	f.Add("-1")
	f.Add("invalid")
	f.Add("")
	f.Add("1.0")
	f.Add("1e10")
	f.Add("0x123")

	f.Fuzz(func(t *testing.T, input string) {
		// Test that PR number parsing doesn't panic
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("PR number parsing panicked with input %q: %v", input, r)
			}
		}()

		// Test strconv.Atoi behavior (used in runList)
		prNumber, err := strconv.Atoi(input)

		if err == nil {
			// If parsing succeeded, the PR number should be reasonable
			if prNumber <= 0 {
				// Non-positive PR numbers are invalid
				t.Logf("Non-positive PR number parsed: %d from input %q", prNumber, input)
			}
			if prNumber > 999999 {
				// Very large PR numbers might be suspicious
				t.Logf("Very large PR number parsed: %d from input %q", prNumber, input)
			}
		}

		// Test formatValidationError doesn't panic with this input
		validationErr := formatValidationError("PR number", input, "must be a valid integer")
		if validationErr == nil {
			t.Errorf("formatValidationError returned nil for input %q", input)
		}

		// Ensure error message is valid UTF-8
		errMsg := validationErr.Error()
		if !utf8.ValidString(errMsg) {
			t.Errorf("Error message contains invalid UTF-8 for input %q: %q", input, errMsg)
		}
	})
}
