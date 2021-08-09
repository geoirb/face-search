package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fasthttp/router"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/google/uuid"
	"github.com/kelseyhightower/envconfig"
	"github.com/valyala/fasthttp"

	"github.com/geoirb/face-search/internal/chromedp"
	"github.com/geoirb/face-search/internal/config"
	search "github.com/geoirb/face-search/internal/face-search"
	"github.com/geoirb/face-search/internal/face-search/server/http"
	"github.com/geoirb/face-search/internal/file"
	"github.com/geoirb/face-search/internal/mongo"
	"github.com/geoirb/face-search/internal/parser"
	"github.com/geoirb/face-search/internal/plugin"
	"github.com/geoirb/face-search/internal/proxy"
	"github.com/geoirb/face-search/internal/response"
	"github.com/geoirb/face-search/internal/result"
)

type configuration struct {
	HttpPort string `envconfig:"HTTP_PORT" default:"8081"`

	ConfigFile string `envconfig:"CONFIG_FILE" default:"./config.yml"`

	StorageConnect     string        `envconfig:"STORAGE_CONNECT" default:"mongodb://face-search:face-search@127.0.0.1:27017"`
	StorageDatabase    string        `envconfig:"STORAGE_DATABASE" default:"face-search"`
	StorageCollection  string        `envconfig:"STORAGE_COLLECTION" default:"result"`
	StorageSaveTimeout time.Duration `envconfig:"STORAGE_SAVE_TIMEOUT" default:"30s"`

	DownloadDir     string `envconfig:"DOWNLOAD_DIR" default:"/tmp/"`
	PluginDirLayout string `envconfig:"PLUGIN_DIR_LAYOUT" default:"/tmp/"`
}

const (
	prefixCfg   = ""
	serviceName = "face-search"
)

func main() {
	logger := log.NewJSONLogger(log.NewSyncWriter(os.Stdout))
	logger = log.WithPrefix(logger, "service", serviceName)
	logger = log.With(logger, "ts", log.DefaultTimestampUTC)

	if time.Since(time.Date(2020, time.August, 15, 0, 0, 0, 0, time.Now().Location())) < 0 {
		level.Error(logger).Log("msg", "trial version")
		return
	}

	var cfg configuration
	if err := envconfig.Process(prefixCfg, &cfg); err != nil {
		level.Error(logger).Log("msg", "configuration", "err", err)
		os.Exit(1)
	}

	level.Error(logger).Log("msg", "initialization")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	storage, err := mongo.NewStorage(
		ctx,
		cfg.StorageConnect,
		cfg.StorageDatabase,
		cfg.StorageCollection,
	)
	if err != nil {
		level.Error(logger).Log("msg", "init mongo", "err", err)
		os.Exit(1)
	}

	searchCfg, err := config.Read(cfg.ConfigFile)
	if err != nil {
		level.Error(logger).Log("msg", "init search config", "err", err)
		os.Exit(1)
	}

	searchConfig := search.Config{
		Timeout: searchCfg.Search.Timeout,
		Actions: make([]search.Action, 0, len(searchCfg.Search.Actions)),
	}
	for _, a := range searchCfg.Search.Actions {
		searchConfig.Actions = append(searchConfig.Actions, search.Action(a))
	}

	file := file.NewFacade(
		cfg.DownloadDir,
	)

	plugin := plugin.NewBuilder(
		proxy.New(),
		cfg.PluginDirLayout,
	)

	resultFacade := result.NewFacade(
		time.Now().Unix,
		uuid.NewString,

		storage,
	)

	chromedp := chromedp.New(
		plugin,
	)

	parser, err := parser.New(
		`<div class="card-vk01-header">([^<]*)<\/div><div class="card-vk01-score">Совпадение: <span class="score-label">([0-9]{1,2}[.]?[0-9]{1,2}%)<\/span><\/div><div class="[^<]*">[^<]*<\/div><div class="card-vk01-geo">[^<]*<\/div><div class="btn-vk01-container"><a href="(https:\/\/vk.com\/[^"]*)" target="_blank" class="btn btn-primary btn-vk01">Профиль<\/a><a href="#" data-target="#modalIMG" data-toggle="modal" class="btn btn-primary btn-vk01" data-imgsrc="[^"]*" data-imghref="(https:\/\/vk.com\/[^"]*)">Фото<\/a>`,
	)
	if err != nil {
		level.Error(logger).Log("msg", "init parser", "err", err)
		os.Exit(1)
	}

	svc := search.NewService(
		searchConfig,
		cfg.StorageSaveTimeout,

		file,
		resultFacade,
		chromedp,
		parser,

		logger,
	)

	svc = search.NewValidator(svc)

	router := router.New()
	http.Routing(router, svc, response.Build)

	httpServer := &fasthttp.Server{
		Handler:          router.Handler,
		DisableKeepalive: true,
	}

	go func() {
		level.Info(logger).Log("msg", "http server turn on", "port", cfg.HttpPort)
		if err := httpServer.ListenAndServe(":" + cfg.HttpPort); err != nil {
			level.Error(logger).Log("msg", "http server turn on", "err", err)
			os.Exit(1)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	level.Info(logger).Log("msg", "received signal", "signal", <-c)

	if err := httpServer.Shutdown(); err != nil {
		level.Info(logger).Log("msg", "http server shoutdown", "err", err)
	}
}
