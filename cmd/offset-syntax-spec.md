# GitLab-style Line Offset Syntax Specification

## Overview
This specification defines the line offset syntax for suggestions that should be applied to lines relative to the current comment line position.

## Syntax
The offset syntax extends the existing `[SUGGEST: code]` format to include relative line positioning:

```
[SUGGEST:<offset>: code]
```

Where `<offset>` is:
- `+N` - Apply suggestion N lines below the current line
- `-N` - Apply suggestion N lines above the current line
- `0` or empty - Apply suggestion to the current line (default behavior)

## Examples

### Positive Offsets (Below Current Line)
- `[SUGGEST:+1: const x = 1;]` - Apply to line 1 below current
- `[SUGGEST:+2: return value;]` - Apply to line 2 below current
- `[SUGGEST:+10: }]` - Apply to line 10 below current

### Negative Offsets (Above Current Line)
- `[SUGGEST:-1: import fs from 'fs';]` - Apply to line 1 above current
- `[SUGGEST:-2: const config = {};]` - Apply to line 2 above current
- `[SUGGEST:-5: function setup() {]` - Apply to line 5 above current

### Zero or Default Offset
- `[SUGGEST:0: fixed code]` - Apply to current line (explicit)
- `[SUGGEST: fixed code]` - Apply to current line (default)

## Parsing Rules

1. **Offset Detection**: Look for pattern `[SUGGEST:<offset>:` where offset matches `[+-]?\d+`
2. **Fallback**: If no offset is specified, default to current line (offset 0)
3. **Validation**: Offsets must be integers in range [-999, +999]
4. **Error Handling**: Invalid offsets should be treated as regular suggestions (no offset)

## Output Format
The GitHub suggestion block should include line range information when offset is specified:

### Current Line (no offset)
```suggestion
suggested code
```

### With Offset
```suggestion:+2
suggested code
```

Or for negative offsets:
```suggestion:-1
suggested code
```

## Implementation Notes

1. **Context Requirement**: Offset suggestions require knowledge of the current line number in the PR diff
2. **Validation**: Must validate that target line exists in the diff context
3. **Error Handling**: If target line is outside diff context, fall back to regular suggestion
4. **Backwards Compatibility**: Existing `[SUGGEST: code]` syntax remains unchanged

## Edge Cases

1. **Out of Bounds**: Offset points to line outside diff context
2. **Invalid Syntax**: Malformed offset (e.g., `[SUGGEST:+: code]`)
3. **Large Offsets**: Very large positive/negative offsets
4. **Zero Lines**: Empty files or no diff context

## Regex Pattern
```go
// Pattern to match [SUGGEST:<offset>: code] syntax
offsetSuggestionPattern := `\[SUGGEST:(([+-]?\d+):)?\s*(.*?)\]`
```

This captures:
1. Optional offset group: `([+-]?\d+):`
2. Code content: `(.*?)`
3. Handles both offset and non-offset cases