package mysql

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"reflect"
	"strconv"
	"strings"
	"time"

	"gopkg.in/svrkit.v1/strutil"
)

const (
	DefaultDBTagKey = "db"
	DBColumnName    = "name"
	DBColumnType    = "type"
	DBColumnIndex   = "index"
	DBColumnCollate = "collate"
	DBColumnDefault = "default"

	StmtCacheIndexSelect = 0
	StmtCacheIndexUpdate = 1
	StmtCacheIndexInsert = 2
	StmtCacheIndexCreate = 3
)

var (
	typeTime        = reflect.TypeOf(time.Time{})
	cachedSourceAst = make(map[string]*SourceAstCache)
)

type SourceAstCache struct {
	fast *ast.File
	fset *token.FileSet
}

// DBIndexMeta 索引元信息
type DBIndexMeta struct {
	Name       string   // 索引名称
	Fields     []string // 包含列
	Priority   []int    // 用于联合索引
	AutoIncr   bool     //
	IsHashType bool     // btree/hash
	OrderDesc  bool     // ASC/DESC
	Unique     bool     // 唯一索引
	Primary    bool     // 主键索引
}

func NewDBIndexMeta(fieldName string) *DBIndexMeta {
	return &DBIndexMeta{
		Fields: []string{fieldName},
	}
}

// SortedFieldIndex 按优先级排序后的字段
func (n *DBIndexMeta) SortedFieldIndex() []int {
	switch len(n.Fields) {
	case 0:
		return nil
	case 1:
		return []int{0}
	}
	var index = make([]int, len(n.Fields))
	for i := 0; i < len(n.Fields); i++ {
		index[i] = i
	}
	for i := 0; i < len(index); i++ {
		for j := i; j > 0 && n.Priority[j] < n.Priority[j-1]; j-- {
			index[j], index[j-1] = index[j-1], index[j]
		}
	}
	return index
}

// table元信息
type DBTableMeta struct {
	Name             string              // 名称
	Comment          string              // 注释
	FieldNames       []string            // 排序的字段名称
	FiledTypeMapping map[string]string   // 字段类型
	FiledNameMapping map[string]string   // 字段名称
	FieldComments    map[string]string   // 字符注释
	FieldExtAttr     map[string][]string // 字段附加属性
	TableAttrKey     []string            // table附加属性
	TableAttrValue   []string            //
	IndexList        []*DBIndexMeta      // 索引
	CachedStmts      []string            //
}

func (m *DBTableMeta) AddFieldAttr(name, s string) {
	m.FieldExtAttr[name] = append(m.FieldExtAttr[name], s)
}

func (m *DBTableMeta) AddTableAttr(k, v string) {
	for i, key := range m.TableAttrKey {
		if key == k {
			m.TableAttrValue[i] = v
			return
		}
	}
	m.TableAttrKey = append(m.TableAttrKey, k)
	m.TableAttrValue = append(m.TableAttrValue, v)
}

func NewDBTableMeta(name string) *DBTableMeta {
	meta := &DBTableMeta{
		Name:             name,
		FiledTypeMapping: make(map[string]string),
		FiledNameMapping: make(map[string]string),
		FieldComments:    make(map[string]string),
		FieldExtAttr:     make(map[string][]string),
		CachedStmts:      make([]string, 4),
	}
	meta.AddTableAttr("ENGINE", "InnoDB")
	meta.AddTableAttr("DEFAULT CHARSET", "utf8mb4")
	return meta
}

// SQLStmtGen SQL statement generator
type SQLStmtGen struct {
	disableNotNull bool
	tableMetas     map[string]*DBTableMeta
}

func NewSQLStmtGen() *SQLStmtGen {
	return &SQLStmtGen{
		tableMetas: make(map[string]*DBTableMeta),
	}
}

func (g *SQLStmtGen) GetMetaList() []*DBTableMeta {
	var list = make([]*DBTableMeta, 0, len(g.tableMetas))
	for _, v := range g.tableMetas {
		list = append(list, v)
	}
	return list
}

// RegisterStruct 注册结构
//
//	 定义格式
//		type MyORM struct {
//		    	ID   int64 	`db:"name=id type=bigint index=primary,auto_incr"`
//				Name string `db:"name=name type=varchar(50) collate=utf8mb4"`
//				Hash string `db:"type=char(40) collate=utf8mb4"`
//		}
func (g *SQLStmtGen) RegisterStruct(tag string, ptr interface{}, srcFile string) {
	if tag == "" {
		tag = DefaultDBTagKey
	}
	var st = reflect.TypeOf(ptr).Elem()
	if st.NumField() <= 0 {
		log.Panicf("parseStruct: %s has no fields", st.Name())
	}
	meta := NewDBTableMeta(st.Name())
	meta.FieldNames = make([]string, 0, st.NumField())
	if srcFile != "" {
		if err := parseComments(meta, srcFile); err != nil {
			log.Printf("parseComments: %s, %v", st.Name(), err)
		}
	}
	for i := 0; i < st.NumField(); i++ {
		var field = st.Field(i)
		var colName, colType string
		tagText := field.Tag.Get(tag)
		tagFields := strutil.ParseKVPairs(tagText, ' ', '=')
		if v := tagFields[DBColumnName]; len(v) > 0 {
			colName = v
		}
		if colName == "" {
			colName = strutil.ToSnakeCase(field.Name)
		}
		if v := tagFields[DBColumnType]; v != "" {
			colType = v
		}
		if colType == "" {
			colType = toMySQLType(field.Type)
		}
		if v := tagFields[DBColumnIndex]; v != "" {
			g.parseIndex(meta, field.Name, v)
		}
		if v := tagFields[DBColumnCollate]; v != "" {
			meta.AddFieldAttr(field.Name, fmt.Sprintf("COLLATE '%s'", v))
		}
		if v := tagFields[DBColumnDefault]; v != "" {
			meta.AddFieldAttr(field.Name, fmt.Sprintf("DEFAULT %s", v))
		}
		meta.FieldNames = append(meta.FieldNames, field.Name)
		meta.FiledTypeMapping[field.Name] = colType
		meta.FiledNameMapping[field.Name] = colName
	}
	g.cacheStmts(meta, st)
	g.tableMetas[meta.Name] = meta
}

// 读取源代码文件解析struct注释
func parseComments(meta *DBTableMeta, filename string) error {
	v, err := getSourceAst(filename)
	if err != nil {
		return err
	}
	file := v.fset.File(1)
	if file == nil {
		return fmt.Errorf("cannot locate file %s", filename)
	}
	for name, obj := range v.fast.Scope.Objects {
		if name != meta.Name {
			continue
		}
		if obj.Kind != ast.Typ {
			continue
		}
		typSpec, ok := obj.Decl.(*ast.TypeSpec)
		if !ok {
			continue
		}
		structType, ok := typSpec.Type.(*ast.StructType)
		if !ok {
			continue
		}
		comment := getStructComments(file, v.fast, structType)
		comment = strings.TrimPrefix(comment, "//")
		meta.Comment = strings.TrimSpace(comment)
		// 每个字段的注释
		for _, field := range structType.Fields.List {
			if len(field.Names) > 0 && field.Comment != nil && len(field.Comment.List) > 0 {
				name1 := field.Names[0].Name
				comment1 := field.Comment.List[0].Text
				comment1 = strings.TrimPrefix(comment1, "//")
				meta.FieldComments[name1] = strings.TrimSpace(comment1)
			}
		}
	}
	return nil
}

func getSourceAst(filename string) (*SourceAstCache, error) {
	if v, found := cachedSourceAst[filename]; found {
		return v, nil
	}
	src, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	fset := token.NewFileSet()
	fast, err := parser.ParseFile(fset, filename, src, parser.ParseComments)
	if err != nil {
		return nil, err
	}
	v := &SourceAstCache{fast, fset}
	cachedSourceAst[filename] = v
	return v, nil
}

// struct注释
func getStructComments(file *token.File, fast *ast.File, structTyp *ast.StructType) string {
	line := file.Line(structTyp.Struct)
	for _, commentGroup := range fast.Comments {
		if len(commentGroup.List) == 1 {
			comment := commentGroup.List[0]
			ln := file.Line(comment.Slash)
			if ln+1 == line {
				return comment.Text
			}
		}
	}
	return ""
}

// 解析字段索引
//
//	type MyORM struct {
//	    	ID   int64 	`db:"index=primary"`
//			Name string `db:"index=unique,order:asc"`
//			Hash string `db:"index=name:idx_hash,type=hash"`
//			Foo  int    `db:"index=name:idx_foo,priority=2"`
//			Bar  int    `db:"index=name:idx_foo,priority=1"`
//	}
func (g *SQLStmtGen) parseIndex(meta *DBTableMeta, name, s string) {
	var index = NewDBIndexMeta(name)
	var opts = strutil.ParseKVPairs(s, ',', ':')
	if v := opts["name"]; v != "" {
		index.Name = v
	} else {
		index.Name = fmt.Sprintf("idx_%s", strings.ToLower(name))
	}
	if _, found := opts["primary"]; found {
		index.Primary = true
		if _, found = opts["auto_incr"]; found {
			index.AutoIncr = true
			meta.AddFieldAttr(name, "AUTO_INCREMENT")
		}
	} else {
		if _, found := opts["unique"]; found {
			index.Unique = true
		}
	}
	if v := opts["order"]; "desc" == strings.ToLower(v) {
		index.OrderDesc = true
	}
	if v := opts["type"]; "hash" == strings.ToLower(v) {
		index.IsHashType = true
	}
	// 联合索引
	if s := opts["priority"]; s != "" {
		priority, _ := strconv.Atoi(s)
		for _, ind := range meta.IndexList {
			if ind.Name == index.Name {
				ind.Priority = append(ind.Priority, priority)
				ind.Fields = append(ind.Fields, index.Fields...)
				return
			}
		}
		index.Priority = append(index.Priority, priority)
	}
	meta.IndexList = append(meta.IndexList, index)
}

// 缓存SQL语句
func (g *SQLStmtGen) cacheStmts(meta *DBTableMeta, st reflect.Type) {
	g.cacheSelectStmt(meta, st)
	g.cacheInsertStmt(meta, st)
	g.cacheUpdateStmt(meta, st)
	g.cacheCreateStmt(meta, st)
}

func (g *SQLStmtGen) cacheSelectStmt(meta *DBTableMeta, st reflect.Type) {
	var buf strings.Builder
	buf.WriteString("SELECT ")
	for i := 0; i < st.NumField(); i++ {
		var field = st.Field(i)
		var name = meta.FiledNameMapping[field.Name]
		if i > 0 {
			buf.WriteByte(',')
		}
		fmt.Fprintf(&buf, "`%s`", name)
	}
	var tblname = strutil.ToSnakeCase(st.Name())
	fmt.Fprintf(&buf, " FROM `%s` ", tblname)
	meta.CachedStmts[StmtCacheIndexSelect] = buf.String()
}

func (g *SQLStmtGen) cacheInsertStmt(meta *DBTableMeta, st reflect.Type) {
	var buf strings.Builder
	var tblname = strutil.ToSnakeCase(st.Name())
	fmt.Fprintf(&buf, "INSERT INTO `%s`(", tblname)
	for i := 0; i < st.NumField(); i++ {
		var field = st.Field(i)
		var name = meta.FiledNameMapping[field.Name]
		if i > 0 {
			buf.WriteByte(',')
		}
		fmt.Fprintf(&buf, "`%s`", name)
	}
	fmt.Fprintf(&buf, ") VALUES(")
	for i := 0; i < st.NumField(); i++ {
		buf.WriteByte('?')
		if i+1 < st.NumField() {
			buf.WriteByte(',')
		}
	}
	buf.WriteByte(')')
	meta.CachedStmts[StmtCacheIndexInsert] = buf.String()
}

func (g *SQLStmtGen) cacheUpdateStmt(meta *DBTableMeta, st reflect.Type) {
	var buf strings.Builder
	var tblname = strutil.ToSnakeCase(st.Name())
	fmt.Fprintf(&buf, "UPDATE `%s` SET ", tblname)
	var cnt = 0
	for i := 0; i < st.NumField(); i++ {
		var field = st.Field(i)
		var name = meta.FiledNameMapping[field.Name]
		if cnt > 0 {
			buf.WriteByte(',')
		}
		cnt++
		fmt.Fprintf(&buf, "`%s`=?", name)
	}
	meta.CachedStmts[StmtCacheIndexUpdate] = buf.String()
}

// create index definition:
//
//	col_name column_definition | {INDEX | KEY} [index_name] [index_type] (key_part,...)
//
// index_type:
//
//	USING {BTREE | HASH}
//
// key_part:
//
//	{col_name [(length)] | (expr)} [ASC | DESC]
func (g *SQLStmtGen) genCreateIndex(meta *DBTableMeta, buf *strings.Builder) {
	for i, index := range meta.IndexList {
		if index.Primary {
			fmt.Fprintf(buf, " PRIMARY KEY ")
		} else if index.Unique {
			fmt.Fprintf(buf, " UNIQUE KEY ")
		} else {
			fmt.Fprintf(buf, " INDEX ")
		}
		if !index.Primary {
			fmt.Fprintf(buf, "`%s` ", index.Name)
		}
		buf.WriteByte('(')
		var sortedIdx = index.SortedFieldIndex()
		for i, idx := range sortedIdx {
			var field = index.Fields[idx]
			if i > 0 {
				buf.WriteByte(',')
			}
			name := meta.FiledNameMapping[field]
			fmt.Fprintf(buf, "`%s`", name)
			if index.OrderDesc {
				buf.WriteString(" DESC")
			}
		}
		buf.WriteByte(')')
		if index.IsHashType {
			buf.WriteString(" USING HASH")
		}
		if i+1 < len(meta.IndexList) {
			buf.WriteByte(',')
		}
		buf.WriteByte('\n')
	}
}

// create table语法 https://dev.mysql.com/doc/refman/5.7/en/create-table.html
func (g *SQLStmtGen) cacheCreateStmt(meta *DBTableMeta, st reflect.Type) {
	var buf strings.Builder
	var tblname = strutil.ToSnakeCase(st.Name())
	fmt.Fprintf(&buf, "CREATE TABLE IF NOT EXISTS `%s`(\n", tblname)
	for i := 0; i < st.NumField(); i++ {
		var field = st.Field(i)
		var name = meta.FiledNameMapping[field.Name]
		var sqlType = meta.FiledTypeMapping[field.Name]
		fmt.Fprintf(&buf, "  `%s` %s", name, sqlType)
		if !g.disableNotNull {
			fmt.Fprintf(&buf, " NOT NULL")
		}
		for _, v := range meta.FieldExtAttr[field.Name] {
			fmt.Fprintf(&buf, " %s", v)
		}
		if comment := meta.FieldComments[field.Name]; comment != "" {
			fmt.Fprintf(&buf, " COMMENT '%s'", comment)
		}
		if i > 0 || len(meta.IndexList) > 0 {
			buf.WriteByte(',')
		}
		buf.WriteByte('\n')
	}
	g.genCreateIndex(meta, &buf)
	buf.WriteByte(')')
	for i, k := range meta.TableAttrKey {
		fmt.Fprintf(&buf, " %s=%s", k, meta.TableAttrValue[i])
	}
	if meta.Comment != "" {
		fmt.Fprintf(&buf, " COMMENT='%s'", meta.Comment)
	}
	buf.WriteByte(';')
	meta.CachedStmts[StmtCacheIndexCreate] = buf.String()
}

func (g *SQLStmtGen) GetMeta(name string) *DBTableMeta {
	return g.tableMetas[name]
}

// 生成SELECT语句
func (g *SQLStmtGen) SelectStmtOf(ptr interface{}) string {
	var st = reflect.TypeOf(ptr).Elem()
	var meta = g.tableMetas[st.Name()]
	if meta == nil {
		log.Panicf("struct %s not registered", st.Name())
	}
	return meta.CachedStmts[StmtCacheIndexSelect]
}

// 生成INSERT语句
func (g *SQLStmtGen) InsertStmtOf(ptr interface{}) string {
	var st = reflect.TypeOf(ptr).Elem()
	var meta = g.tableMetas[st.Name()]
	if meta == nil {
		log.Panicf("struct %s not registered", st.Name())
	}
	return meta.CachedStmts[StmtCacheIndexInsert]
}

// 生成CREATE语句
func (g *SQLStmtGen) CreateStmtOf(ptr interface{}) string {
	var st = reflect.TypeOf(ptr).Elem()
	var meta = g.tableMetas[st.Name()]
	if meta == nil {
		log.Panicf("struct %s not registered", st.Name())
	}
	return meta.CachedStmts[StmtCacheIndexCreate]
}

// 生成UPDATE语句
func (g *SQLStmtGen) UpdateStmtOf(ptr interface{}) string {
	var st = reflect.TypeOf(ptr).Elem()
	var meta = g.tableMetas[st.Name()]
	if meta == nil {
		log.Panicf("struct %s not registered", st.Name())
	}
	return meta.CachedStmts[StmtCacheIndexUpdate]
}

func (g *SQLStmtGen) UpdateQueryOf(ptr interface{}, keys []string) (string, []interface{}) {
	var value = reflect.ValueOf(ptr).Elem()
	var st = value.Type()
	var meta = g.tableMetas[st.Name()]
	if meta == nil {
		log.Panicf("struct %s not registered", st.Name())
	}
	var buf strings.Builder
	var tblname = strutil.ToSnakeCase(st.Name())
	var args = make([]interface{}, 0, st.NumField())
	fmt.Fprintf(&buf, "UPDATE `%s` SET ", tblname)
	var cnt = 0
	for i := 0; i < st.NumField(); i++ {
		var field = st.Field(i)
		if findStringInArray(field.Name, keys) >= 0 {
			continue
		}
		if cnt > 0 {
			buf.WriteByte(',')
		}
		cnt++
		var name = meta.FiledNameMapping[field.Name]
		fmt.Fprintf(&buf, "`%s`=?", name)
		args = append(args, value.Field(i).Interface())
	}
	if len(keys) > 0 {
		buf.WriteString(" WHERE ")
	}
	for i, key := range keys {
		field, ok := st.FieldByName(key)
		if !ok {
			log.Panicf("field %s.%s not found", st.Name(), key)
		}
		v := value.FieldByName(key)
		if i > 0 {
			buf.WriteString(" AND ")
		}
		var name = meta.FiledNameMapping[field.Name]
		fmt.Fprintf(&buf, "`%s`=?", name)
		args = append(args, v.Interface())
	}
	return buf.String(), args
}

func findStringInArray(s string, array []string) int {
	for i, elem := range array {
		if elem == s {
			return i
		}
	}
	return -1
}

func toMySQLType(field reflect.Type) string {
	switch field.Kind() {
	case reflect.Bool:
		return "tinyint(1)"
	case reflect.Int8:
		return "tinyint"
	case reflect.Int16:
		return "smallint"
	case reflect.Int32:
		return "int"
	case reflect.Int64, reflect.Int:
		return "bigint"
	case reflect.Uint8:
		return "tinyint unsigned"
	case reflect.Uint16:
		return "smallint unsigned"
	case reflect.Uint32:
		return "int unsigned"
	case reflect.Uint64, reflect.Uint:
		return "bigint unsigned"
	case reflect.Float32:
		return "float"
	case reflect.Float64:
		return "double"
	case reflect.String:
		return "varchar(255)" // default 255 length
	default:
		if field == typeTime {
			return "datetime"
		}
	}
	return "text"
}

var _gen = NewSQLStmtGen()

func G() *SQLStmtGen {
	return _gen
}
