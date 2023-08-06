package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	"smark.freecoop.net/grafana-email/config"
	"smark.freecoop.net/grafana-email/model"
)

func main() {
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
			log.Println("Remove schedule id:%d, cronID:%d", v, s.CronID)
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

		log.Println("var ---:", s.Variables)
		fmt.Println(s.ScheduleCronString())
		if s.ID != 0 {
			var tmpS model.Schedule
			db.Where("id = ?", s.ID).Find(&tmpS)
			s.CronID = tmpS.CronID
			cr.Remove(cron.EntryID(s.CronID))
			id, _ := cr.AddJob(s.ScheduleCronString(), s)
			s.CronID = int(id)
			s.EncodeDBModel()
			db.Updates(&s)
			log.Println("Update schedule id:%d,raw cronID:%d, new cronID:%d", tmpS.CronID, s.ID, id)
		} else {
			db.Create(&s)
			id, _ := cr.AddJob(s.ScheduleCronString(), s)
			s.CronID = int(id)
			s.EncodeDBModel()
			db.Updates(&s)
			log.Println("New schedule id:%d,cronID:%d", s.ID, id)
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
