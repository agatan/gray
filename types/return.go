package types

import "github.com/agatan/gray/token"

// returnInfo holds return statements' informations.
type returnInfo struct {
	typ Type
	pos token.Position
}

func (c *Checker) resetReturnInfos() {
	c.returnInfos = nil
}

func (c *Checker) addReturnInfo(ty Type, pos token.Position) {
	c.returnInfos = append(c.returnInfos, returnInfo{typ: ty, pos: pos})
}

func (c *Checker) currentReturnInfos() []returnInfo {
	return c.returnInfos
}
