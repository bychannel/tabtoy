package javasrc

// 报错行号+3
const templateText = `// Generated by {{.ToolName}}
// DO NOT EDIT!!
// Version: {{.Version}}
package {{.PackageName}};
import java.util.ArrayList;
import java.util.HashMap;
import java.util.List;
import java.util.Map;

public class {{.CombineStructName}} {	{{range $sn, $objName := $.Types.EnumNames}}
	 public enum {{$objName}} { {{range $fi,$field := $.Types.AllFieldByName $objName}}
			{{$field.FieldName}}({{$field.Value}}), // {{$field.Name}} {{end}}
		;
		{{$objName}}(int v) {
		   this.{{$objName}} = v;
		}
		public int get{{$objName}}() {
		   return {{$objName}};
		}
		
		private final int {{$objName}};

		public static {{$objName}} fromInt( int v ){
			switch (v){ {{range $fi,$field := $.Types.AllFieldByName $objName}}
			  case {{$field.Value}}:
				return {{$field.FieldName}}; {{end}}
			  default:
				return null;
			}
		}

	 }
	{{end}}
	{{range $sn, $objName := $.Types.StructNames}}
	public class {{$objName}} { {{range $fi,$field := $.Types.AllFieldByName $objName}}	
		public {{JavaType $field false}} {{$field.FieldName}} = {{JavaDefaultValue $ $field}}; // {{$field.Name}}; {{end}}
	}
	{{end}}
	{{range $ti, $tab := $.Datas.AllTables}}
	public List<{{$tab.HeaderType}}> {{$tab.HeaderType}} = new ArrayList<>(); // table: {{$tab.HeaderType}} {{end}}

	// Indices {{range $ii, $idx := GetIndices $}}
	public Map<{{JavaType $idx.FieldInfo true}}, {{$idx.Table.HeaderType}}> {{$idx.Table.HeaderType}}By{{$idx.FieldInfo.FieldName}} = new HashMap<{{JavaType $idx.FieldInfo true}}, {{$idx.Table.HeaderType}}>(); // table: {{$idx.Table.HeaderType}} {{end}}
	{{if HasKeyValueTypes $}}
	//{{range $ti, $name := GetKeyValueTypeNames $}} table: {{$name}}
	public {{$name}} GetKeyValue_{{$name}}() {
		return {{$name}}.get(0);
	}
	{{end}}{{end}}
	public interface TableEvent{
		void OnPreProcess( );
		void OnPostProcess( );
	}
	// Handlers
	private List<TableEvent> eventHandlers = new ArrayList<TableEvent>();

	// 注册加载后回调
	public void RegisterEventHandler(TableEvent ev ){
		eventHandlers.add(ev);
	}

	// 清除索引和数据, 在处理前调用OnPostProcess, 可能抛出异常
	public void ResetData()  {

		for( TableEvent ev : eventHandlers){
			ev.OnPreProcess();
		}
		{{range $ti, $tab := $.Datas.AllTables}}
		{{$tab.HeaderType}}.clear(); {{end}}
		{{range $ii, $idx := GetIndices $}}
		{{$idx.Table.HeaderType}}By{{$idx.FieldInfo.FieldName}}.clear(); {{end}}	
	}

	// 构建索引, 需要捕获OnPostProcess可能抛出的异常
	public void  BuildData()  {
		{{range $ii, $idx := GetIndices $}}
		for( {{$idx.Table.HeaderType}} v:{{$idx.Table.HeaderType}} ) {
			{{$idx.Table.HeaderType}}By{{$idx.FieldInfo.FieldName}}.put(v.{{$idx.FieldInfo.FieldName}}, v);
		}{{end}}

		for( TableEvent ev : eventHandlers){
			ev.OnPostProcess();
		}
	}
}

`
