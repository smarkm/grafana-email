package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/denisbrodbeck/machineid"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"smark.freecoop.net/grafana-email/config"
	"smark.freecoop.net/grafana-email/datasource"
	"smark.freecoop.net/grafana-email/model"
	"smark.freecoop.net/grafana-email/pdf"
)

var machineID string
var scheduleLimit = 2
var (
	version bool
)

func parseArgs() {
	flag.BoolVar(&version, "version", false, "Show version")
	flag.Parse()

	if version {
		log.Println("Free v1.0.0 version with 2 email schedules limit")
		os.Exit(0)
	}
}

func main() {
	machineID = generateMachineID()
	parseArgs()
	VerifyLisence()

	config.Init("./config.json")
	db := model.InitDB()
	router := gin.Default()
	cr := cron.New()
	cr.Start()
	var ss []model.Schedule
	db.Find(&ss)
	for i := 0; i < len(ss); i++ {
		s := ss[i]

		id, err := cr.AddJob(s.ScheduleCronString(), s)
		if err != nil {
			log.Printf("Add schedule got error: %s", err.Error())
		}
		s.CronID = int(id)
		db.Updates(&s)
	}
	log.Println("Update schedules", len(cr.Entries()), len(ss))
	// Define a route for the root URL "/"
	router.GET("/api/schedule/getScheduleList", func(c *gin.Context) {
		var ss []model.Schedule
		db.Where(&model.Schedule{User: "admin"}).Find(&ss)
		c.JSON(200, gin.H{
			"code": 0,
			"msg":  "success",
			"data": ss,
		})
	})
	router.GET("/api/schedule/delete", func(c *gin.Context) {
		ids := c.QueryArray("ids")
		log.Println(ids)
		for _, v := range ids {
			id, _ := strconv.Atoi(v)
			var s model.Schedule
			db.Where("id = ?", id).Find(&s)
			db.Delete(&s)
			cr.Remove(cron.EntryID(s.CronID))
			log.Printf("Remove schedule id:%s, cronID:%d", v, s.CronID)
		}
		c.JSON(200, gin.H{
			"code": 0,
			"msg":  "success",
			"data": nil,
		})
	})
	router.POST("/api/schedule/add", func(c *gin.Context) {
		var s model.Schedule
		if err := c.ShouldBindJSON(&s); err != nil {
			c.JSON(http.StatusOK, gin.H{"code": 1, "msg": "error" + err.Error()})
			return
		}
		if len(cr.Entries()) > scheduleLimit && !strings.EqualFold("now", s.Recurs) {
			msg := "Your backend only support " + strconv.Itoa(scheduleLimit) + " email schedule with test model, please contact provider to get full version"
			c.JSON(http.StatusOK, gin.H{"code": 1, "msg": msg})
			return
		}
		if s.Recurs == "now" {
			err := s.SendEmail()
			if err != nil {
				msg := "Faild send email: " + err.Error()
				c.JSON(http.StatusOK, gin.H{"code": 1, "msg": msg})
				return
			}
		} else if s.ID != 0 {
			var tmpS model.Schedule
			db.Where("id = ?", s.ID).Find(&tmpS)
			s.CronID = tmpS.CronID
			cr.Remove(cron.EntryID(s.CronID))
			id, _ := cr.AddJob(s.ScheduleCronString(), s)
			s.CronID = int(id)
			s.EncodeDBModel()
			db.Updates(&s)
			log.Printf("Update schedule id:%d,raw cronID:%d, new cronID:%d", tmpS.CronID, s.ID, id)
		} else {
			db.Create(&s)
			id, _ := cr.AddJob(s.ScheduleCronString(), s)
			s.CronID = int(id)
			s.EncodeDBModel()
			db.Updates(&s)
			log.Printf("New schedule id:%d,cronID:%d", s.ID, id)
		}
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success"})
	})
	router.GET("/api/schedule/tasks", func(c *gin.Context) {
		//data, _ := json.Marshal(cr.Entries())
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "success", "data": cr.Entries()})
	})
	router.GET("/api/schedule/pdf", func(c *gin.Context) {
		orgID, _ := c.GetQuery("orgId")
		dashboardID, _ := c.GetQuery("uid")
		queryParams := make(map[string]string)

		for key, values := range c.Request.URL.Query() {
			if len(values) > 0 {
				queryParams[key] = values[0]
			}
		}
		title, panels := datasource.DashboardPanels(orgID, dashboardID)
		log.Printf("dashboard:%s,%v", title, panels)
		pd := pdf.InitPDF(title)
		for i, v := range panels {
			pid := strconv.Itoa(v)
			bytes := datasource.PanelImage(orgID, dashboardID, pid, queryParams)
			if i == 0 {
				pdf.InsertImage(pid, pd, bytes, 70)
			} else {
				pdf.InsertImageInNewPage(pid, pd, bytes)
			}
		}
		pdfFile := title + ".pdf"
		err := pd.OutputFileAndClose(pdfFile)
		if err != nil {
			log.Panicln("Generate PDF got error:" + err.Error())
		}
		fileContent, err := ioutil.ReadFile(pdfFile)
		if err != nil {
			fmt.Println("Error reading file:", err)
			return
		}
		c.Header("Content-Type", "application/pdf")
		// Set the content disposition to make the browser download the file.
		c.Header("Content-Disposition", "attachment; filename="+pdfFile)
		// Write the PDF content as the response body.
		c.Data(http.StatusOK, "application/pdf", fileContent)
	})

	// Run the server on port 8080
	router.Run(":" + strconv.Itoa(config.Instance.HTTPPort))

	// Schedule the job to run every minute

	// Start the cron scheduler
	//eyJrIjoiNGd0RElKaE1OcWtnY3Z1dDBhYTJiSDJTanNPaEVQN1oiLCJuIjoidGVzdCIsImlkIjoxfQ==

	select {}
}
func VerifyLisence() {
	userCode := "cb95d602-6eaf-4bfe-9e9d-18882738d49e"
	if strings.EqualFold(userCode, machineID) {
		allowTime := "2023-08-08 15:04:05"
		parsedTime, _ := time.Parse("2006-01-02 15:04:05", allowTime)
		if time.Now().After(parsedTime) {
			log.Println("Your lisence got expired by" + allowTime + ", running with test model only " + strconv.Itoa(scheduleLimit) + " schedules allowed")
		} else {
			scheduleLimit = math.MaxInt64
			//pass lisence verify
			log.Println("Your lisence will be expire by " + allowTime)
		}

	} else {
		log.Println("No lisenced machine")

		os.Exit(0)
	}
}
func generateMachineID() string {
	// 获取机器的唯一标识符，例如 MAC 地址
	machineID, err := machineid.ID()
	if err != nil {
		log.Panicln("Failed to get machineID")
		os.Exit(1)
	}
	log.Println("machineID: " + machineID)
	return machineID
}
