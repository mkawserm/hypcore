package app

import (
	"fmt"
	"github.com/logrusorgru/aurora"
	"github.com/mkawserm/hypcore/package/models"
	"github.com/spf13/viper"
	"time"
)

func AddSuperUser(v *viper.Viper, userName string, password string) {
	hc := PrepareServer(v)
	if hc != nil {
		hc.Setup()
		hcc := hc.GetContext()
		user := models.NewSuperUser(userName, password)
		updated := false
		if hcc.IsObjectExists(user) {
			hcc.GetObject(user)
			updated = true
		}

		user.UpdatedAt = time.Now().UnixNano()

		added := hcc.AddObject(user)

		if added {
			if updated {
				fmt.Printf(aurora.BrightGreen("Super user <%s> updated.\n").String(), userName)
			} else {
				fmt.Printf(aurora.BrightGreen("Super user <%s> added.\n").String(), userName)
			}

		} else {
			if updated {
				fmt.Printf(aurora.BrightRed("Failed to update <%s> user.\n").String(), userName)
			} else {
				fmt.Printf(aurora.BrightRed("Failed to add <%s> user.\n").String(), userName)
			}
		}
	}
}
