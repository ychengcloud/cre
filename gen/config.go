package gen

const (
	DefaultIDName = "id"
)

const (
	LoaderMysql    = "mysql"
	LoaderPostgres = "postgres"
)

const (
	TplModeSingle = "single"
	TplModeMulti  = "multi"
)

type Type string

const (
	Bool    Type = "bool"
	Binary       = "binary"
	Bit          = "bit"
	Int8         = "int8"
	Uint8        = "uint8"
	Int16        = "int16"
	Uint16       = "uint16"
	Int32        = "int32"
	Uint32       = "uint32"
	Int64        = "int64"
	Uint64       = "uint64"
	Float32      = "float32"
	Float64      = "float64"
	String       = "string"
	Time         = "time"
	Enum         = "enum"
	UUID         = "uuid"
	JSON         = "json"
)

type Config struct {
	Project   string         `yaml:"project"` // the name of the project.
	Package   string         `yaml:"package"` // the name of the generated package.
	Header    string         `yaml:"header"`
	Dialect   string         `yaml:"dialect"` // the name of the dialect.
	DSN       string         `yaml:"dsn"`
	Overwrite bool           `yaml:"overwrite"`
	Delim     Delim          `yaml:"delim"`   // 模板变量标识符
	Root      string         `yaml:"root"`    // 模板根目录
	GenRoot   string         `yaml:"genRoot"` // 生成根目录
	Attrs     map[string]any `yaml:"attrs"`   // 其他配置项

	// Templates 所有的 Template Path 需要保证唯一，实际模板文件路径仅为更好的组织文件
	Templates []*Template `yaml:"templates"`

	Tables []*Table `yaml:"tables"`
}

type Delim struct {
	Left  string `yaml:"left"`
	Right string `yaml:"right"`
}

type Template struct {
	Path    string `yaml:"path"`    // Path 模板相对路径,相对于Root
	GenPath string `yaml:"genPath"` // GenPath 生成路径，相对于GenRoot
	Format  string `yaml:"format"`  // Format 生成文件名格式
	// Mode 生成模式, 可选值: "single", "multi"
	// 默认: single 模式下, 所有表数据生成一个文件
	// multi 模式下, 每个表数据生成一个文件
	Mode string `yaml:"mode"`
}

type Table struct {
	Name   string   `yaml:"name"`
	Skip   bool     `yaml:"skip"` // Skip 忽略表
	Fields []*Field `yaml:"fields"`
}

type Field struct {
	Name       string    `yaml:"name"`
	Type       Type      `yaml:"type"` // 指定字段类型，优先级高于数据库定义
	Nullable   bool      `yaml:"nullable"`
	Optional   bool      `yaml:"optional"`
	Comment    string    `yaml:"comment"`
	Alias      string    `yaml:"alias"`
	Skip       bool      `yaml:"skip"` // Skip 忽略表
	Sortable   bool      `yaml:"sortable"`
	Filterable bool      `yaml:"filterable"`
	Operations []string  `yaml:"operations"`
	Remote     bool      `yaml:"remote"`
	Relation   *Relation `yaml:"relation"`
}

// Relation represents	a Relation definition.
type Relation struct {
	Name      string     `yaml:"name"`
	Type      string     `yaml:"type"`
	Field     string     `yaml:"field"`
	RefTable  string     `yaml:"ref_table"`
	RefField  string     `yaml:"ref_field"`
	JoinTable *JoinTable `yaml:"join_table"` // 当 Type 为 ManyToMany 时, JoinTable 不为空
	Inverse   bool       `yaml:"inverse"`
}

type JoinTable struct {
	Name     string `yaml:"name"`
	Table    string `yaml:"table"`
	RefTable string `yaml:"ref_table"`
	Field    string `yaml:"field"`
	RefField string `yaml:"ref_field"`
}
