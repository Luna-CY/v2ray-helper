package configurator

var Configure struct {
	Home   string
	Server struct {
		Address string `mapstructure:"address"`
		Port    string `mapstructure:"port"`
		Release bool   `mapstructure:"release"`
	} `mapstructure:"server"`
	Acme struct {
		Email string `mapstructure:"email"`
	} `mapstructure:"acme"`
	Auth struct {
		Access  string `mapstructure:"access"`
		Manager string `mapstructure:"manager"`
	} `mapstructure:"auth"`
	Mail struct {
		Notice struct {
			Enable bool   `mapstructure:"enable"`
			To     string `mapstructure:"to"`
		} `mapstructure:"notice"`
		Smtp struct {
			Address  string `mapstructure:"address"`
			Port     string `mapstructure:"port"`
			User     string `mapstructure:"user"`
			Password string `mapstructure:"password"`
			Secret   string `mapstructure:"secret"`
		} `mapstructure:"smtp"`
	} `mapstructure:"mail"`
}
