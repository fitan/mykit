module github.com/fitan/mykit

go 1.17

require (
	github.com/bytedance/go-tagexpr/v2 v2.9.5
	github.com/fsnotify/fsnotify v1.5.4
	github.com/go-kit/kit v0.12.0
	github.com/go-kit/log v0.2.1
	github.com/go-resty/resty/v2 v2.7.0
	github.com/google/gops v0.3.25
	github.com/google/wire v0.5.0
	github.com/gorilla/handlers v1.5.1
	github.com/gorilla/mux v1.8.0
	github.com/mattn/go-colorable v0.1.13
	github.com/pkg/errors v0.9.1
	github.com/prometheus/client_golang v1.13.0
	github.com/pterm/pterm v0.12.47
	github.com/pyroscope-io/pyroscope v0.29.0
	github.com/slok/go-http-metrics v0.10.0
	github.com/spf13/cobra v1.5.0
	github.com/spf13/viper v1.13.0
	github.com/tidwall/gjson v1.14.3
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.32.0
	go.opentelemetry.io/otel v1.11.1
	go.opentelemetry.io/otel/exporters/jaeger v1.10.0
	go.opentelemetry.io/otel/sdk v1.10.0
	go.uber.org/zap v1.23.0
	gopkg.in/natefinch/lumberjack.v2 v2.0.0
	gorm.io/driver/mysql v1.3.6
	gorm.io/gorm v1.23.10
)

replace go.opentelemetry.io/otel v1.11.1 => go.opentelemetry.io/otel v1.10.0

require (
	github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d
	github.com/hashicorp/consul/api v1.12.0
	go-micro.dev/v4 v4.8.1
	go.opentelemetry.io/otel/trace v1.11.1
)

require go.opentelemetry.io/contrib/instrumentation/github.com/gorilla/mux/otelmux v0.32.0

require (
	github.com/cockroachdb/errors v1.9.0
	github.com/google/uuid v1.3.0
	go.opentelemetry.io/otel/metric v0.32.0 // indirect
)
