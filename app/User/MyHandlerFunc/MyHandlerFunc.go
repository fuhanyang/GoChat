package MyHandlerFunc

type MyHandlerFunc func(c *MyContext)

type stack []int

type MyContext struct {
	KV       map[string]interface{}
	Handlers []MyHandlerFunc
	stack    stack
}

func (c *MyContext) Push(v int) {
	c.stack = append(c.stack, v)
}
func (c *MyContext) Pop() int {
	l := c.stack
	if len(l) == 0 {
		return 0
	}
	v := l[len(l)-1]
	l = l[:len(l)-1]
	return v
}
func (l stack) len() int {
	return len(l)
}
func LoadHandlers(handlers ...MyHandlerFunc) *MyContext {
	var c MyContext
	for _, handler := range handlers {
		c.Handlers = append(c.Handlers, handler)
	}
	c.KV = make(map[string]interface{})
	c.stack = make(stack, 0, len(handlers))
	return &c
}
func (c *MyContext) Next() {
	for {
		index := c.Pop()
		c.Push(index + 1)
		c.Handlers[index](c)
		if index == len(c.Handlers)-1 {
			c.Pop()
			break
		}
		if c.stack.len() == 0 {
			break
		}
	}
}
func (c *MyContext) Abort() {
	c.stack = c.stack[:0]
}
func (c *MyContext) Apply() {
	if len(c.Handlers) == 0 {
		return
	}
	c.Next()
}
func CreateTask(ctx *MyContext) func() {
	var task func()
	task = func() {
		ctx.Apply()
	}
	return task
}
