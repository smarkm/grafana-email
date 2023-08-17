package model

import (
	"reflect"
	"testing"

	"gorm.io/gorm"
	"smark.freecoop.net/grafana-email/config"
)

func TestInitDB(t *testing.T) {
	tests := []struct {
		name string
		want *gorm.DB
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := InitDB(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InitDB() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSchedule_ScheduleCronString(t *testing.T) {
	tests := []struct {
		name string
		s    *Schedule
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.ScheduleCronString(); got != tt.want {
				t.Errorf("Schedule.ScheduleCronString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSchedule_Run(t *testing.T) {
	tests := []struct {
		name string
		s    Schedule
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.s.Run()
		})
	}
}

func TestSendEmail(t *testing.T) {
	config.Init("../config.json")

	type args struct {
		s Schedule
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test basic email",
			args: args{s: Schedule{
				Subject:     "Test emai",
				Recipients:  "smark@ossera.com",
				Message:     "<b>Hello</b>",
				OrgID:       "1",
				DashboardID: "Gfgpou3Vk",
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.s.SendEmail()
		})
	}
}
