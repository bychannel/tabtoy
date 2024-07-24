package pbsrc

import (
	"github.com/bychannel/tabtoy/v3/gen"
	"github.com/bychannel/tabtoy/v3/model"
	"github.com/davyxu/protoplus/codegen"
)

func Generate(globals *model.Globals) (data []byte, err error) {

	cg := codegen.NewCodeGen("pbsrc").
		RegisterTemplateFunc(codegen.UsefulFunc).
		RegisterTemplateFunc(gen.UsefulFunc).
		RegisterTemplateFunc(UsefulFunc)

	err = cg.ParseTemplate(templateText, globals).Error()
	if err != nil {
		return
	}

	err = cg.WriteBytes(&data).Error()

	return
}
