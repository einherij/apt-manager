package repositories

import (
	"database/sql"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/suite"
	"path"
	"runtime"
	"testing"

	"github.com/einherij/apt-manager/pkg/config"
)

const testDBName = "apt_manager_test"

type RepositoriesSuite struct {
	suite.Suite

	psql     *sql.DB
	masterDB *sql.DB
}

func TestRepositoriesSuite(t *testing.T) {
	suite.Run(t, new(RepositoriesSuite))
}

func (s *RepositoriesSuite) SetupSuite() {
	conf := config.NewConfig()
	s.Assert().NoError(conf.ParseEnv())

	fmt.Println(conf.Postgres.PostgresConnection())
	masterDB, err := sql.Open("postgres", conf.Postgres.PostgresConnection())
	s.Assert().NoError(err)

	_, err = masterDB.Exec("CREATE DATABASE " + testDBName)
	s.Assert().NoError(err)
	s.masterDB = masterDB

	testPGConfig := conf.Postgres
	testPGConfig.DB = testDBName
	fmt.Println(testPGConfig.PostgresConnection())
	psql, err := sql.Open("postgres", testPGConfig.PostgresConnection())
	s.Assert().NoError(err)

	m, err := migrate.New(
		fmt.Sprintf("file://%s/migrations", GetProjectRootPath()),
		testPGConfig.PostgresConnection(),
	)
	s.Assert().NoError(err)
	err = m.Up()
	s.Assert().NoError(err)
	srcErr, dbErr := m.Close()
	s.Assert().NoError(srcErr)
	s.Assert().NoError(dbErr)

	s.psql = psql
}

func (s *RepositoriesSuite) TearDownSuite() {
	s.Assert().NoError(s.psql.Close())

	_, err := s.masterDB.Exec("DROP DATABASE " + testDBName)
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
