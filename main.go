package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
)

var (
	DOMAIN   string
	USERLIST string
	WORKERS  int
	LOGJSON  string
)

func CheckUser(users chan string, wg *sync.WaitGroup) {
	defer wg.Done()

	for username := range users {
		domain := strings.Split(DOMAIN, ".")

		userPrincipalName := strings.Replace(username, ".", "_", -1)
		for _, c := range domain {
			userPrincipalName += fmt.Sprintf("_%s", c)
		}

		url := fmt.Sprintf("https://%s-my.sharepoint.com/personal/%s/_layouts/15/onedrive.aspx", domain[0], userPrincipalName)

		req, err := http.NewRequest("HEAD", url, nil)
		if err != nil {
			log.Errorf("error creating request - err: %v", err)
			continue
		}

		client := http.Client{}

		rsp, err := client.Do(req)
		if err != nil {
			log.Errorf("error sending request - err: %v", err)
			continue
		}

		fmt.Printf("[*] Checking user - status_code=%d domain=%s username=%s \n", rsp.StatusCode, DOMAIN, username)

		log.WithFields(log.Fields{
			"username":    username,
			"domain":      DOMAIN,
			"status_code": rsp.StatusCode,
		}).Info("checking user")
	}
}

func InitLogger(logFile string) (*os.File, error) {
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return file, err
	}

	log.SetOutput(file)
	log.SetFormatter(&log.JSONFormatter{})

	return file, nil
}

func init() {
	flag.StringVar(&DOMAIN, "d", "", "Domain")
	flag.StringVar(&USERLIST, "u", "", "Userlist")
	flag.IntVar(&WORKERS, "w", 3, "Workers")
	flag.StringVar(&LOGJSON, "j", "results.json", "Log results to JSON file")
	flag.Parse()
}

func main() {
	var logfile *os.File

	var err error
	if logfile, err = InitLogger(LOGJSON); err != nil {
		log.Fatalf("error creating log file - err: %v", err)
	}
	defer logfile.Close()

	userlist, err := os.Open(USERLIST)
	if err != nil {
		log.Fatalf("error opening userlist - err: %v", err)
	}
	defer userlist.Close()

	var wg sync.WaitGroup
	wg.Add(WORKERS)

	users := make(chan string)
	for w := 0; w < WORKERS; w++ {
		go CheckUser(users, &wg)
	}

	scanner := bufio.NewScanner(userlist)
	for scanner.Scan() {
		users <- scanner.Text()
	}

	close(users)
	wg.Wait()
}
