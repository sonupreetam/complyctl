# Go Project Style Guide

This style guide outlines the best practices to ensure consistency and readability across the codebase.

## General Guidelines

- **File Naming**: Use lowercase letters and underscores for file names (e.g., `my_file.go`).
- **Package Names**: Use short, concise, and lowercase names for packages. Avoid underscores and mixed caps.
- **Imports**: Group standard library imports separately from third-party imports and project-specific imports. Use blank lines to separate these groups.
- **Error Handling**: Always check for errors and handle them appropriately. Return errors to the caller when necessary.
- **Testing**: Write tests for your code. Use descriptive names for test functions and include edge cases.

## Code Formatting

- **Indentation**: Use tabs for indentation.
- **Line Length**: Limit lines to 99 characters when possible.
- **Braces**: Place opening braces on the same line as the statement (e.g., `if`, `for`, `func`).
- **Spacing**: Use a single space after commas, colons, and semicolons. Do not add spaces before commas, colons, and semicolons.

## Additional Guidelines

- **Empty Line at End of File**: Ensure that all files include an empty line at the end. This helps with version control diffs and adheres to POSIX standards.

By following these guidelines, we can maintain a clean, readable, and maintainable codebase.
