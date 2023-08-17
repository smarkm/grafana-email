package config

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

func TestConver(t *testing.T) {
	Init("../config")

	var r map[int]string
	v, _ := json.Marshal(*Instance)
	fmt.Println(string(v))
	log.Println(r)
}
