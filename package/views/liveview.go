package views

import (
	"github.com/golang/glog"
	core2 "github.com/mkawserm/hypcore/package/core"
	"github.com/mkawserm/hypcore/package/cors"
	"net/http"
)

type LiveView struct {
	Context *core2.HContext
}

func (lView *LiveView) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	glog.Infof("PATH: " + r.URL.Path + "\n")
	//log.Printf("PATH: " + r.URL.Path)
	//log.Printf("Processing Middleware in the Live View.")
	//for _, mi := range l.Context.MiddlewareList {
	//	next := mi.ServeHTTP(l.Context, r, w)
	//	if next == false {
	//		return
	//	}
	//}
	//log.Printf("Middleware processing complete.")
	if !cors.CheckCROSAndStepForward(lView.Context.CORSOptions, w, r) {
		glog.Infof("CORS!!!\n")
		return
	}

	glog.Infof("Processing Middleware in the LiveView\n")
	for _, mi := range lView.Context.MiddlewareList {
		next := mi.ServeHTTP(lView.Context, r, w)
		if next == false {
			return
		}
	}
	glog.Infof("Middleware processing complete\n")

	var output []byte
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if lView.Context.GetIsLive() == true {
		output = []byte("{\"isLive\":true,\"success_code\":\"HCHS200\"}")
	} else {
		output = []byte("{\"isLive\":false,\"success_code\":\"HCHS200\"}")
	}

	_, _ = w.Write(output)
}
