package file_operations

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
)

// Generate function calling statement by funName and arguments
func GenerateCallExpr(funName string, args []string, isNew bool, assertType string) string {
	lenOfArgs := len(args)
	if isNew {
		funName += "("
	} else {
		funName += "\\("
	}
	for index, arg := range args {
		funName += arg
		if index != lenOfArgs-1 {
			funName += ", "
		}
	}
	if isNew {
		funName += ")"
	} else {
		funName += "\\)"
		lenOfAssert := len(assertType)
		if lenOfAssert != 0 {
			funName += "\\.\\("
			funName += assertType
			funName += "\\)"
		}
	}
	return funName
}

func ReplaceOriginFuncByFile(file, origin, target string) {
	output, needHandle, err := readFile(file, origin, target)
	if err != nil {
		panic(err)
	}
	if needHandle {
		err = writeCallExprToFile(file, output)
		if err != nil {
			panic(err)
		}
		fmt.Println(origin, "has been replaced with", target)

	} else {
		fmt.Println("Can't find ", origin)
	}
}

// Read the file line by line to match origin and replace by target
func readFile(filePath, origin, target string) ([]byte, bool, error) {
	f, err := os.OpenFile(filePath, os.O_RDONLY, 0644)
	if err != nil {
		return nil, false, err
	}
	defer f.Close()
	reader := bufio.NewReader(f)
	needHandle := false
	output := make([]byte, 0)
	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				return output, needHandle, nil
			}
			return nil, needHandle, err
		}

		if ok, _ := regexp.Match(origin, line); ok {
			fmt.Println("Statement match success!")
			reg := regexp.MustCompile(origin)
			newByte := reg.ReplaceAll(line, []byte(target))
			output = append(output, newByte...)
			output = append(output, []byte("\n")...)
			if !needHandle {
				needHandle = true
			}
		} else {
			output = append(output, line...)
			output = append(output, []byte("\n")...)
		}
	}
	return output, needHandle, nil
}

// Write target function calling statement to the file
func writeCallExprToFile(filePath string, input []byte) error {
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0600)
	defer f.Close()
	if err != nil {
		return err
	}
	writer := bufio.NewWriter(f)
	_, err = writer.Write(input)
	if err != nil {
		return err
	}
	writer.Flush()
	return nil
}
