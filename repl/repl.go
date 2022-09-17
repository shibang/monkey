package repl

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/shibang/monkey/evaluator"
	"github.com/shibang/monkey/lexer"
	"github.com/shibang/monkey/object"
	"github.com/shibang/monkey/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()
	macroEnv := object.NewEnvironment()

	for {
		fmt.Fprint(out, PROMPT)
		line := readLine(scanner)
		if line == "" {
			return
		}
		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		evaluator.DefineMacros(program, macroEnv)
		expanded := evaluator.ExpandMacros(program, macroEnv)

		evaluated := evaluator.Eval(expanded, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}

func readLine(scanner *bufio.Scanner) string {
	var out bytes.Buffer
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasSuffix(line, "\\") {
			out.WriteString(line)
			break
		}
		out.WriteString(strings.TrimSuffix(line, "\\"))
	}
	return out.String()
}
