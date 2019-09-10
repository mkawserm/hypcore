package app

import (
	"fmt"
	"github.com/logrusorgru/aurora"
	"github.com/mkawserm/hypcore/package/models"
	"github.com/spf13/viper"
)

func AddSuperUser(v *viper.Viper, userName string, password string) {
	hc := PrepareServer(v)
	if hc != nil {
		hc.Setup()
		hcc := hc.GetContext()
		added := hcc.AddObject(models.NewSuperUser(userName, password))

		if added {
			fmt.Printf(aurora.BrightGreen("Super user <%s> added.\n").String(), userName)
		} else {
			fmt.Printf(aurora.BrightRed("Failed to add <%s> user.\n").String(), userName)
		}
	}
}
