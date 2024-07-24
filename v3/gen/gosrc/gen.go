package gosrc

import (
	"github.com/bychannel/tabtoy/v3/gen"
	"github.com/bychannel/tabtoy/v3/model"
	"github.com/bychannel/tabtoy/v3/report"
	"github.com/davyxu/protoplus/codegen"
)

func Generate(globals *model.Globals) (data []byte, err error) {

	cg := codegen.NewCodeGen("gosrc").
		RegisterTemplateFunc(codegen.UsefulFunc).
		RegisterTemplateFunc(gen.UsefulFunc).
		RegisterTemplateFunc(UsefulFunc)

	err = cg.ParseTemplate(templateText, globals).Error()
	if err != nil {
		return
	}

	err = cg.FormatGoCode().Error()
	if err != nil {
		report.Log.Infoln(cg.Code())
		return
	}

	err = cg.WriteBytes(&data).Error()

	return
}
