package query

import (
	"strconv"
	"time"

	"github.com/romberli/das/internal/app/query"
	"github.com/romberli/das/pkg/message"
	msgquery "github.com/romberli/das/pkg/message/query"
	"github.com/romberli/go-util/constant"
)

const (
	startTimeJSON = "start_time"
	endTimeJSON   = "end_time"
	limitJSON     = "limit"
	offsetJSON    = "offset"
)

func GetConfig(dataMap map[string]string) (*query.Config, error) {
	config := query.NewConfigWithDefault()

	// get start time
	startTimeStr, exists := dataMap[startTimeJSON]
	if exists {
		startTime, err := time.ParseInLocation(constant.TimeLayoutSecond, startTimeStr, time.Local)
		if err != nil {
			return nil, err
		}

		config.SetStartTime(startTime)
	}
	// get end time
	endTimeStr, exists := dataMap[endTimeJSON]
	if exists {
		endTime, err := time.ParseInLocation(constant.TimeLayoutSecond, endTimeStr, time.Local)
		if err != nil {
			return nil, err
		}

		config.SetEndTime(endTime)
	}
	// get limit
	limitStr, exists := dataMap[limitJSON]
	if exists {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			return nil, err
		}

		config.SetLimit(limit)
	}
	// get offset
	offsetStr, exists := dataMap[offsetJSON]
	if exists {
		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			return nil, err
		}

		config.SetLimit(offset)
	}
	// validate config
	if !config.IsValid() {
		return nil, message.NewMessage(msgquery.ErrQueryConfigNotValid, config.GetStartTime(), config.GetEndTime(), config.GetLimit())
	}

	return config, nil
}
