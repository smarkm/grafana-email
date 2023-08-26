package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

var Instance *Config = &Config{
	GrafanaUrl: "http://localhost:3001",
	OrgAPIKeys: make(map[string]string),
	EmailFrom:  "",
	SmtpPort:   587,
	HTTPPort:   8001,
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
	PdfTitleAlign string            `json:"pdfTitleAlign"`
	PDFFontPath   string            `json:"fontPath"`
	DBUser        string            `json:"dbUser"`
	DBPassword    string            `json:"dbPassword"`
	DBHost        string            `json:"dbHost"`
	Database      string            `json:"database"`
	HTTPPort      int               `json:"httpPort"`
}

func Init(file string) *Config {
	filePath := "./config.json"
	if file != "" {
		filePath = file
	}

	// Read the file content
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal("Failed to read the file:", err)
		os.Exit(0)
	}

	// Print the content of the file as a string
	err = json.Unmarshal(content, &Instance)
	if err != nil {
		log.Println("Failed to parse config file: " + filePath)
		os.Exit(0)
	}
	if Instance.DebugModel {
		fmt.Println("Debug: " + string(content))
	}

	return Instance
}
