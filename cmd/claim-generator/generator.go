package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"

	"workout-plan/config"
	"workout-plan/plan"
)

var log = logrus.New()

func main() {
	conf, err := config.GetConfig()
	handleError(err)

	log.SetLevel(logrus.WarnLevel)

	err = plan.InitializeExerciseDefinitions(conf.Plans.DefinitionsFile, log)
	handleError(err)

	logrus.New()
	err = plan.InitializePlans(conf.Plans.Directory, log)
	handleError(err)

	var (
		subject     string
		rawLifetime string
	)

	plansSingleton := plan.GetPlansInstance()
	plans := plansSingleton.GetAll()
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("\nAvailable plans are: ")
	for _, planItem := range plans {
		fmt.Println(" - " + planItem.ID)
	}
	fmt.Println()

	fmt.Print("Enter Subject: ")
	subject, _ = reader.ReadString('\n')
	subject = strings.TrimRight(subject, "\n")

	fmt.Print("Enter Lifetime in seconds: ")
	rawLifetime, _ = reader.ReadString('\n')
	rawLifetime = strings.TrimRight(rawLifetime, "\n")

	lifetime, err := strconv.Atoi(rawLifetime)
	handleError(err)

	access := make([]struct {
		Type string
		Name string
	}, 0)

	fmt.Print("Enter comma seperated allowed plans: ")
	allowedPlans, _ := reader.ReadString('\n')
	allowedPlans = strings.TrimRight(allowedPlans, "\n")

	for _, allowedPlan := range strings.Split(allowedPlans, ",") {
		access = append(access, struct {
			Type string
			Name string
		}{Type: "plan", Name: allowedPlan})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iss":    conf.Auth.Token.Issuer,
		"sub":    subject,
		"aud":    conf.Auth.Token.Service,
		"exp":    time.Now().Unix() + int64(lifetime),
		"nbf":    time.Now().Unix(),
		"iat":    time.Now().Unix(),
		"access": access,
	})

	tokenString, err := token.SignedString([]byte(conf.Auth.Token.Secret))
	fmt.Printf("Token String:\n%v\n", tokenString)
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
