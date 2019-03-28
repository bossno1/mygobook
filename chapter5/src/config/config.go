package config
import(
	"strings"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"github.com/lexkong/log"
)
type Config struct {
	Name string
}

func Init (cfg string) error {
	c := Config{
		Name : cfg,
	}

	//初始化配置文件 

	if err := c.initConfig(); err != nil {
		return err
	}
	//监控配置文件变化并热加载程序
	c.watchConfig()
 	// 初始化日志包
 	c.initLog()
	return nil
}
func (c *Config) initLog() {
    passLagerCfg := log.PassLagerCfg {
        Writers:        viper.GetString("log.writers"),
        LoggerLevel:    viper.GetString("log.logger_level"),
        LoggerFile:     viper.GetString("log.logger_file"),
        LogFormatText:  viper.GetBool("log.log_format_text"),
        RollingPolicy:  viper.GetString("log.rollingPolicy"),
        LogRotateDate:  viper.GetInt("log.log_rotate_date"),
        LogRotateSize:  viper.GetInt("log.log_rotate_size"),
        LogBackupCount: viper.GetInt("log.log_backup_count"),
    }

    log.InitWithConfig(&passLagerCfg)
}  
func (c *Config) initConfig() error {
	if c.Name != "" { 
		viper.SetConfigFile(c.Name) //如果指定了配置文件 ， 则解析
	} else {
		viper.AddConfigPath(".") //当前目录 
		viper.SetConfigName("config")
	}
	viper.SetConfigType("yaml")
	viper.AutomaticEnv() //读取匹配的环境变量
	viper.SetEnvPrefix("APISERVER") //
	replacer := strings.NewReplacer(".", "_")
	viper.SetEnvKeyReplacer(replacer)
	if err := viper.ReadInConfig();  err != nil{ 
		return err
	}
	return nil
}
//监控配置文件 并热加载 
func (c *Config) watchConfig(){
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event){
		log.Info("Config file changed:" + e.Name)
	})
}