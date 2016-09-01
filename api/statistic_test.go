package api

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/vmware/harbor/tests/apitests/apilib"
	"strconv"
	"testing"
)

func TestStatisticGet(t *testing.T) {

	fmt.Println("Testing Statistic API")
	assert := assert.New(t)

	apiTest := newHarborAPI()

	//prepare for test

	admin := &usrInfo{"admin", "Harbor12345"}

	var myProCount, pubProCount, totalProCount int32
	result, err := apiTest.StatisticGet(*admin)
	if err != nil {
		t.Error("Error while get statistic information", err.Error())
		t.Log(err)
	} else {
		myProCount = result.MyProjectCount
		pubProCount = result.PublicProjectCount
		totalProCount = result.TotalProjectCount
	}
	//post project
	var project apilib.Project
	project.ProjectName = "statistic_project"
	project.Public = 1

	//case 2: admin successful login, expect project creation success.
	fmt.Println("case 2: admin successful login, expect project creation success.")
	reply, err := apiTest.ProjectsPost(*admin, project)
	if err != nil {
		t.Error("Error while creat project", err.Error())
		t.Log(err)
	} else {
		assert.Equal(reply, int(201), "Case 2: Project creation status should be 201")
	}

	//get and compare
	result, err = apiTest.StatisticGet(*admin)
	if err != nil {
		t.Error("Error while get statistic information", err.Error())
		t.Log(err)
	} else {
		assert.Equal(myProCount+1, result.MyProjectCount, "MyProjectCount should be equal")
		assert.Equal(int32(2), result.MyRepoCount, "MyRepoCount should be equal")
		assert.Equal(pubProCount+1, result.PublicProjectCount, "PublicProjectCount should be equal")
		assert.Equal(int32(2), result.PublicRepoCount, "PublicRepoCount should be equal")
		assert.Equal(totalProCount+1, result.TotalProjectCount, "TotalProCount should be equal")
		assert.Equal(int32(2), result.TotalRepoCount, "TotalRepoCount should be equal")

	}

	//get the project
	var projects []apilib.Project
	httpStatusCode, projects, err := apiTest.ProjectsGet(project.ProjectName, 1)
	if err != nil {
		t.Error("Error while search project by proName and isPublic", err.Error())
		t.Log(err)
	} else {
		assert.Equal(int(200), httpStatusCode, "httpStatusCode should be 200")
		project.ProjectId = projects[0].ProjectId
	}

	//delete the project
	projectID := strconv.Itoa(int(project.ProjectId))
	httpStatusCode, err = apiTest.ProjectsDelete(*admin, projectID)
	if err != nil {
		t.Error("Error while delete project", err.Error())
		t.Log(err)
	} else {
		assert.Equal(int(200), httpStatusCode, "Case 1: Project creation status should be 200")
		//t.Log(result)
	}

	fmt.Printf("\n")

}