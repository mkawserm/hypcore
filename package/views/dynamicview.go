package views

import (
	"github.com/golang/glog"
	"github.com/mkawserm/hypcore/package/constants"
	core2 "github.com/mkawserm/hypcore/package/core"
	"github.com/mkawserm/hypcore/package/cors"
	"github.com/mkawserm/hypcore/package/gqltypes"
	"github.com/mkawserm/hypcore/package/mcodes"
	"net/http"
	"regexp"
)

type DynamicView struct {
	Context *core2.HContext
}

func (dView *DynamicView) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	glog.Infof("PATH: " + r.URL.Path + "\n")

	if !cors.CheckCROSAndStepForward(dView.Context.CORSOptions, w, r) {
		glog.Infof("CORS!!!\n")
		return
	}

	glog.Infof("Processing Middleware in the DynamicView \n")
	for _, mi := range dView.Context.MiddlewareList {
		next := mi.ServeHTTP(dView.Context, r, w)
		if next == false {
			return
		}
	}
	glog.Infof("Middleware processing complete\n")

	// check for auth
	if dView.Context.HasAuthVerify() {
		uid := ""
		ok := false

		h := httpGetHeader(r.Header, constants.HeaderAuthorizationCanonical)
		if h == "" {
			errorType := gqltypes.NewErrorType()
			errorType.Group = mcodes.DynamicViewGroupCode
			errorType.Code = mcodes.DynamicViewNoAuthorizationHeaderFound
			errorType.MessageType = "DynamicViewException"
			errorType.AddStringMessage("Oops! No Authorization header found !!!")
			GraphQLSmartErrorMessage(w, errorType, 400)

			return
		} else {
			var dataMap map[string]interface{}

			dataMap, ok = dView.Context.AuthVerify.Verify([]byte(h),
				[]byte(dView.Context.AuthPublicKey),
				[]byte(dView.Context.AuthSecretKey),
				dView.Context.AuthBearer)

			if ok {
				if uniqueId, found := dataMap["uid"]; found {
					uid = uniqueId.(string)
				}
			}
		}

		if ok {
			if uid == "" {
				errorType := gqltypes.NewErrorType()
				errorType.Group = mcodes.DynamicViewGroupCode
				errorType.Code = mcodes.DynamicViewNoUIDFoundFromToken
				errorType.MessageType = "DynamicViewException"
				errorType.AddStringMessage("Oops! No UID found from AuthVerifyInterface !!!")
				GraphQLSmartErrorMessage(w, errorType, 400)
				return
			}

		} else { // Failed to authorize. not ok
			errorType := gqltypes.NewErrorType()
			errorType.Group = mcodes.DynamicViewGroupCode
			errorType.Code = mcodes.DynamicViewInvalidAuthorizationData
			errorType.MessageType = "DynamicViewException"
			errorType.AddStringMessage("Oops! Invalid Authorization data !!!")
			GraphQLSmartErrorMessage(w, errorType, 400)
			return
		}
	}

	glog.Infof("Dynamic route dispatch started\n")
	for _, route := range dView.Context.RouteList {
		var rc = regexp.MustCompile(route.Pattern)
		if rc.Match([]byte(r.URL.Path)) {
			route.HttpHandlerObject.ServeHTTP(dView.Context, r, w)
			return
		}
	}

	glog.Infof("Dynamic route dispatch failed to find a route\n")
	glog.Errorf("Showing 404 Http Error\n")

	errorType := gqltypes.NewErrorType()
	errorType.Group = mcodes.DynamicViewGroupCode
	errorType.Code = mcodes.DynamicViewPathNotFound
	errorType.MessageType = "DynamicViewException"
	errorType.AddStringMessage("Oops!!!")
	GraphQLSmartErrorMessage(w, errorType, 404)
}
