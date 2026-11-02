package rule

import (
	"fmt"
	"log/slog"
	"strconv"
	"sync"

	"github.com/dop251/goja"
)

const (
	exportedSrc = `const exportSrc = (typeof rules !== 'undefined') ? rules : [rule]`
)

type jsRuntime struct {
	vm              *goja.Runtime
	onRequestHeader []goja.Callable
}

type jsRule struct {
	program  *goja.Program
	exported *goja.Program
	pool     sync.Pool
}

func NewJsRule(name string, content string) (*jsRule, error) {
	prog, err := goja.Compile(name, content, true)
	if err != nil {
		return nil, fmt.Errorf("new js rule compile: %w", err)
	}

	exported, err := goja.Compile("export", exportedSrc, true)
	if err != nil {
		return nil, fmt.Errorf("exported js rule compile: %w", err)
	}

	js := &jsRule{
		program:  prog,
		exported: exported,
	}

	rt, err := js.newRuntime()
	if err != nil {
		return nil, err
	}

	js.pool.Put(rt)
	js.pool.New = func() any {
		rt, err := js.newRuntime()
		if err != nil { // already validated at startup; shouldn't happen
			slog.Error("err new pool", err)
			return nil
		}
		return rt
	}

	return js, nil
}

func (j *jsRule) newRuntime() (*jsRuntime, error) {
	vm := goja.New()

	if _, err := vm.RunProgram(j.program); err != nil {
		return nil, fmt.Errorf("run script: %w", err)
	}

	_, err := vm.RunProgram(j.exported)
	if err != nil {
		return nil, fmt.Errorf("script must define rules or rule: %w", err)
	}

	rt := &jsRuntime{
		vm: vm,
	}
	return rt, nil
}
