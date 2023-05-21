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
	Project   string         `yaml:"project" mapstructure:"project"` // the name of the project.
	Package   string         `yaml:"package" mapstructure:"package"` // the name of the generated package.
	Header    string         `yaml:"header" mapstructure:"header"`
	Dialect   string         `yaml:"dialect" mapstructure:"dialect"` // the name of the dialect.
	DSN       string         `yaml:"dsn" mapstructure:"dsn"`
	Overwrite bool           `yaml:"overwrite" mapstructure:"overwrite"`
	Delim     Delim          `yaml:"delim" mapstructure:"delim"`     // 模板变量标识符
	Root      string         `yaml:"root" mapstructure:"root"`       // 模板根目录
	GenRoot   string         `yaml:"genRoot" mapstructure:"genRoot"` // 生成根目录
	Attrs     map[string]any `yaml:"attrs" mapstructure:"attrs"`     // 其他配置项

	// Templates 所有的 Template Path 需要保证唯一，实际模板文件路径仅为更好的组织文件
	Templates []*Template `yaml:"templates" mapstructure:"templates"`

	Tables []*Table `yaml:"tables" mapstructure:"tables"`
}

type Delim struct {
	Left  string `yaml:"left" mapstructure:"left"`
	Right string `yaml:"right" mapstructure:"right"`
}

type Template struct {
	Path    string `yaml:"path" mapstructure:"path"`       // Path 模板相对路径,相对于Root
	GenPath string `yaml:"genPath" mapstructure:"genPath"` // GenPath 生成路径，相对于GenRoot
	Format  string `yaml:"format" mapstructure:"format"`   // Format 生成文件名格式
	// Mode 生成模式, 可选值: "single", "multi"
	// 默认: single 模式下, 所有表数据生成一个文件
	// multi 模式下, 每个表数据生成一个文件
	Mode string `yaml:"mode" mapstructure:"mode"`
}

type Table struct {
	Name   string   `yaml:"name" mapstructure:"name"`
	Skip   bool     `yaml:"skip" mapstructure:"skip"` // Skip 忽略表
	Fields []*Field `yaml:"fields" mapstructure:"fields"`
}

type Field struct {
	Name       string    `yaml:"name" mapstructure:"name"`
	Type       Type      `yaml:"type" mapstructure:"type"` // 指定字段类型，优先级高于数据库定义
	Nullable   bool      `yaml:"nullable" mapstructure:"nullable"`
	Optional   bool      `yaml:"optional" mapstructure:"optional"`
	Comment    string    `yaml:"comment" mapstructure:"comment"`
	Alias      string    `yaml:"alias" mapstructure:"alias"`
	Skip       bool      `yaml:"skip" mapstructure:"skip"` // Skip 忽略表
	Sortable   bool      `yaml:"sortable" mapstructure:"sortable"`
	Filterable bool      `yaml:"filterable" mapstructure:"filterable"`
	Operations []string  `yaml:"operations" mapstructure:"operations"`
	Remote     bool      `yaml:"remote" mapstructure:"remote"`
	Relation   *Relation `yaml:"relation" mapstructure:"relation"`
}

// Relation represents	a Relation definition.
type Relation struct {
	Name      string     `yaml:"name" mapstructure:"name"`
	Type      string     `yaml:"type" mapstructure:"type"`
	Field     string     `yaml:"field" mapstructure:"field"`
	RefTable  string     `yaml:"ref_table" mapstructure:"ref_table"`
	RefField  string     `yaml:"ref_field" mapstructure:"ref_field"`
	JoinTable *JoinTable `yaml:"join_table" mapstructure:"join_table"` // 当 Type 为 ManyToMany 时, JoinTable 不为空
	Inverse   bool       `yaml:"inverse" mapstructure:"inverse"`
}

type JoinTable struct {
	Name     string `yaml:"name" mapstructure:"name"`
	Table    string `yaml:"table" mapstructure:"table"`
	RefTable string `yaml:"ref_table" mapstructure:"ref_table"`
	Field    string `yaml:"field" mapstructure:"field"`
	RefField string `yaml:"ref_field" mapstructure:"ref_field"`
}
