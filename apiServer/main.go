package main

import (
	"fmt"
	"github.com/BambooTuna/go-server-lib/config"
	"github.com/BambooTuna/go-server-lib/connection/mysql"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"sync"
	"time"

	"github.com/BambooTuna/go-server-lib/metrics"
)

const namespace = "trade_bot"

func main() {
	wg := new(sync.WaitGroup)
	wg.Add(4)

	m := metrics.CreateMetrics(namespace)
	go func() {
		health := m.Gauge("health", map[string]string{})
		health.Set(200)
		ticker := time.NewTicker(time.Minute * 1)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				health.Set(200)
			}
		}
	}()

	mysqlConnection := mysql.GormConnection()
	defer mysqlConnection.Close()

	go func() {
		serverPort := config.GetEnvString("PORT", "18080")
		r := gin.Default()
		r.GET("/", func(ctx *gin.Context) { ctx.Status(200) })
		r.GET("/health", func(ctx *gin.Context) { ctx.Status(200) })
		_ = r.Run(fmt.Sprintf(":%s", serverPort))
		wg.Done()
	}()

	// monitoring metrics, process
	go func() {
		processCollector := prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{Namespace: namespace})
		prometheus.MustRegister(m, processCollector)
		http.Handle("/metrics", promhttp.Handler())
		_ = http.ListenAndServe(":2112", nil)
		wg.Done()
	}()
	wg.Wait()
}
