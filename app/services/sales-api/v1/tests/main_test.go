package tests

import (
	"fmt"
	"testing"

	"ardanlabs/service/business/data/dbtest"
)

func TestMain(m *testing.M) {
	var err error
	c, err = dbtest.StartDB()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer dbtest.StopDB(c)

	m.Run()
}
