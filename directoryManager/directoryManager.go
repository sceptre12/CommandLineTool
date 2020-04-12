 package directoryManager

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"commandLineTool/utils"
	"fmt"
	"strings"
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

func ListAllProjects(){
	navigateThroughDir(filepath.Join(getUserHome(),WorkSpaceName), nil, nil,nil,nil,2)
}

func SearchForProject(projectName, language string){
	isCorrectFolder := func (folderName string) bool{
		return folderName == projectName
	}

	isCorrectLanguageFolder := func (path string) bool{
		return strings.Contains(filepath.ToSlash(path),"/" + language)
	}

	cbFunc := func (_, dirName string){
		fmt.Println("Project Found: " + dirName)
	}

	if projectName != "" && language != "" {
		navigateThroughDir(filepath.Join(getUserHome(),WorkSpaceName), cbFunc,isCorrectLanguageFolder,isCorrectFolder, nil,2)
	}else if projectName != "" && language == "" {
		navigateThroughDir(filepath.Join(getUserHome(),WorkSpaceName), cbFunc,nil,isCorrectFolder, nil,2)
	}else if projectName == "" && language !="" {
		navigateThroughDir(filepath.Join(getUserHome(),WorkSpaceName), nil,isCorrectLanguageFolder, nil, nil,2)
	}

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


func navigateThroughDir(dir string, cbFunc func(fileName, dirName string),isCorrectLanguageFolder, isCorrectFolder, isCorrectFile func(string)bool, maxLevel int) (bool,error) {
	if !doesDirExist(dir){
		panic("Dir does not exist: " + dir)
	}

	basePathLength := len(strings.Split(filepath.ToSlash(dir),"/"))

	// The walk function literally walks through all files and dir starting from and including the root path
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error{
		if err != nil{
			fmt.Println("Error Occurred: " + err.Error())
			return err
		}

		_, endOfPath := filepath.Split(dir)
		if endOfPath == info.Name() {
			return nil
		}

		// Don't dig pass the max folder level
		if len(strings.Split(filepath.ToSlash(path),"/")) - basePathLength > maxLevel{
			return filepath.SkipDir
		}


		if info.IsDir() {
			// Ensures the search only occurs in the language dir
			if isCorrectLanguageFolder != nil && !isCorrectLanguageFolder(path){
				return filepath.SkipDir
			}

			if isCorrectFolder != nil && isCorrectFolder(info.Name()){
				cbFunc("",path)
				return filepath.SkipDir
			}

			// Prints out the dir info if the conditions hold true
			if isCorrectLanguageFolder != nil && isCorrectFolder == nil || isCorrectLanguageFolder == nil && isCorrectFolder == nil{

				isProjectName := checkAndPrintProjName(info.Name())

				if !isProjectName {
					fmt.Println(" - " + info.Name())
				}
			}
		}else{
			if isCorrectFolder != nil && isCorrectFile(info.Name()){
				cbFunc(info.Name(),"")
				return filepath.SkipDir
			}else{
				fmt.Println(" -- " + info.Name())
			}
		}

		return nil
	})

	if err != nil {
		return false, err
	}

	return true, nil
}

func checkAndPrintProjName(folder string) bool{
	isFolderAProjectFolder := false
	switch folder{
	case utils.Js:
		fmt.Println("\nJs Projects")
		isFolderAProjectFolder = true
	case utils.Java:
		fmt.Println("\nJava Projects")
		isFolderAProjectFolder = true
	case utils.Go:
		fmt.Println("\nGo Projects")
		isFolderAProjectFolder = true
	case utils.Python:
		fmt.Println("\nPython Projects")
		isFolderAProjectFolder = true
	}

	return isFolderAProjectFolder
}
