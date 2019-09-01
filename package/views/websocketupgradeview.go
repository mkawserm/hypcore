package views

import (
	"github.com/gobwas/ws"
	"github.com/golang/glog"
	core2 "github.com/mkawserm/hypcore/package/core"
	"net/http"
)

type WebSocketUpgradeView struct {
	Context *core2.HContext
}

func (wsu *WebSocketUpgradeView) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path == string(wsu.Context.WebSocketUpgradePath) {

		if r.Method != http.MethodGet {
			httpBadRequest(w, []byte("Oops! WebSocket upgrade must be done using get request !!!"))
			return
		}

		if r.ProtoMajor < 1 || (r.ProtoMajor == 1 && r.ProtoMinor < 1) {
			httpBadRequest(w, []byte("Oops! Bad protocol !!!"))
			return
		}

		if r.Host == "" {
			httpBadRequest(w, []byte("Oops! No Host found !!!"))
			return
		}

		if u := httpGetHeader(r.Header, core2.HeaderUpgradeCanonical); u != "websocket" && !core2.StrEqualFold(u, "websocket") {
			httpBadRequest(w, []byte("Oops! No Upgrade header found !!!"))
			return
		}

		if c := httpGetHeader(r.Header, core2.HeaderConnectionCanonical); c != "Upgrade" && !core2.StrHasToken(c, "upgrade") {
			httpBadRequest(w, []byte("Oops! No Connection header found !!!"))
			return
		}

		if wsu.Context.HasAuth() {
			uid := ""
			ok := false

			h := httpGetHeader(r.Header, core2.HeaderAuthorizationCanonical)

			if h == "" {
				httpBadRequest(w, []byte("Oops! No Authorization header found !!!"))
				return

			} else {
				uid, ok = wsu.Context.Auth.GetUID([]byte(h), wsu.Context.AuthBearer)
			}

			if ok {
				if uid == "" {
					httpBadRequest(w, []byte("Oops! No UID found from AuthInterface !!!"))
					return
				}

				conn, _, _, err := ws.UpgradeHTTP(r, w)

				if err != nil {
					// NOTE: Failed to do upgrade handshake
					return
				}

				if err := wsu.Context.AddConnection(conn); err != nil {
					glog.Errorln("Failed to add connection %v", err)
					_ = conn.Close()
				} else {
					//NOTE Only if has websocket auth
					//connection added to the container now we'll map it to specific user based on authorization
					wsu.Context.AddUser(uid, core2.WebsocketFileDescriptor(conn))
				}

			} else { // Failed to authorize. not ok
				httpBadRequest(w, []byte("Oops! Invalid Authorization data !!!"))
			}

		} else { // No AuthInterface. So just upgrade connection

			conn, _, _, err := ws.UpgradeHTTP(r, w)

			if err != nil {
				// NOTE: Failed to do upgrade handshake
				return
			}

			if err := wsu.Context.AddConnection(conn); err != nil {
				glog.Errorln("Failed to add connection %v", err)
				_ = conn.Close()
			}

		} // End of No AuthInterface

	} else { // r.URL.Path != "/ws"
		httpNotFound(w, []byte("Oops!!!"))
	}
}
