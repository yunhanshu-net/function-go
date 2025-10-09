package runner

type MsgServices interface {
}

type IMsgServices struct {
}

func (c *Context) GetMessageService() MsgServices {

	return &IMsgServices{}
}
