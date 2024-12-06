package commentdensity

import (
	"flag"
	"go/ast"
	"go/token"
	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

//nolint:gochecknoglobals
var flagSet flag.FlagSet

//nolint:gochecknoglobals
var (
	minCommentDensity int
	minLinesOfCode    int
)

const (
	defaultMinCommentDensity = 10
	defaultMinLinesOfCode    = 10
)

//nolint:gochecknoinits
func init() {
	flagSet.IntVar(&minCommentDensity, "minCommentDensity", defaultMinCommentDensity, "percentage of comments required for functions in relation to the method length")
	flagSet.IntVar(&minLinesOfCode, "minLinesOfCode", defaultMinLinesOfCode, "minimum lines of codes for methods to require a comment")
}

var Analyzer = &analysis.Analyzer{
	Name:     "gocommentdensity",
	Doc:      "Checks that a given function has an appropriate amount of comments.",
	Run:      run,
	Requires: []*analysis.Analyzer{inspect.Analyzer},
	Flags:    flagSet,
}

func run(pass *analysis.Pass) (interface{}, error) {
	var commentsInFile []*ast.CommentGroup
	inspect := func(node ast.Node) bool {
		funcDecl, ok := node.(*ast.FuncDecl)
		if !ok {
			return true
		}

		linesInFunction := countLinesInFunction(funcDecl, pass.Fset)
		if linesInFunction < minLinesOfCode {
			return true
		}

		linesOfComment := 0
		if funcDecl.Doc != nil {
			linesOfComment = countCommentLines(funcDecl.Doc, pass.Fset)
		}

		linesOfCommentsInMethodBody := determineInlineComments(funcDecl, commentsInFile, pass.Fset)
		linesOfComment += linesOfCommentsInMethodBody
		linesInFunction -= linesOfCommentsInMethodBody

		if float64(linesOfComment)/float64(linesInFunction)*100 < float64(minCommentDensity) {
			pass.Reportf(node.Pos(), "function '%s' should not have enough coverage. Lines of code: %d, expected lines of comment: %d", funcDecl.Name.Name, linesInFunction, linesOfComment)
		}

		return true
	}

	for _, f := range pass.Files {
		commentsInFile = f.Comments
		ast.Inspect(f, inspect)
	}
	return nil, nil
}

// determineInlineComments determines the number of lines of comments that are part of the method body.
// These comments are not returned as part of the AST of a FuncDecl.
// But all comments within a given file are available in the file's comments.
// This function determines the number of lines of comment within a method body by checking the comments in the file.
func determineInlineComments(f *ast.FuncDecl, commentsInFile []*ast.CommentGroup, fset *token.FileSet) int {
	commentLines := 0
	for _, comment := range commentsInFile {
		if (comment.Pos() >= f.Pos()) && (comment.End() <= f.End()) {
			commentLines += countCommentLines(comment, fset)
		}
	}
	return commentLines
}

// countLinesInFunction counts the lines between the start and end of a given function declaration.
// Note that this method also includes comments.
func countLinesInFunction(funcDecl *ast.FuncDecl, fset *token.FileSet) int {
	return fset.Position(funcDecl.End()).Line - fset.Position(funcDecl.Pos()).Line
}

// countCommentLines counts the lines covered by comments in a given AST node.
// This method takes into account that a command can span multiple lines using the /* */ syntax.
func countCommentLines(node ast.Node, fset *token.FileSet) int {
	commentLines := 0

	// Traverse the node to find all comment groups
	ast.Inspect(node, func(n ast.Node) bool {
		if commentGroup, ok := n.(*ast.CommentGroup); ok {
			for _, comment := range commentGroup.List {
				start := fset.Position(comment.Pos()).Line
				end := fset.Position(comment.End()).Line
				commentLines += end - start + 1
			}
		}
		return true
	})

	return commentLines
}
