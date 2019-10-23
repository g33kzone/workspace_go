package model

//Struct used to create trigger
type BuildTrigger struct {
	ProductName    string `db:"product_name" json:"product_name"`       //product_name
	ComponentName  string `db:"component_name" json:"component_name"`   //component_name
	ProjectID      string `db:"project_id" json:"project_id"`           //project_id
	RepositoryName string `db:"repository_name" json:"repository_name"` //repository_name
	RepositoryUrl  string `db:"repository_url" json:"repository_url"`   //repository_url
	DockerRegistry string `db:"docker_registry" json:"docker_registry"` //docker_registry
	BranchName     string `db:"banch_name" json:"branch_name"`          //branch_name
	Filename       string `db:"file_name" json:"file_name"`             //file_name
	Repo           string
	UserName       string
}

//Struct used to store Build Details in Build Table
type Build struct {
	BuildSeq       int    `db:"build_seq" json:"build_seq"`             //build_seq
	ProductName    string `db:"product_name" json:"product_name"`       //product_name
	ComponentName  string `db:"component_name" json:"component_name"`   //component_name
	ProjectID      string `db:"project_id" json:"project_id"`           //project_id
	RepositoryName string `db:"repository_name" json:"repository_name"` //repository_name
	RepositoryUrl  string `db:"repository_url" json:"repository_url"`   //repository_url
	CommitID       string `db:"commit_id" json:"commit_id"`             //commit_id
	CommitMsg      string `db:"commit_msg" json:"commit_msg"`           //commit_msg
	BranchName     string `db:"branch_name" json:"branch_name"`         //branch_name
	BuildID        string `db:"build_id" json:"build_id"`               //build_id
	BuildStatus    string `db:"build_status" json:"build_status"`       //build_status
	ImageID        string `db:"image_id" json:"image_id"`               //image_id
	DockerRegistry string `json:"docker_registry"`                      //docker_registry
	TriggerID      string `json:"trigger_id"`                           //trigger_id
	StartTime      string `db:"start_time" json:"start_time"`           //start_time
	EndTime        string `db:"end_time" json:"end_time"`               //end_time
}

//Store Request deatils to create trigger
type ComponentResponse struct {
	ComponentName  string `json:"component_name"`
	ProductName    string `json:"product_name"`
	ProjectID      string `json:"project_id"`
	RepositoryName string `json:"repository_name"`
	RepositoryURL  string `json:"repository_url"`
	TriggerID      string `json:"trigger_id"`
}

type Commit struct {
	CommitID  string `json:"hash"`
	CommitMsg string `json:"message"`
}
type CommitBranch struct {
	Branch string `json:"name"`
}
type Change struct {
	Closed       bool         `json:"closed"`
	Commits      []Commit     `json:"commits"`
	CommitBranch CommitBranch `json:"new"`
}

type Push struct {
	Changes []Change `json:"changes"`
}
type Html struct {
	RepositoryUrl string `json:"href"`
}
type Link struct {
	Html Html `json:"html"`
}
type Repository struct {
	Name string `json:"full_name"`
	Link Link   `json:"links"`
}

//Bitbucket Response
type Response struct {
	Push       Push       `json:"push"`
	Repository Repository `json:"repository"`
}

type GitCommits struct {
	CommitId  string `json:"id"`
	CommitMsg string `json:"message"`
}

type GitRepository struct {
	Name      string `json:"name"`
	RepoUrl   string `json:"url"`
	GitBranch string `json:"master_branch"`
}

//Github Response
type GitResponse struct {
	GitCommits    []GitCommits  `json:"commits"`
	GitRepository GitRepository `json:"repository"`
}
