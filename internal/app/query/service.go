package query

import (
	"github.com/romberli/das/internal/dependency/query"
	"github.com/romberli/go-util/common"
	"github.com/romberli/go-util/constant"
)

const queryQueriesStruct = "Queries"

var _ query.Service = (*Service)(nil)

type Service struct {
	config  *Config
	dasRepo *DASRepo
	queries []query.Query
}

// NewService returns a new *Service
func NewService(config *Config, dasRepo *DASRepo) *Service {
	return newService(config, dasRepo)
}

// NewServiceWithDefault returns a new *Service with default repository
func NewServiceWithDefault(config *Config) *Service {
	return newService(config, NewDASRepoWithGlobal())
}

// newService returns a new *Service
func newService(config *Config, dasRepo *DASRepo) *Service {
	return &Service{
		config:  config,
		dasRepo: dasRepo,
	}
}

// GetConfig returns the config of query
func (s *Service) GetConfig() *Config {
	return s.config
}

// GetQueries returns the query slice
func (s *Service) GetQueries() []query.Query {
	return s.queries
}

// GetByMySQLClusterID gets the query slice by the mysql cluster identity
func (s *Service) GetByMySQLClusterID(mysqlClusterID int) error {
	var err error

	querier := NewQuerierWithGlobal(s.GetConfig())
	s.queries, err = querier.GetByMySQLClusterID(mysqlClusterID)
	if err != nil {
		return err
	}

	return s.Save(mysqlClusterID, constant.DefaultRandomInt, constant.DefaultRandomInt, constant.DefaultRandomString)
}

// GetByMySQLServerID gets the query slice by the mysql server identity
func (s *Service) GetByMySQLServerID(mysqlServerID int) error {
	var err error

	querier := NewQuerierWithGlobal(s.GetConfig())
	s.queries, err = querier.GetByMySQLServerID(mysqlServerID)
	if err != nil {
		return err
	}

	return s.Save(constant.DefaultRandomInt, mysqlServerID, constant.DefaultRandomInt, constant.DefaultRandomString)
}

// GetByDBID gets the query slice by the mysql server identity and the db identity
func (s *Service) GetByDBID(mysqlServerID int, dbID int) error {
	var err error

	querier := NewQuerierWithGlobal(s.GetConfig())
	s.queries, err = querier.GetByDBID(mysqlServerID, dbID)
	if err != nil {
		return err
	}

	return s.Save(constant.DefaultRandomInt, mysqlServerID, dbID, constant.DefaultRandomString)
}

// GetBySQLID gets the query by the mysql server identity and the sql identity
func (s *Service) GetBySQLID(mysqlServerID int, sqlID string) error {
	var err error

	querier := NewQuerierWithGlobal(s.GetConfig())
	s.queries, err = querier.GetBySQLID(mysqlServerID, sqlID)
	if err != nil {
		return err
	}

	return s.Save(constant.DefaultRandomInt, mysqlServerID, constant.DefaultRandomInt, sqlID)
}

// Save the query info into DAS repo
func (s *Service) Save(mysqlClusterID, mysqlServerID, dbID int, sqlID string) error {

	return s.dasRepo.Save(mysqlClusterID, mysqlServerID, dbID, sqlID,
		s.GetConfig().GetStartTime(), s.GetConfig().GetEndTime(), s.GetConfig().GetLimit(), s.GetConfig().GetOffset())
}

// Marshal marshals Service.Queries to json bytes
func (s *Service) Marshal() ([]byte, error) {

	return s.MarshalWithFields(queryQueriesStruct)
}

// MarshalWithFields marshals only specified fields of the Service to json bytes
func (s *Service) MarshalWithFields(fields ...string) ([]byte, error) {

	return common.MarshalStructWithFields(s, fields...)
}
