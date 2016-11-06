package codegen

func (c *Context) mangle(s string) string {
	return "gray." + s
}
