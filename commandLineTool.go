// Package CommandLineTool will be used to help me automate some system tasks
//
package main

import (
	"commandLineTool/directoryManager"
	"commandLineTool/utils"
	"flag"
	"fmt"
	"os"
)

func main() {
	setupUserInputs()
}

//
// https://github.com/jawher/mow.cli
func setupUserInputs() {
	// SubCommands
	createProject := flag.NewFlagSet(utils.CreateProject, flag.ContinueOnError)
	launchProject := flag.NewFlagSet(utils.LaunchProject, flag.ContinueOnError)
	deleteProject := flag.NewFlagSet(utils.DeleteProject, flag.ContinueOnError)
	searchForProject := flag.NewFlagSet(utils.SearchForProject, flag.ContinueOnError)

	// Create Project subCommands
	projectName := createProject.String("name", "", "Project Name (Required)")
	projectType := createProject.String("language", "", "Projects main programming Language Ex input: (js|go|python|ruby) (Required)")
	launchDefaultEditor := createProject.Bool("launchDefaultEditor", false, "Should vscode be launched")


	// // Launch Project subcommands
	openProject := launchProject.String("name", "", "Project to be opened (Required)")
	searchArea := launchProject.String("language", "", "Projects main programming Language Ex input: (js|go|python|ruby)")

	// Delete Project subcommands
	projectToBeDeleted := deleteProject.String("name", "", "Project to be delete ")
	projectLanguage := deleteProject.String("language", "", "Projects main programming Language Ex input: (js|go|python|ruby)")

	// Search for project
	//projectToBeSearched := searchForProject.String("name", "", "Project to be Searched ")
	//projectSearchArea := searchForProject.String("language", "", "Projects main programming Language Ex input: (js|go|python|ruby)")

	// Verify that a subcommand has been provided
	// os.Arg[0] is the main command
	// os.Arg[1] will be the subcommand
	if len(os.Args) < 2 {
		fmt.Println("You need to enter one of the sub commands: createProject | launchProject | deleteProject | searchForProject ")
		os.Exit(1)
	}

	/**
	Switch on the sub command
	Pare the flags
	os.Args[2:] will be all arguments starting after the subcommand at os.Args[1]
	*/

	switch os.Args[1] {
	case utils.CreateProject:
		createProject.Parse(os.Args[2:])
	case utils.LaunchProject:
		launchProject.Parse(os.Args[2:])
	case utils.DeleteProject:
		deleteProject.Parse(os.Args[2:])
	case utils.SearchForProject:
		searchForProject.Parse(os.Args[2:])
	case "help", "--help", "-help":
		printDefaultVals(createProject,launchProject,deleteProject,searchForProject)
	default:
		printDefaultVals(createProject,launchProject,deleteProject,searchForProject)
		panic("Wrong arg passed: " +  os.Args[1],)
	}


	if createProject.Parsed() {
		handleCreateProject(createProject,*projectName,*projectType,*launchDefaultEditor)
	}else if deleteProject.Parsed(){
		handleRemoveProject(deleteProject, *projectToBeDeleted, *projectLanguage)
	}else if launchProject.Parsed(){
		handleLaunchProject(createProject,*openProject,*searchArea)
	}else if searchForProject.Parsed(){
		fmt.Println("Searching for project not handled")
	}

}

func printDefaultVals(createProject *flag.FlagSet,launchProject *flag.FlagSet,deleteProject *flag.FlagSet,searchForProject *flag.FlagSet){
	fmt.Println("Create Project Commands:")
	createProject.PrintDefaults()
	fmt.Println("\nLaunch Project Commands:")
	launchProject.PrintDefaults()
	fmt.Println("\nDelete Project Commands:")
	deleteProject.PrintDefaults()
	fmt.Println("\nSearch For Project Commands:")
	searchForProject.PrintDefaults()
}

func isNotValidArgs( name string, langType string) bool{
	return  ( name == "" || langType == "" )  || (langType != utils.Js && langType != utils.Go && langType != utils.Java && langType != utils.Python)
}

func errorCheck(params *flag.FlagSet,name,langType string){
	if isNotValidArgs(name,langType) {
		params.PrintDefaults()
		panic("One of these args aren't valid: " + name + " : " + langType)
	}
}

func handleCreateProject(params *flag.FlagSet, name string, langType string, shouldLaunch bool){
		// Checking to make sure the required params are filled
		errorCheck(params,name,langType)

		// Create Project folder and activate (not working)
		if directoryManager.CreateProject(name,langType,shouldLaunch){
			fmt.Println("Success...")
		}
}

func handleRemoveProject(params *flag.FlagSet, name string, langType string){
	errorCheck(params,name,langType)

	if directoryManager.RemoveProject(name,langType){
		fmt.Println("Success...")
	}
}

func handleLaunchProject(params *flag.FlagSet, projectName string, langType string){
	// Checking to make sure the required params are filled
	errorCheck(params,projectName,langType)

	if directoryManager.LaunchProject(projectName,langType){
		fmt.Println("Success...")
	}
}