package app

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"pubsub-assignment/internal/config"
	"pubsub-assignment/internal/domain"
	"pubsub-assignment/internal/domain/pubsub"
	"pubsub-assignment/internal/domain/queueservice"
	"pubsub-assignment/internal/server"
	"pubsub-assignment/internal/server/schema"
	"strconv"
	"strings"
	"time"
)

//go:embed testdata
var testFiles embed.FS

type app struct {
	config *config.Config
	server *server.Server

	queueService domain.FileQueueService
	broker       domain.Broker
	subscriber   domain.Sub
}

func (a *app) initConfig() {
	c, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("cannot initialize config for app: %v\n", err)
	}
	a.config = c
}

func (a *app) initApp() {
	a.initConfig()
	a.initBroker()
	a.initSubscriber()
	a.initQueueService()
	a.initHTTPServer()
}

func (a *app) initQueueService() {
	a.queueService = queueservice.New(a.broker, a.subscriber)
}

func (a *app) initBroker() {
	a.broker = pubsub.NewBroker()
}

func (a *app) initSubscriber() {
	a.subscriber = pubsub.NewSubscriber(uuid.NewString())
	a.broker.Subscribe(a.subscriber, domain.FileTopic)
}

func (a *app) initHTTPServer() {
	s, err := server.New(a.config, a.queueService)
	if err != nil {
		log.Fatalf("cannot init server: %v\n", err)
	}
	a.server = s
}

func (a *app) start() {
	go a.RWFiles()
	a.server.Start()
}

func Run() {
	a := app{}
	a.initApp()
	a.start()
}

func (a *app) RWFiles() {
	url := "http://" + net.JoinHostPort(a.config.Server.Host, a.config.Server.Port) + "/line"
	for i := 0; i < 3; i++ {
		//i := rand.Intn(3)
		time.Sleep(time.Duration(i+1) * time.Second)
		filename := "file" + strconv.Itoa(i+1) + ".txt"
		f, err := testFiles.ReadFile("testdata/" + filename)
		if err != nil {
			log.Fatal(err)
		}

		payload := schema.File{
			Name:  filename,
			Lines: strings.Split(string(f), "\n"),
		}
		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			log.Fatal(err)
		}

		req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonPayload))
		if err != nil {
			log.Fatal(err)
		}

		req.Header.Set("Content-Type", "application/json")
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(resp.Status)
		_ = resp.Body.Close()
	}
	// Read file
	for i := 0; i < 3; i++ {
		resp, err := http.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		var fileResp schema.FileResponse
		if err := json.Unmarshal(body, &fileResp); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("file No.%d", i)
		fmt.Println(fileResp)

		file, err := os.Create(a.config.ConsumedDirPath + fileResp.Name)
		if err != nil {
			log.Fatal(err)
		}
		_, err = file.WriteString(fileResp.Content)
		if err != nil {
			log.Fatal(err)
		}
		_ = resp.Body.Close()
		_ = file.Close()
	}
}
