package db

import (
	"database/sql"
	"fmt"
	"log"
	"odj-deliver-cloudbuild/config"
	"odj-deliver-cloudbuild/model"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" //db utils
)

type Db struct {
	dbConn *sqlx.DB
	err    error
}

//InitDBConnection  initializes db connection
func (dbInfo *Db) InitDBConnection(postgresConf config.PostgresConfig, attempt uint8) {
	log.Println(fmt.Sprintf("User %s connecting to db %s:%v", postgresConf.User, postgresConf.Server, postgresConf.Port))

	db, err := sqlx.Connect("postgres",
		fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
			postgresConf.User, postgresConf.Password, postgresConf.Server, postgresConf.Port, postgresConf.Database))

	if err != nil {
		dbInfo.handleDBConnectionError(postgresConf, attempt)
	}
	dbInfo.dbConn = db
	dbInfo.err = err
	log.Println(fmt.Sprintf("User %s connected to db %s:%v", postgresConf.User, postgresConf.Server, postgresConf.Port))
}

//handleDBConnectionError retires connecting to DB if attempts are not exhausted as per user defined limit
func (dbInfo *Db) handleDBConnectionError(postgresConf config.PostgresConfig, attempt uint8) {
	if attempt < postgresConf.DBRetryAttempts {
		time.Sleep(time.Millisecond * 500)
		dbInfo.InitDBConnection(postgresConf, attempt+1)
	}
	log.Println("Connection Failed")
	panic(dbInfo.err)
}

//CloseDBConnection is to close DB connection if needed
func (dbInfo *Db) CloseDBConnection() {
	dbInfo.dbConn.Close()
}

//InsertBuildDetails is to insert commit details in build table
func (dbInfo *Db) InsertCloudBuild(build model.Build) (bool, error) {
	fmt.Println("build:" + build.ProductName)
	fmt.Println("build:" + build.ComponentName)
	sqlQuery := `INSERT INTO build (product_name, component_name, project_id, repo_name, repo_url, branch, commit_id, commit_message) VALUES (:product_name, :component_name, :project_id, :repository_name, :repository_url, :branch_name, :commit_id, :commit_msg)`
	_, dbInfo.err = dbInfo.dbConn.NamedExec(sqlQuery, build)
	if dbInfo.err != nil {
		fmt.Println("Error while Inserting Build Details: ", build)
		fmt.Println("Error Encountered: ", dbInfo.err)
		return false, dbInfo.err
	}
	fmt.Println("Inserted Successfully Build Details ", build)
	return true, nil
}

//Read Product and component details from component table
func (dbInfo *Db) ReadProducComponentAndDockerRegistry(repo_name string) model.Build {
	var build model.Build
	sqlStatement := "SELECT product_name, component_name, docker_registry, trigger_id FROM component WHERE repo_name=$1"
	row := dbInfo.dbConn.QueryRow(sqlStatement, repo_name)
	err := row.Scan(&build.ProductName, &build.ComponentName, &build.DockerRegistry, &build.TriggerID)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Zero rows found")
		} else {
			fmt.Println("Connection Error")
			panic(err)
		}
	} else {
		fmt.Println("fetched:", build)
	}
	return build

}

//Read ProjectId from Product table
func (dbInfo *Db) ReadProjectId(productName string) string {
	var projectId string
	sqlStatement := `SELECT project_id FROM product WHERE product_name=$1`
	row := dbInfo.dbConn.QueryRow(sqlStatement, productName)
	err := row.Scan(&projectId)

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Zero rows found")
		} else {
			fmt.Println("Connection Got Error")
			panic(err)
		}
	} else {
		fmt.Println("fetched:", projectId)
	}
	return projectId
}

//UpdateComponent is to insert TriggerId Component table
func (dbInfo *Db) UpdateComponentData(component model.ComponentResponse) (bool, error) {

	sqlQuery := `UPDATE public.component SET trigger_id=$1 WHERE component_name=$2 AND repo_name=$3`
	_, dbInfo.err = dbInfo.dbConn.Exec(sqlQuery, component.TriggerID, component.ComponentName, component.RepositoryName)
	if dbInfo.err != nil {
		fmt.Println("Error while Updating Build Details: ", component)
		fmt.Println("Error Encountered: ", dbInfo.err)
		return false, dbInfo.err
	}
	fmt.Println("Updated Successfully Build Details ", component)
	return true, nil
}

//Read BuildSeq from Build table
func (dbInfo *Db) ReadBuildSeq(build model.Build) int {
	var buildSeq int
	sqlStatement := `SELECT build_seq FROM build WHERE product_name=$1 and component_name=$2 and commit_id=$3`
	row := dbInfo.dbConn.QueryRow(sqlStatement, build.ProductName, build.ComponentName, build.CommitID)
	err := row.Scan(&buildSeq)

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Zero rows found")
		} else {
			fmt.Println("Connection Got Error")
			panic(err)
		}
	} else {
		fmt.Println("fetched:", buildSeq)
	}
	return buildSeq
}

//UpdateCloudBuild is to insert BuildId and BuildStatus in build table
func (dbInfo *Db) UpdateCloudBuild(build model.Build) (bool, error) {

	sqlQuery := `UPDATE public.build SET build_status=$1, build_id=$2, image_id=$3 WHERE build.build_seq=$4`
	_, dbInfo.err = dbInfo.dbConn.Exec(sqlQuery, build.BuildStatus, build.BuildID, build.ImageID, build.BuildSeq)
	if dbInfo.err != nil {
		fmt.Println("Error while Updating Build Details: ", build)
		fmt.Println("Error Encountered: ", dbInfo.err)
		return false, dbInfo.err
	}
	fmt.Println("Updated Successfully Build Details ", build)
	return true, nil
}
