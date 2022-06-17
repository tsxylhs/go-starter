package code

type Context map[string]interface{}

func (ctx *Context) Get(key string) interface{} {
	return (*ctx)[key]
}

func (ctx *Context) MustGet(key string) interface{} {
	v := (*ctx)[key]

	if v == nil {
		panic("key " + key + " not present in context")
	}
	return v
}

func (ctx *Context) Set(key string, value interface{}) {
	(*ctx)[key] = value
}
