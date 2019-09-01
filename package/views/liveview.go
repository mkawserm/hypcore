package views

import (
	core2 "github.com/mkawserm/hypcore/package/core"
	"net/http"
)

type LiveView struct {
	Context *core2.HContext
}

func (lView *LiveView) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//log.Printf("PATH: " + r.URL.Path)
	//log.Printf("Processing Middleware in the Live View.")
	//for _, mi := range l.Context.MiddlewareList {
	//	next := mi.ServeHTTP(l.Context, r, w)
	//	if next == false {
	//		return
	//	}
	//}
	//log.Printf("Middleware processing complete.")

	var output []byte
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if lView.Context.GetIsLive() == true {
		output = []byte("{\"isLive\":true,\"success_code\":\"200\"}")
	} else {
		output = []byte("{\"isLive\":false,\"success_code\":\"200\"}")
	}

	_, _ = w.Write(output)
}
