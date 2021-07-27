package util

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
)

var Config = viper.New()

func PrintWelcome() bool {
	b := !Config.IsSet("print.welcome") || Config.GetBool("print.welcome")
	if b {
		// http://patorjk.com/software/taag/
		log.Print(`
 ___   ___   ___  _____   __    _    
| |_) / / \ | |_)  | |   / /\  | |   
|_|   \_\_/ |_| \  |_|  /_/--\ |_|__                           
                                                                    
██╗  ██╗██╗  ██╗███╗   ██╗██╗ ██████╗    ██████╗ ██████╗ ███╗   ███╗
╚██╗██╔╝██║  ██║████╗  ██║██║██╔════╝   ██╔════╝██╔═══██╗████╗ ████║
 ╚███╔╝ ███████║██╔██╗ ██║██║██║        ██║     ██║   ██║██╔████╔██║
 ██╔██╗ ██╔══██║██║╚██╗██║██║██║        ██║     ██║   ██║██║╚██╔╝██║
██╔╝ ██╗██║  ██║██║ ╚████║██║╚██████╗██╗╚██████╗╚██████╔╝██║ ╚═╝ ██║
╚═╝  ╚═╝╚═╝  ╚═╝╚═╝  ╚═══╝╚═╝ ╚═════╝╚═╝ ╚═════╝ ╚═════╝ ╚═╝     ╚═╝
                                                                    `)
	}
	return b
}

func initConfig() {
	Config.SetConfigName("config")
	Config.AddConfigPath("./conf")

	if err := Config.ReadInConfig(); err != nil {
		panic(fmt.Errorf("Fatal error config file: %s ", err))
	}
}
