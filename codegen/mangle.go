package codegen

import "fmt"

func (c *Context) mangle(s string) string {
	return fmt.Sprintf("%s.%s", c.moduleName, s)
}

func (c *Context) mainFuncName() string {
	return c.moduleName + ".main"
}
