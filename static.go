package gostudy

import (
	"fmt"
	"net/http"
	"strings"
)

var Static map[string]string = make(map[string]string)

func serverStatic(w http.ResponseWriter, r *http.Request) bool {
	for prefix, static := range Static {
		fmt.Println(prefix, static, r.URL.Path)
		if strings.HasPrefix(r.URL.Path, prefix) {
			file := static + r.URL.Path[len(prefix):]
			fmt.Println("file:" + file)
			http.ServeFile(w, r, file)
			return true
		}
	}
	return false
}
