package model

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"time"

	"gopkg.in/gomail.v2"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"smark.freecoop.net/grafana-email/config"
	"smark.freecoop.net/grafana-email/datasource"
	"smark.freecoop.net/grafana-email/pdf"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	user := config.Instance.DBUser
	passwd := config.Instance.DBPassword
	dbHost := config.Instance.DBHost
	database := config.Instance.Database
	dsn := user + ":" + passwd + "@tcp(" + dbHost + ")/" + database + "?charset=utf8mb4&parseTime=True&loc=Local"
	//dsn := "root:1234@tcp(localhost:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"

	// Open a MySQL database connection.
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	DB = db
	if err != nil {
		panic("failed to connect to database")
	}
	db.AutoMigrate(&Schedule{})
	return DB
}

const (
	RecursDaily   = "Daily"
	RecursWeekly  = "Weekly"
	RecursMonthly = "Monthly"
)

type Schedule struct {
	gorm.Model
	DashboardID    string            `json:"dashboardID" gorm:"type:text"`
	DayOfMonth     int               `json:"dayOfMonth"`
	DayOfWeek      int               `json:"dayOfWeek"`
	DayTime        string            `json:"dayTime"`
	ExpireTime     string            `json:"expireTime"`
	From           string            `json:"from"`
	IsCSV          bool              `json:"isCSV"`
	IsExcel        bool              `json:"isExcel"`
	IsPDF          bool              `json:"isPDF"`
	IsPreview      bool              `json:"isPreview"`
	IsZip          bool              `json:"isZip"`
	LasteEmailTime string            `json:"lasteEmailTime"`
	Message        string            `json:"message"`
	Name           string            `json:"name"`
	NextEmailTime  string            `json:"nextEmailTime"`
	OrgID          string            `json:"orgID"`
	Recipients     string            `json:"recipients"`
	RecipientsBCC  string            `json:"recipients_BCC"`
	RecipientsCC   string            `json:"recipients_CC"`
	Recurs         string            `json:"recurs"`
	ReportType     string            `json:"reportType"`
	Result         string            `json:"result"`
	Scheduler      string            `json:"scheduler"`
	Subject        string            `json:"subject"`
	To             string            `json:"to"`
	User           string            `json:"user"`
	PanelIds       string            `json:"panelIds"`
	Plugins        string            `json:"plugins"`
	Variables      map[string]string `json:"variables" gorm:"-"`
	VarsData       string
	CronID         int
}

func (s *Schedule) EncodeDBModel() {
	data, _ := json.Marshal(s.Variables)
	s.VarsData = string(data)
}

func (s *Schedule) DecodeDBModel() {
	json.Unmarshal([]byte(s.VarsData), &s.Variables)
}

func (s *Schedule) ScheduleCronString() string {
	rs := ""
	dt := strings.Split(s.DayTime, ":")
	switch s.Recurs {
	case RecursDaily:
		rs = strings.Join([]string{dt[1], dt[0], "*", "*", "*"}, " ")
	case RecursWeekly:
		rs = strings.Join([]string{dt[1], dt[0], "*", "*", strconv.Itoa(s.DayOfWeek)}, " ")
	case RecursMonthly:
		rs = strings.Join([]string{dt[1], dt[0], strconv.Itoa(s.DayOfMonth), "*", "*"}, " ")
	}
	return rs
}

func (s Schedule) Run() {

	var ls Schedule
	DB.Where("id = ?", s.ID).Find(&ls)
	sdata, err := json.Marshal(ls)
	log.Println("run job,"+string(sdata), err)

	//send email
	s.SendEmail()

	ls.Result = "success"
	ls.NextEmailTime = time.Now().String()
	ls.LasteEmailTime = time.Now().String()
	DB.Updates(&ls)
}

func (s *Schedule) SendEmail() error {
	m := gomail.NewMessage()
	m.SetHeader("From", config.Instance.EmailFrom)
	m.SetHeader("To", strings.Split(s.Recipients, ",")...)
	m.SetHeader("Subject", s.Subject)
	m.SetBody("text/html", s.Message)
	// Create a new SMTP
	c := config.Instance
	d := gomail.NewDialer(c.SmtpHost, c.SmtpPort, c.EmailUserName, c.EmailPasswd)
	panels := datasource.DashboardPanels(s.OrgID, s.DashboardID)
	pd := pdf.InitPDF()
	for i, v := range panels {
		pid := strconv.Itoa(v)
		bytes := datasource.PanelImage(s.OrgID, s.DashboardID, pid, nil)
		if i == 0 {
			pdf.InsertImage("Test"+pid, pd, bytes, 70)
		} else {
			pdf.InsertImageInNewPage("Test"+pid, pd, bytes)
		}
	}
	pdfFile := s.DashboardID + ".pdf"
	pd.OutputFileAndClose(pdfFile)
	m.Attach(pdfFile)
	// Send the email
	if err := d.DialAndSend(m); err != nil {
		log.Fatal("Error sending email:", err)
		return err
	}

	log.Println("Email sent successfully.")
	return nil
}
