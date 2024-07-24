package gen

import "github.com/bychannel/tabtoy/v3/model"

type GenSingleFile func(globals *model.Globals) (data []byte, err error)

type GenCustom func(globals *model.Globals, param string) (err error)
