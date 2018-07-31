package gostudy

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

//https://blog.csdn.net/qibin0506/article/details/52614290

func NewServer(addr string, port int) *HttpServer {
	return &HttpServer{
		httpAddr: addr,
		httpPort: port,
		//mux:      &httpMux{router: make(map[string]map[string]reflect.Type)},
		mux: &httpMux{router: make(map[string]reflect.Type)},
	}
}

type httpMux struct {
	//router map[string]map[string]reflect.Type
	router map[string]reflect.Type
	spool  sync.Pool
}

type HttpServer struct {
	httpAddr string
	httpPort int
	mux      *httpMux
}

const (
	controllerSubfix = "Controller"
	actionSubfix     = "Action"
)

//add controller
func (this *HttpServer) Router(c interface{}) {
	reflectType := reflect.TypeOf(c)

	typeVal := reflectType.Elem()

	// myReflectVal := reflect.ValueOf(c)
	// myReflectType := reflect.Indirect(myReflectVal).Type()
	//获取controller name  只有reflectType才可以获取
	var controllerName string
	hasSubfix := strings.HasSuffix(typeVal.Name(), controllerSubfix)
	if hasSubfix {
		controllerName = strings.TrimSuffix(typeVal.Name(), controllerSubfix)
	} else {
		controllerName = strings.TrimSpace(typeVal.Name())
	}
	if _, ok := this.mux.router[controllerName]; ok {
		return
	} else {
		//this.mux.router[controllerName] = make(map[string]reflect.Type)
		this.mux.router[controllerName] = reflectType.Elem()
	}

	//this.mux.addAction(controllerName, reflectType)
	// this.mux.addAction(controllerName, reflectType, myReflectType)
}

func (this *HttpServer) Run() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	addr := this.httpAddr + ":" + strconv.Itoa(this.httpPort)
	fmt.Printf("addr is %s\n", addr)
	err := http.ListenAndServe(addr, this.mux)
	if err != nil {
		panic(err)
	}
}

// add method
func (this *httpMux) addAction(controllerName string, rt reflect.Type) {
	for i := 0; i < rt.NumMethod(); i++ {
		method := rt.Method(i).Name
		if strings.HasSuffix(method, actionSubfix) {
			action := strings.TrimSuffix(method, actionSubfix)
			fmt.Printf("action:%s", action)
			// this.router[controllerName][action] = rt.Elem()
		}
	}
}

func (this *httpMux) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if serverStatic(rw, r) {
		return
	}
	path := strings.TrimPrefix(r.URL.Path, "/")
	if path == "" || path == "/" {
		//              rw.WriteHeader(http.StatusForbidden)
		http.Error(rw, "forbidden", http.StatusForbidden)
		return
	}

	// ctx := this.spool.Get().(*Context)
	// defer this.spool.Put(ctx)
	ctx := &Context{}
	ctx.Config(rw, r)

	rPath := strings.Split(path, "/")
	cname := strings.Title(rPath[0])

	//controller 是否存在
	if controller, ok := this.router[cname]; ok {
		var actionName string
		if len(rPath) == 1 || rPath[1] == "" {
			actionName = "Index"
		} else {
			actionName = strings.Title(rPath[1])
		}

		vc := reflect.New(controller)

		methodName := actionName + actionSubfix

		sconstroller := vc.Interface().(IApp)
		sconstroller.Init(ctx)
		method := vc.MethodByName(methodName)
		method.Call(nil)
		// if controller, ok := this.router[cname][actionName]; ok {
		// 	vc := reflect.New(controller)

		// 	methodName := actionName + actionSubfix

		// 	r.ParseForm()

		// 	sconstroller := vc.Interface().(IApp)
		// 	sconstroller.Init(ctx)
		// 	method := vc.MethodByName(methodName)
		// 	method.Call(nil)
		// 	// fmt.Println([]reflect.Value{})
		// 	// method.Call([]reflect.Value{})
		// }
	} else {
		http.NotFound(rw, r)
		return
	}
}
