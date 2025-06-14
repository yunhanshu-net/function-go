package syscall

import (
	"github.com/yunhanshu-net/pkg/constants/syscall"
)

type Request struct {
	CallbackType syscall.Type `json:"callback_type"`
	Data         interface{}  `json:"data"`
}

//func (s *Request) DecodeData(el interface{}) error {
//	marshal, err := json.Marshal(s.Data)
//	if err != nil {
//		return err
//	}
//	err = json.Unmarshal(marshal, &el)
//	if err != nil {
//		return err
//	}
//	return nil
//}

type Response struct {
	Data interface{} `json:"data"`
}
