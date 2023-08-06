package config

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

func TestConver(t *testing.T) {
	v := "{\"1\":1}"

	var r map[int]string
	e := json.Unmarshal([]byte(v), &r)
	fmt.Println(e)
	log.Println(r)
}
