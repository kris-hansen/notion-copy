# notion-copy Development & Critique Pipeline

Two agentic loops: first develops the project, second critiques and improves it.

## Variables

```yaml
vars:
  MAX_DEV_ITERATIONS: "10"
  MAX_CRITIQUE_ITERATIONS: "5"
```

## Workflow

```yaml
steps:
  # Phase 1: Project Development Loop
  - name: project-developer
    type: loop
    provider: claude-code
    until: "implementation complete and all features from README.md are working"
    max_iterations: $MAX_DEV_ITERATIONS
    steps:
      - name: analyze_requirements
        type: prompt
        prompt: |
          Read README.md and analyze what needs to be built for notion-copy.
          
          Current state: Check what code exists in cmd/, internal/, and main.go.
          
          Create a prioritized implementation plan. Focus on:
          1. Core CLI structure (cobra/flags)
          2. Notion API client
          3. Pull command (Notion → Markdown)
          4. Push command (Markdown → Notion)
          5. Block type mapping
          
          Output a clear plan for the next implementation step.

      - name: implement_next_feature
        type: prompt
        prompt: |
          Based on the current codebase state, implement the next feature or fix.
          
          Guidelines:
          - Write idiomatic Go code
          - Add proper error handling
          - Create tests for new functionality
          - Update go.mod if new dependencies needed
          
          After implementing, run `go build` to verify it compiles.
          Run `go test ./...` to check tests pass.

      - name: check_progress
        type: prompt
        prompt: |
          Review the current implementation status:
          
          1. Does `go build` succeed?
          2. Do tests pass with `go test ./...`?
          3. Compare implemented features against README.md spec
          4. List what's complete vs what's remaining
          
          If all features are implemented and working, respond with "IMPLEMENTATION COMPLETE".
          Otherwise, describe what needs to be done next.

  # Phase 2: Critique & Improvement Loop  
  - name: project-critic
    type: loop
    provider: claude-code
    until: "all critical issues resolved and code quality is excellent"
    max_iterations: $MAX_CRITIQUE_ITERATIONS
    steps:
      - name: audit_codebase
        type: prompt
        prompt: |
          Perform a thorough code review of the notion-copy codebase:
          
          1. **Architecture**: Is the code well-organized? Proper separation of concerns?
          2. **Error Handling**: Are errors handled gracefully? Good user messages?
          3. **Edge Cases**: What happens with empty pages, rate limits, network errors?
          4. **Testing**: Is test coverage adequate? Any missing test cases?
          5. **Documentation**: Are public functions documented? Is usage clear?
          6. **Security**: API key handling? Input validation?
          7. **Performance**: Any obvious bottlenecks for large exports?
          
          Create a prioritized list of issues to fix.

      - name: fix_critical_issues
        type: prompt
        prompt: |
          Address the most critical issues found in the audit:
          
          - Fix bugs and edge cases
          - Improve error messages
          - Add missing tests
          - Refactor problematic code
          
          After each fix, verify with `go build` and `go test ./...`.

      - name: evaluate_quality
        type: prompt
        prompt: |
          Re-evaluate the codebase quality:
          
          1. Run `go vet ./...` and fix any issues
          2. Check `gofmt -d .` for formatting
          3. Review test coverage
          4. Verify all README features work correctly
          
          If code quality is now excellent with no critical issues, respond with "QUALITY APPROVED".
          Otherwise, describe remaining concerns.
```

## Execution

Run with:
```bash
comanda run ./comanda/initial-dev.md --verbose
```
