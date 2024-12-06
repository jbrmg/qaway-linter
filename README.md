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
        functions:
          - targets:
              - packages: [ "github.com/myrepo/mypkg" ]
                minLinesOfCode: 10
            params:
              minCommentDensity: 0.1
              requireHeadlineComment: true
              minHeadlineDensity: 0.05
              trivialCommentThreshold: 0.5
              minLoggingDensity: 0.05
        interfaces:
          - targets:
              - packages: [ "github.com/myrepo/mypkg" ]
            params:
              requireHeadlineComments: true
              requireMethodComments: true
```

4. Execute the custom version by running `./custom-gcl run` in your project's root directory.

## Exclusions

Add `// nolint:qawaylinter` to the line you want to exclude from the linter.

