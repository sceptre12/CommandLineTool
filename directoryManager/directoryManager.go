 package directoryManager

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"commandLineTool/utils"
	"fmt"
)

const WorkSpaceName = "LearningWorkSpace"

func CreateProject(projectName, langType string, shouldActivate bool) bool{
	 workingDir := filepath.Join(getUserHome(),WorkSpaceName,langType)

	 // Creates the working dir
	checkAndCreateLangDir(workingDir)

	 projectDir := filepath.Join(workingDir,projectName)

	 // Check if the projectName already exists
	 if doesDirExist(projectDir) {
	 		fmt.Println("Project: ", projectName, " Already exists")
		 return false
	 }

	createDir(projectDir)

	 // Perform some special project initialization here
	if shouldActivate {
		activateProject(langType,projectDir)
	}

	 return true
}

func RemoveProject(projectName, langType string) bool{
	projectDir := filepath.Join(getUserHome(),WorkSpaceName,langType,projectName)

	if !doesDirExist(projectDir) {
		panic("Project does not exist: " + projectName)
	}

	removeDir(projectDir)

	return true
}

func LaunchProject(projectName, langType string ) bool{
	projectDir := filepath.Join(getUserHome(),WorkSpaceName,langType,projectName)

	if !doesDirExist(projectDir) {
		panic("Project does not exist: " + projectName)
	}

	activateProject(langType,projectDir)
	return true
}

// HELPER FUNCTIONS

func activateProject(langType,projectDir string){
	var projectEditor , args string
	switch langType{
	case utils.Js:
		projectEditor = "code"
		args = "--folder-uri"
	}

	if projectEditor != ""{
		cmdPtr := exec.Command(projectEditor,args,projectDir)
		var out bytes.Buffer
		cmdPtr.Stdout = &out
		err := cmdPtr.Run()
		if err != nil {
			panic("Error executing program " + err.Error())
		}

		output := out.String()
		if output != ""{
			fmt.Printf("in all caps: %q\n", output)
		}
	}
}

func getUserHome() (userDir string){
	userDir , err := os.UserHomeDir()
	if err != nil {
		panic("Error Getting user home dir:" + err.Error())
	}
	return
}

func checkAndCreateLangDir(langDir string){
	if !doesDirExist(langDir){
		createDir(langDir)
	}
}

func doesDirExist(dir string) bool{
	doesItExit := true
	if _, err := os.Stat(dir); os.IsNotExist(err){
		doesItExit = false
	}
	return doesItExit
}

func createDir(dir string){
	isError := os.MkdirAll(dir, os.ModePerm)
	if isError != nil {
		panic("Error Creating Dir: " + isError.Error())
	}
}

func removeDir(dir string){
	error := os.RemoveAll(dir)
	if error != nil {
		panic("Error removing dir :" + dir)
	}
}

func listDirFiles(dir string) {
	if !doesDirExist(dir){
		panic("Dir does not exist" + dir)
	}
	err := filepath.Walk(dir,func(path string, info os.FileInfo, err error) error{
		if err != nil{
			fmt.Println("Error Occurred: " + err.Error())
			return err
		}
		if info.IsDir(){
			fmt.Println("Dir: ", info.Name())
		}else {
			fmt.Println("File: ", info.Name())
		}

		return err
	})

	if err != nil {
		panic("Error occurred: " + err.Error())
	}

}

// TODO finish the below
func findDir(name,dir string, level int)bool{
	hasFound := false
	currentDirLevel := level -1


	err := filepath.Walk(dir,func(path string, info os.FileInfo, err error) error{
		if err != nil{
			fmt.Println("Error Occurred: " + err.Error())
			return err
		}
		if info.IsDir() && !hasFound{
			hasFound = findDir(name,filepath.Join(dir,info.Name()), level -1)
		}

		return err
	})

	if currentDirLevel == 0 {
		return false
	}

	if err != nil {
		panic("Error occurred: " + err.Error())
	}

	return hasFound
}

func digIntoDir(name string,dir string, level int) bool{
	currentLevel := level -1
	if FileInfo, err := os.Stat(dir); err != nil {
		if err != nil{
			panic("Error with file " + err.Error())
		}
		if FileInfo.Name() == name  {
			return false
		}else if currentLevel == 0 {
			return true
		} else if FileInfo.IsDir(){
			return digIntoDir(name,dir,currentLevel)
		}
	}
	return false
}