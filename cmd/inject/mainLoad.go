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
	injects map[string]injectType
}

// loadFileContent loads the container definition
func loadFileContent(ctx context.Context, wd, fileName string) (*loadResult, error) {
	filePath := path.Join(wd, fileName)
	injections := &loadResult{
		injects: make(map[string]injectType),
	}

	if file, err := os.Open(filePath); err != nil {
		return nil, err
	} else {
		defer file.Close()

		scanner := bufio.NewScanner(file)

		isFuncFound, bracketsLevel := false, 0
		var injectType injectType = notSet
		var sb strings.Builder

		for scanner.Scan() {
			if !isFuncFound {
				if strings.EqualFold(scanner.Text(), "func BuildServiceProvider() {") {
					isFuncFound = true
				}

				continue
			}

			if text := strings.TrimSpace(scanner.Text()); text == "}" {
				return injections, nil
			} else {
				// TODO: support utf8?
				for i := 0; i < len(text); i++ {
					if injectType == notSet {
						if incr, ok := HasPrefix(text, "AddSingleton", i); ok {
							i += incr
							injectType = Singleton
						} else if incr, ok = HasPrefix(text, "AddTransient", i); ok {
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

func HasPrefix(text, prefix string, startIndex int) (int, bool) {
	// TODO: support utf8?
	for i := 0; i < len(prefix) && startIndex+i < len(text); i++ {
		if prefix[i] != text[startIndex+i] {
			return 0, false
		}
	}

	// -1 for the array increment
	return len(prefix) - 1, true
}
