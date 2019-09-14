package views

import (
	"github.com/gobwas/ws"
	"github.com/golang/glog"
	core2 "github.com/mkawserm/hypcore/package/core"
	"github.com/mkawserm/hypcore/package/mcodes"
	"net/http"
)

type WebSocketUpgradeView struct {
	Context *core2.HContext
}

func (wsu *WebSocketUpgradeView) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path == string(wsu.Context.WebSocketUpgradePath) {

		if r.Method != http.MethodGet {
			GraphQLErrorMessage(w, []byte("Oops! WebSocket upgrade must be done using get request !!!"),
				mcodes.WebSocketUpgradeBadRequestMethod, 400)
			return
		}

		if r.ProtoMajor < 1 || (r.ProtoMajor == 1 && r.ProtoMinor < 1) {
			GraphQLErrorMessage(w, []byte("Oops! Bad protocol !!!"),
				mcodes.WebSocketBadProtocol, 400)
			return
		}

		if r.Host == "" {
			GraphQLErrorMessage(w, []byte("Oops! No Host found !!!"),
				mcodes.WebSocketNoHostFound, 400)
			return
		}

		if u := httpGetHeader(r.Header, core2.HeaderUpgradeCanonical); u != "websocket" && !core2.StrEqualFold(u, "websocket") {
			GraphQLErrorMessage(w, []byte("Oops! No Upgrade header found !!!"),
				mcodes.WebSocketNoUpgradeHeaderFound, 400)
			return
		}

		if c := httpGetHeader(r.Header, core2.HeaderConnectionCanonical); c != "Upgrade" && !core2.StrHasToken(c, "upgrade") {
			GraphQLErrorMessage(w, []byte("Oops! No Connection header found !!!"),
				mcodes.WebSocketNoConnectionHeaderFound, 400)
			return
		}

		if wsu.Context.HasAuth() {
			uid := ""
			ok := false
			group := ""

			h := httpGetHeader(r.Header, core2.HeaderAuthorizationCanonical)

			if h == "" {
				GraphQLErrorMessage(w, []byte("Oops! No Authorization header found !!!"),
					mcodes.NoAuthorizationHeaderFound, 400)
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
					GraphQLErrorMessage(w, []byte("Oops! No UID found from AuthVerifyInterface !!!"),
						mcodes.NoUIDFromAuthVerifyInterface, 400)
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
				GraphQLErrorMessage(w, []byte("Oops! Invalid Authorization data !!!"),
					mcodes.InvalidAuthorizationData, 400)
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
		GraphQLErrorMessage(w, []byte("Oops!!!"), mcodes.HttpNotFound, 404)
	}
}
