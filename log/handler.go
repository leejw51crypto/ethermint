package log

import (
	"fmt"

	"github.com/ethereum/go-ethereum/log"
	ethlog "github.com/ethereum/go-ethereum/log"
	tmlog "github.com/tendermint/tendermint/libs/log"
)

var ethermintLogger *tmlog.Logger = nil

func FuncHandler(fn func(r *ethlog.Record) error) log.Handler {
	return funcHandler(fn)
}

type funcHandler func(r *ethlog.Record) error

func (h funcHandler) Log(r *log.Record) error {
	return h(r)
}

func NewHandler(logger tmlog.Logger) ethlog.Handler {

	ethermintLogger = &logger

	return FuncHandler(func(r *log.Record) error {
		(*ethermintLogger).Debug(fmt.Sprintf("[EVM] %v", r))
		return nil
	})
}
