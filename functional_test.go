// Functional test will use the real integration environment.package main
// We'll have real database server and real HTTP session.

package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/suite"
)

type FunctionalTestSuite struct {
	suite.Suite
	hostIP      string
	servicePath string
}

func TestFunctionalTestSuite(t *testing.T) {
	s := new(FunctionalTestSuite)
	s.hostIP = os.Getenv("ADDRESS")
	s.servicePath = "/"
	if len(os.Getenv("SERVICE_PATH")) > 0 {
		s.servicePath = os.Getenv("SERVICE_PATH")
	}
	suite.Run(t, s)
}

func (s FunctionalTestSuite) Test_Index_ReturnsStatus200() {
	address := fmt.Sprintf("http://%s%s/hello", s.hostIP, s.servicePath)
	log.Printf("Sending a request to %s\n", address)
	resp, err := http.Get(address)

	s.NoError(err)
	s.Equal(200, resp.StatusCode, "ADDR: ", address)
}

func (s FunctionalTestSuite) Test_IndexGreetings_ReturnsStatus204() {
	address := fmt.Sprintf("http://%s%s/greetings", s.hostIP, s.servicePath)
	log.Printf("Sending a request to %s\n", address)
	resp, err := http.Get(address)

	s.NoError(err)
	s.Equal(204, resp.StatusCode, "ADDR: %s", address)
}
