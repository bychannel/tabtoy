package printer

import (
	"fmt"

	"github.com/bychannel/tabtoy/util"
	"github.com/bychannel/tabtoy/v2/i18n"
	"github.com/bychannel/tabtoy/v2/model"
)

func valueWrapperLua(g *Globals, t model.FieldType, n *model.Node) string {

	switch t {
	case model.FieldType_String:
		return util.StringWrap(util.StringEscape(n.Value))
	case model.FieldType_Enum:
		if g.LuaEnumIntValue {
			return fmt.Sprintf("%d", n.EnumValue)
		} else {
			return fmt.Sprintf("\"%s\"", n.Value)
		}

	}

	return n.Value
}

type luaPrinter struct {
}

func (self *luaPrinter) Run(g *Globals) *Stream {

	stream := NewStream()

	stream.Printf("-- Generated by github.com/davyxu/tabtoy\n")
	stream.Printf("-- Version: %s\n", g.Version)

	if g.LuaTabHeader != "" {
		stream.Printf("\n%s\n", g.LuaTabHeader)
	}

	stream.Printf("\nlocal tab = {\n")

	for tabIndex, tab := range g.Tables {

		if !tab.LocalFD.MatchTag(".lua") {
			log.Infof("%s: %s", i18n.String(i18n.Printer_IgnoredByOutputTag), tab.Name())
			continue
		}

		if !printTableLua(g, stream, tab) {
			return nil
		}

		// 根字段分割
		if tabIndex < len(g.Tables)-1 {
			stream.Printf(", ")
		}

		stream.Printf("\n\n")
	}

	// local tab = {
	stream.Printf("}\n\n")

	if !genLuaIndexCode(stream, g.CombineStruct) {
		return stream
	}

	// 生成枚举
	if !genLuaEnumCode(g, stream, g.FileDescriptor) {
		return stream
	}

	stream.Printf("\nreturn tab")

	return stream
}

func printTableLua(g *Globals, stream *Stream, tab *model.Table) bool {

	stream.Printf("	%s = {\n", tab.LocalFD.Name)

	// 遍历每一行
	for rIndex, r := range tab.Recs {

		// 每一行开始
		stream.Printf("		{ ")

		// 遍历每一列
		for rootFieldIndex, node := range r.Nodes {

			if node.IsRepeated {
				stream.Printf("%s = { ", node.Name)
			} else {
				stream.Printf("%s = ", node.Name)
			}

			// 普通值
			if node.Type != model.FieldType_Struct {

				if node.IsRepeated {

					// repeated 值序列
					for arrIndex, valueNode := range node.Child {

						stream.Printf("%s", valueWrapperLua(g, node.Type, valueNode))

						// 多个值分割
						if arrIndex < len(node.Child)-1 {
							stream.Printf(", ")
						}

					}
				} else {
					// 单值
					valueNode := node.Child[0]

					stream.Printf("%s", valueWrapperLua(g, node.Type, valueNode))

				}

			} else {

				// 遍历repeated的结构体
				for structIndex, structNode := range node.Child {

					// 结构体开始
					stream.Printf("{ ")

					// 遍历一个结构体的字段
					for structFieldIndex, fieldNode := range structNode.Child {

						// 值节点总是在第一个
						valueNode := fieldNode.Child[0]

						stream.Printf("%s= %s", fieldNode.Name, valueWrapperLua(g, fieldNode.Type, valueNode))

						// 结构体字段分割
						if structFieldIndex < len(structNode.Child)-1 {
							stream.Printf(", ")
						}

					}

					// 结构体结束
					stream.Printf(" }")

					// 多个结构体分割
					if structIndex < len(node.Child)-1 {
						stream.Printf(", ")
					}

				}

			}

			if node.IsRepeated {
				stream.Printf(" }")
			}

			// 根字段分割
			if rootFieldIndex < len(r.Nodes)-1 {
				stream.Printf(", ")
			}

		}

		// 每一行结束
		stream.Printf(" 	}")

		if rIndex < len(tab.Recs)-1 {
			stream.Printf(",")
		}

		stream.Printf("\n")

	}

	// Sample = {
	stream.Printf("	}")

	return true

}

func anyFieldOutput(d *model.Descriptor) bool {
	for _, fd := range d.Fields {

		if fd.Meta.GetBool("LuaValueMapperString") {
			return true
		}

		if fd.Meta.GetBool("LuaStringMapperValue") {
			return true
		}

	}

	return false
}

// 收集需要构建的索引的类型
func genLuaEnumCode(g *Globals, stream *Stream, globalFile *model.FileDescriptor) bool {

	stream.Printf("\ntab.Enum = {\n")

	// 遍历字段
	for _, d := range globalFile.Descriptors {

		if d.Kind != model.DescriptorKind_Enum {
			continue
		}

		if anyFieldOutput(d) {
			stream.Printf("	%s = {\n", d.Name)

			for _, fd := range d.Fields {

				if fd.Meta.GetBool("LuaValueMapperString") {
					stream.Printf("		[%d] = \"%s\",\n", fd.EnumValue, fd.Name)
				}

			}

			for _, fd := range d.Fields {

				if fd.Meta.GetBool("LuaStringMapperValue") {
					stream.Printf("		%s = %d,\n", fd.Name, fd.EnumValue)
				}

			}

			stream.Printf("	},\n")
		}

	}

	stream.Printf("}\n")

	return true

}

// 收集需要构建的索引的类型
func genLuaIndexCode(stream *Stream, combineStruct *model.Descriptor) bool {

	// 遍历字段
	for _, fd := range combineStruct.Fields {

		// 这个字段被限制输出
		if fd.Complex != nil && !fd.Complex.File.MatchTag(".lua") {
			continue
		}

		// 对CombineStruct的XXDefine对应的字段
		if combineStruct.Usage == model.DescriptorUsage_CombineStruct {

			// 这个结构有索引才创建
			if fd.Complex != nil && len(fd.Complex.Indexes) > 0 {

				// 索引字段
				for _, key := range fd.Complex.Indexes {
					mapperVarName := fmt.Sprintf("tab.%sBy%s", fd.Name, key.Name)

					stream.Printf("\n-- %s\n", key.Name)
					stream.Printf("%s = {}\n", mapperVarName)
					stream.Printf("for _, rec in pairs(tab.%s) do\n", fd.Name)
					stream.Printf("\t%s[rec.%s] = rec\n", mapperVarName, key.Name)
					stream.Printf("end\n")
				}

			}

		}

	}

	return true

}

func init() {

	RegisterPrinter("lua", &luaPrinter{})

}
