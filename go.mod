module github.com/romberli/das

go 1.16

replace (
	github.com/ClickHouse/clickhouse-go v1.4.7 => github.com/romberli/clickhouse-go v1.4.4-0.20210902113008-bb38dc6f756d
	github.com/spf13/pflag v1.0.5 => github.com/romberli/pflag v1.0.6-alpha
)

require (
	github.com/alecthomas/template v0.0.0-20190718012654-fb15b899a751
	github.com/asaskevich/govalidator v0.0.0-20210307081110-f21760c49a8d
	github.com/gin-gonic/gin v1.7.4
	github.com/hashicorp/go-multierror v1.1.1
	github.com/hashicorp/go-version v1.3.0
	github.com/jinzhu/now v1.1.2
	github.com/romberli/go-util v0.3.11
	github.com/romberli/log v1.0.20
	github.com/spf13/cast v1.4.1
	github.com/spf13/cobra v1.2.1
	github.com/spf13/viper v1.9.0
	github.com/stretchr/testify v1.7.0
	github.com/swaggo/files v0.0.0-20210815190702-a29dd2bc99b2
	github.com/swaggo/gin-swagger v1.3.2
	github.com/swaggo/swag v1.7.3
	go.uber.org/zap v1.19.1
)
