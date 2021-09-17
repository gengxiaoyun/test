package healthcheck

import (
	"time"
)

type ItemConfig interface {
	// GetID returns the identity
	GetID() int
	// GetItemName returns the item name
	GetItemName() string
	// GetItemWeight returns the item weight
	GetItemWeight() int
	// GetLowWatermark returns the low watermark
	GetLowWatermark() float64
	// GetHighWatermark returns the high watermark
	GetHighWatermark() float64
	// GetUnit returns the unit
	GetUnit() float64
	// GetScoreDeductionPerUnitHigh returns the score deduction per unit high
	GetScoreDeductionPerUnitHigh() float64
	// GetMaxScoreDeductionHigh returns the max score deduction high
	GetMaxScoreDeductionHigh() float64
	// GetScoreDeductionPerUnitMedium returns the score deduction per unit medium
	GetScoreDeductionPerUnitMedium() float64
	// GetMaxScoreDeductionMedium returns the max score deduction medium
	GetMaxScoreDeductionMedium() float64
	// GetDelFlag returns the delete flag
	GetDelFlag() int
	// GetCreateTime returns the create time
	GetCreateTime() time.Time
	// GetLastUpdateTime returns the last update time
	GetLastUpdateTime() time.Time
}

type EngineConfig interface {
	// GetItemConfig returns the item config
	GetItemConfig(item string) ItemConfig
	// Validate validates if engine configuration is valid
	Validate() error
}

type Variable interface {
	// GetName returns the name of the variable
	GetName() string
	// GetName returns the value of the variable
	GetValue() string
}

type Table interface {
	// GetSchema returns the table schema
	GetSchema() string
	// GetName returns the table name
	GetName() string
	// GetRows returns the table rows
	GetRows() int
	// GetSize returns the table size
	GetSize() float64
}

type FileSystem interface {
	GetMountPoint() string
	GetDevice() string
}

type PrometheusData interface {
	// GetTimestamp returns the timestamp
	GetTimestamp() string
	// GetValue returns the value
	GetValue() float64
}
