package main

import (
	"flag"
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
	"smark.freecoop.net/grafana-email/model"
)

var machineID string
var scheduleLimit = 2

func main() {
	machineID = generateMachineID()

	var version bool
	flag.BoolVar(&version, "version", false, "Show version")
	flag.Parse()

	if version {
		log.Println("Free v1.0.0 version with 2 email schedules limit")
		os.Exit(0)
	}

	VerifyLisence()
	config.Init()
	db := model.InitDB()
	router := gin.Default()
	cr := cron.New()
	cr.Start()
	var ss []model.Schedule
	db.Find(&ss)
	for i := 0; i < len(ss); i++ {
		s := ss[i]
		id, _ := cr.AddJob(s.ScheduleCronString(), s)
		s.CronID = int(id)
		db.Updates(&s)
	}
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
		if len(cr.Entries()) >= scheduleLimit {
			msg := "Your backend only support " + strconv.Itoa(scheduleLimit) + " email schedule"
			c.JSON(http.StatusOK, gin.H{"code": 1, "msg": msg})
			return
		}
		if s.ID != 0 {
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

	// Run the server on port 8080
	router.Run(":8001")

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
			log.Println("Your lisence got expired")
			os.Exit(0)
		}

		//pass lisence verify
		scheduleLimit = math.MaxInt64
		log.Println("Your lisence will be expire by " + allowTime)
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
