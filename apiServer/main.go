package main

import (
	"fmt"
	"github.com/BambooTuna/go-server-lib/config"
	"github.com/BambooTuna/go-server-lib/connection/mysql"
	"github.com/BambooTuna/go-server-lib/connection/redis"
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

	redisConnection := redis.RedisConnection(0)
	defer redisConnection.Close()

	go func() {
		serverPort := config.GetEnvString("PORT", "18080")
		r := gin.Default()
		r.GET("/", func(ctx *gin.Context) { ctx.Status(200) })
		r.GET("/health", func(ctx *gin.Context) { ctx.Status(200) })
		r.GET("/health/mysql", func(ctx *gin.Context) {
			type Record struct{ ID string }
			var record Record
			sql := "select * from sample_table;"
			if err := mysqlConnection.Raw(sql).Scan(&record).Error; err != nil {
				ctx.String(500, fmt.Sprintf("Mysql select error (%s)", err.Error()))
				return
			}
			if record.ID != "f0c28384-3aa4-3f87-9fba-66a0aa62c504" {
				ctx.String(400, fmt.Sprintf("Mysql selected value (%s) != (f0c28384-3aa4-3f87-9fba-66a0aa62c504)", record.ID))
			}
			ctx.Status(200)
		})
		r.GET("/health/redis", func(ctx *gin.Context) {
			ctx.String(200, fmt.Sprintf("Redis DBSize is (%d)", redisConnection.DBSize().Val()))
		})

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
