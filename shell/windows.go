package shell

import (
	. "IxDShell/config"
	"fmt"
	"github.com/gen2brain/dlgs"
	"github.com/zserge/lorca"
	"log"
	"time"
)

func StartWindows() {
	localhost := fmt.Sprintf("%s/login?t=%d&client=true", CONF.WebAddr, time.Now().Unix())
	ui, err := lorca.New(localhost, "", 1100, 618, "--autoplay-policy=no-user-gesture-required")
	if err != nil {
		log.Fatal(err)
	}
	defer ui.Close()
	err = ui.Bind("openFileDialog", func() {
		//filename, ret, err := dlgs.File("Choose directory", "", true)
		filename, ret, err := dlgs.File("Select file", "", false)
		if err != nil {
			log.Println(err)
		}
		if filename != "" && ret {
			s := fmt.Sprintf(`externalInvokeOpen("%s")`, filename)
			_ = ui.Eval(s)
			if err != nil {
				log.Println(err)
			}
		}
	})
	if err != nil {
		log.Println(err)
		log.Fatal(err)
	}
	<-ui.Done()
}
