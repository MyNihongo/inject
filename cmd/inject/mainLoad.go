package main

import (
	"bufio"
	"context"
	"errors"
	"os"
	"path"
	"strings"
)

type loadResult struct {
	imports []*importStmt
	injects map[string]injectType
}

type importStmt struct {
	alias string
	path  string
}

// loadFileContent loads the container definition
func loadFileContent(ctx context.Context, wd, fileName string) (*loadResult, error) {
	filePath := path.Join(wd, fileName)
	injections := &loadResult{
		imports: make([]*importStmt, 0),
		injects: make(map[string]injectType),
	}

	if file, err := os.Open(filePath); err != nil {
		return nil, err
	} else {
		defer file.Close()

		scanner := bufio.NewScanner(file)

		isFuncFound, areImportsFound, bracketsLevel := false, false, 0
		var injectType injectType = notSet
		var sb strings.Builder

		for scanner.Scan() {
			if !isFuncFound {
				const importDecl = "import"
				if text := strings.TrimSpace(scanner.Text()); strings.EqualFold(text, "func BuildServiceProvider() {") {
					isFuncFound = true
				} else if areImportsFound {
					if text == ")" {
						areImportsFound = false
					} else {
						// block of imports
						if importStmt := createImport(text); importStmt != nil {
							injections.imports = append(injections.imports, importStmt)
						}
					}
				} else if len(injections.imports) == 0 && strings.HasPrefix(text, importDecl) {
					if strings.HasSuffix(text, "(") {
						areImportsFound = true
					} else {
						// single import
						if importStmt := createImport(text[len(importDecl)+1:]); importStmt != nil {
							injections.imports = append(injections.imports, importStmt)
						}
					}
				}

				continue
			}

			if text := strings.TrimSpace(scanner.Text()); text == "}" {
				return injections, nil
			} else {
				// TODO: support utf8?
				for i := 0; i < len(text); i++ {
					if injectType == notSet {
						if incr, ok := hasPrefix(text, "AddSingleton", i); ok {
							i += incr
							injectType = Singleton
						} else if incr, ok = hasPrefix(text, "AddTransient", i); ok {
							i += incr
							injectType = Transient
						}

						continue
					}

					if text[i] == '(' {
						bracketsLevel++

						// new block started
						if bracketsLevel == 1 {
							sb.Reset()
						} else {
							sb.WriteByte(text[i])
						}
					} else if text[i] == ')' {
						bracketsLevel--

						// current block ended
						if bracketsLevel == 0 {
							decl := sb.String()
							if lastIndex := len(decl) - 1; decl[lastIndex] == ',' {
								decl = decl[0:lastIndex]
							}

							injections.injects[decl] = injectType
							sb.Reset()
							injectType = notSet
						} else {
							sb.WriteByte(text[i])
						}
					} else {
						sb.WriteByte(text[i])
					}
				}
			}
		}

		if err = scanner.Err(); err != nil {
			return nil, err
		} else {
			return nil, errors.New("container declaration not found")
		}
	}
}

func hasPrefix(text, prefix string, startIndex int) (int, bool) {
	// TODO: support utf8?
	for i := 0; i < len(prefix) && startIndex+i < len(text); i++ {
		if prefix[i] != text[startIndex+i] {
			return 0, false
		}
	}

	// -1 for the array increment
	return len(prefix) - 1, true
}

func createImport(text string) *importStmt {
	for i := 0; i < len(text); i++ {
		if text[i] == '"' {
			return &importStmt{
				alias: strings.TrimSpace(text[:i]),
				path: strings.TrimFunc(text[i+1:], func(r rune) bool {
					return r == ' ' || r == '"'
				}),
			}
		}
	}

	return nil
}
