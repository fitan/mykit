package myctx

type CtxKey int

const (
	CtxGormDbKey = iota
	CtxGormScopesKey
)
