/**
 * Auth :   liubo
 * Date :   2018/11/1 16:26
 * Comment: 通用回调
 */

package single

import "fmt"

type IAction interface {
	Call()
	OnPanic() bool
}
func MakeCommonAction(callback func(...interface{}), args ...interface{}) IAction {
	ret := &commonActionCallback{}
	ret.Callback = callback
	ret.Args = append(make([]interface{}, 0), args...)
	return ret
}

type Action struct {
	callback IAction
	done     chan bool
	id       int
}
func (self *Action) safeExec() {
	defer func() {
		r := recover()
		if r != nil {
			self.safeRecover(r)
		}
	}()

	self.callback.Call()
}
func (self *Action) safeRecover(r interface{}) {

	defer func() {
		r := recover()
		if r != nil {
			fmt.Println("尝试恢复崩溃的时候，又崩溃了！")
		}
	}()

	if !self.callback.OnPanic() {
		fmt.Println(r)
	}
}


type commonActionCallback struct {
	Callback func(...interface{})
	Args     []interface{}
}
func (self *commonActionCallback) Call() {
	self.Callback(self.Args...)
}
func (self *commonActionCallback) OnPanic() bool {
	return false
}
