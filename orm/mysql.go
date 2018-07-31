package orm

type baseMysql struct {
	dbBase
}

// func (this *baseMysql) Insert(md interface{}) {
// 	fmt.Printf("%v", md)
// 	val := reflect.ValueOf(md)
// 	ind := reflect.Indirect(val)
// 	typ := ind.Type()

// 	fmt.Printf("valueOf:%v\n Indirect:%v\n Type:%v\n type name %v", val, ind, typ, val.NumField())
// }

// create new mysql dbBaser.
func NewBaseMysql() IDbBaser {
	b := new(baseMysql)
	return b
}
