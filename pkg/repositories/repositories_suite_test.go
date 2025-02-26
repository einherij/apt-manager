package repositories

import (
	"database/sql"
	"fmt"
	"path"
	"runtime"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"

	"github.com/einherij/apt-manager/pkg/config"
)

const testDBName = "apt_manager_test"

type RepositoriesSuite struct {
	suite.Suite

	psql *sql.DB
}

func TestRepositoriesSuite(t *testing.T) {
	suite.Run(t, new(RepositoriesSuite))
}

func (s *RepositoriesSuite) SetupSuite() {
	conf := config.NewConfig()
	s.Assert().NoError(conf.ParseEnv())

	psql, err := sql.Open("postgres", conf.Postgres.PostgresConnection())
	s.Assert().NoError(err)

	_, err = psql.Exec("CREATE DATABASE " + testDBName)
	s.Assert().NoError(err)

	conf.Postgres.DB = testDBName
	psql, err = sql.Open("postgres", conf.Postgres.PostgresConnection())
	s.Assert().NoError(err)

	m, err := migrate.New(
		fmt.Sprintf("file://%s/migrations", GetProjectRootPath()),
		conf.Postgres.PostgresConnection(),
	)
	s.Assert().NoError(err)
	err = m.Up()
	s.Assert().NoError(err)

	s.psql = psql
}

func (s *RepositoriesSuite) TearDownSuite() {
	s.Assert().NoError(s.psql.Close())

	conf := config.NewConfig()
	s.Assert().NoError(conf.ParseEnv())
	psql, err := sql.Open("postgres", conf.Postgres.PostgresConnection())
	s.Assert().NoError(err)
	_, err = psql.Exec("DROP DATABASE " + testDBName)
	s.Assert().NoError(err)
}

func GetProjectRootPath() string {
	const projectName = "apt-manager"
	_, currentFile, _, _ := runtime.Caller(1)
	for i := 10; path.Base(currentFile) != projectName && i > 0; i-- {
		currentFile = path.Dir(currentFile)
	}
	return currentFile
}
