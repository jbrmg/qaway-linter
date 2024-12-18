# QAway Linter

## Usage

1. Create a file called `.custom-gcl.yml` in your projects root directy with the following content:

```yaml
version: v1.62.2 # TODO: update to latest version (see https://github.com/golangci/golangci-lint/releases/tag/v1.62.2)
plugins:
  # a plugin from local source
  - module: 'qawaylinter'
    path: /Users/jan.bender/Documents/qaware/qalabs/qaway-linter
```

2. Execute `golangci-lint custom` to build a custom version of golangci-lint with the qawaylinter plugin.

3. Add configuration to your `.golangci.yml`:

```yaml
linters:
  enable:
    # add to existing list if linters.disable-all is set to true
    - qawaylinter

linter-settings:
  custom:
    qawaylinter:
      type: "module"
      description: "Checks for appropriate documentation in code"
      settings:
        rules:
          - packages: [ "github.com/myorg/myrepo" ]
            # This rule demonstrates all available configuration options
            # If a parameter is not set, it is not enforced.
            functions:
              filters:
                # Apply parameters only to functions with at least 10 lines of code
                minLinesOfCode: 10
              params:
                # A method must have at least 10% of comments (headline + inline) compared to its lines of code
                minCommentDensity: 0.1
                # A headline comment is required for every method
                requireHeadlineComment: true
                # Trivial comments (similarity to method name) are not allowed. 
                # The threshold indicates the similarity to the method name.
                # A higher threshold indicates a higher similarity, resulting in less warnings.
                trivialCommentThreshold: 0.5
                # Amount of logging statements compared to lines of code. 
                minLoggingDensity: 0.0
            interfaces:
              params:
                # A headline comment is required for every interface
                requireHeadlineComment: true
                # A comment is required for every method in an interface
                requireMethodComment: true
            structs:
              params:
                # A headline comment is required for every struct
                requireHeadlineComment: true
                # A comment is required for every field in a struct
                requireFieldComment: false
          - packages: [ "github.com/myorg/myrepo/subpkg" ] # rules for subpackage override super packages
            functions:
              filters:
                minLinesOfCode: 20
              params:
                trivialCommentThreshold: 0.5
                minLoggingDensity: 0.1
```

4. Execute the custom version by running `./custom-gcl run` in your project's root directory.

## Exclusions

Add `// nolint:qawaylinter` to the line you want to exclude from the linter.

