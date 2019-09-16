package views

import (
	"github.com/gobwas/ws"
	"github.com/golang/glog"
	"github.com/mkawserm/hypcore/package/constants"
	core2 "github.com/mkawserm/hypcore/package/core"
	"github.com/mkawserm/hypcore/package/cors"
	"github.com/mkawserm/hypcore/package/gqltypes"
	"github.com/mkawserm/hypcore/package/mcodes"
	"net/http"
)

type WebSocketUpgradeView struct {
	Context *core2.HContext
}

func (wsu *WebSocketUpgradeView) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if !cors.CheckCROSAndStepForward(wsu.Context.CORSOptions, w, r) {
		glog.Infoln("CORS!!!")
		return
	}

	if r.URL.Path == string(wsu.Context.WebSocketUpgradePath) {

		if r.Method != http.MethodGet {
			errorType := gqltypes.NewErrorType()
			errorType.Group = mcodes.GraphQLWSUpgradeGroupCode
			errorType.Code = mcodes.GraphQLWSUpgradeRequestMethodError
			errorType.MessageType = "GraphQLWSUpgradeException"
			errorType.AddStringMessage("Oops! WebSocket upgrade must be done using get request !!!")
			GraphQLSmartErrorMessage(w, errorType, 400)
			return
		}

		if r.ProtoMajor < 1 || (r.ProtoMajor == 1 && r.ProtoMinor < 1) {
			errorType := gqltypes.NewErrorType()
			errorType.Group = mcodes.GraphQLWSUpgradeGroupCode
			errorType.Code = mcodes.GraphQLWSUpgradeBadProtocol
			errorType.MessageType = "GraphQLWSUpgradeException"
			errorType.AddStringMessage("Oops! Bad protocol !!!")
			GraphQLSmartErrorMessage(w, errorType, 400)
			return
		}

		if r.Host == "" {
			errorType := gqltypes.NewErrorType()
			errorType.Group = mcodes.GraphQLWSUpgradeGroupCode
			errorType.Code = mcodes.GraphQLWSUpgradeHostNotFound
			errorType.MessageType = "GraphQLWSUpgradeException"
			errorType.AddStringMessage("Oops! No Host found !!!")
			GraphQLSmartErrorMessage(w, errorType, 400)
			return
		}

		if u := httpGetHeader(r.Header, constants.HeaderUpgradeCanonical); u != "websocket" && !core2.StrEqualFold(u, "websocket") {
			errorType := gqltypes.NewErrorType()
			errorType.Group = mcodes.GraphQLWSUpgradeGroupCode
			errorType.Code = mcodes.GraphQLWSUpgradeNoUpgradeHeaderFound
			errorType.MessageType = "GraphQLWSUpgradeException"
			errorType.AddStringMessage("Oops! No Upgrade header found !!!")
			GraphQLSmartErrorMessage(w, errorType, 400)
			return
		}

		if c := httpGetHeader(r.Header, constants.HeaderConnectionCanonical); c != "Upgrade" && !core2.StrHasToken(c, "upgrade") {
			errorType := gqltypes.NewErrorType()
			errorType.Group = mcodes.GraphQLWSUpgradeGroupCode
			errorType.Code = mcodes.GraphQLWSUpgradeNoConnectionHeaderFound
			errorType.MessageType = "GraphQLWSUpgradeException"
			errorType.AddStringMessage("Oops! No Connection header found !!!")
			GraphQLSmartErrorMessage(w, errorType, 400)

			return
		}

		if wsu.Context.HasAuthVerify() {
			uid := ""
			ok := false
			group := ""

			h := httpGetHeader(r.Header, constants.HeaderAuthorizationCanonical)

			if h == "" {
				h = r.URL.Query().Get("token")
				if h != "" {
					h = wsu.Context.AuthBearer + " " + h
				}
			}

			if h == "" {
				errorType := gqltypes.NewErrorType()
				errorType.Group = mcodes.GraphQLWSUpgradeGroupCode
				errorType.Code = mcodes.GraphQLWSUpgradeNoAuthorizationHeaderFound
				errorType.MessageType = "GraphQLWSUpgradeException"
				errorType.AddStringMessage("Oops! No Authorization header found !!!")
				GraphQLSmartErrorMessage(w, errorType, 400)
				return

			} else {
				var dataMap map[string]interface{}

				dataMap, ok = wsu.Context.AuthVerify.Verify([]byte(h),
					[]byte(wsu.Context.AuthPublicKey),
					[]byte(wsu.Context.AuthSecretKey),
					wsu.Context.AuthBearer)

				if ok {
					if uniqueId, found := dataMap["uid"]; found {
						uid = uniqueId.(string)
					}
					if groupString, found := dataMap["group"]; found {
						group = groupString.(string)
					}
				}
			}

			if ok {
				if uid == "" {
					errorType := gqltypes.NewErrorType()
					errorType.Group = mcodes.GraphQLWSUpgradeGroupCode
					errorType.Code = mcodes.GraphQLWSUpgradeNoUIDFoundFromToken
					errorType.MessageType = "GraphQLWSUpgradeException"
					errorType.AddStringMessage("Oops! No UID found from AuthVerifyInterface !!!")
					GraphQLSmartErrorMessage(w, errorType, 400)
					return
				}

				conn, _, _, err := ws.UpgradeHTTP(r, w)

				if err != nil {
					// NOTE: Failed to do upgrade handshake
					return
				}

				if err := wsu.Context.AddConnection(conn); err != nil {
					glog.Errorf("Failed to add connection %v\n", err)
					_ = conn.Close()
				} else {
					//NOTE Only if has websocket auth
					//connection added to the container now we'll map it to specific user based on authorization
					wsu.Context.AddUser(uid, group, core2.WebsocketFileDescriptor(conn))
				}

			} else { // Failed to authorize. not ok
				errorType := gqltypes.NewErrorType()
				errorType.Group = mcodes.GraphQLWSUpgradeGroupCode
				errorType.Code = mcodes.GraphQLWSUpgradeInvalidAuthorizationData
				errorType.MessageType = "GraphQLWSUpgradeException"
				errorType.AddStringMessage("Oops! Invalid Authorization data !!!")
				GraphQLSmartErrorMessage(w, errorType, 400)
			}

		} else { // No AuthVerifyInterface. So just upgrade connection

			conn, _, _, err := ws.UpgradeHTTP(r, w)

			if err != nil {
				// NOTE: Failed to do upgrade handshake
				return
			}

			if err := wsu.Context.AddConnection(conn); err != nil {
				glog.Errorf("Failed to add connection %v\n", err)
				_ = conn.Close()
			}

		} // End of No AuthVerifyInterface

	} else { // r.URL.Path != "/ws"
		errorType := gqltypes.NewErrorType()
		errorType.Group = mcodes.GraphQLWSUpgradeGroupCode
		errorType.Code = mcodes.GraphQLWSUpgradePathNotFound
		errorType.MessageType = "GraphQLWSUpgradeException"
		errorType.AddStringMessage("Oops!!!")
		GraphQLSmartErrorMessage(w, errorType, 404)
	}
}
