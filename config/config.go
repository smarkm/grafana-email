package config

var Instance *Config = &Config{
	GrafanaUrl: "http://localhost:3001",
	OrgAPIKeys: make(map[string]string),
	EmailFrom:  "",
	SmtpPort:   587,
}

type Config struct {
	DebugModel    bool              `json:"debugModel"`
	GrafanaUrl    string            `json:"grafanaBaseUrl"`
	OrgAPIKeys    map[string]string `json:"orgAPIKeys"` //
	EmailFrom     string            `json:"emailFrom"`
	EmailUserName string            `json:"emailUserName"`
	EmailPasswd   string            `json:"emailPassword"`
	SmtpHost      string            `json:"smtpHost"`
	SmtpPort      int               `json:"smtpPort"`
}

func Init() *Config {
	Instance.OrgAPIKeys["1"] = "eyJrIjoiVnVEbGhodlh2bHllV3J6bkFCcjN5TEZiUkFhY01RNkkiLCJuIjoiYSIsImlkIjoxfQ=="
	Instance.EmailPasswd = "vktytcrfxcdojbha"
	Instance.EmailFrom = "1187650061@qq.com"
	Instance.EmailUserName = "1187650061@qq.com"
	Instance.SmtpHost = "smtp.qq.com"
	return Instance
}
