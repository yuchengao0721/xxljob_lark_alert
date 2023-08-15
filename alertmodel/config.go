package alertmodel

type Config struct {
	Interval    int    `mapstructure:"interval"`
	LarkToken   string `mapstructure:"lark_tokens"`
	Application struct {
		Name string `mapstructure:"name"`
		Port int    `mapstructure:"port"`
	} `mapstructure:"application"`
	Instances []struct {
		ID       string            `mapstructure:"id"`
		Address  string            `mapstructure:"address"`
		DB       string            `mapstructure:"db"`
		Username string            `mapstructure:"username"`
		Password string            `mapstructure:"password"`
		Labels   map[string]string `mapstructure:"labels"`
	} `mapstructure:"instances"`
}
