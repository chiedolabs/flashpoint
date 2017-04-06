package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"github.com/mitchellh/go-homedir"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

var r *rand.Rand

// Our json format
type Config struct {
	Project string `json:"project"`
	Apps    []struct {
		Name          string            `json:"name"`
		ParentAppName string            `json:"parent_app_name"`
		Path          string            `json:"path"`
		Env           map[string]string `json:"env"`
		Scripts       map[string]string `json:"scripts"`
	} `json:"apps"`
}

func init() {
	r = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func main() {
	if len(os.Args) >= 2 && os.Args[1] == "--version" {
		fmt.Println("0.0.3")
		os.Exit(0)
	}

	/////////////////////////////////////////////////
	// PREPARE SOME OUTPUT COLORS
	/////////////////////////////////////////////////
	boldGreen := color.New(color.FgGreen, color.Bold)
	boldWhite := color.New(color.FgWhite, color.Bold)
	boldBlue := color.New(color.FgBlue, color.Bold)
	boldRed := color.New(color.FgRed, color.Bold)

	/////////////////////////////////////////////////
	// START THE TIMER
	////////////////////////////////////////////////

	start := time.Now()

	if len(os.Args) >= 3 && os.Args[2] == "create" {
		/////////////////////////////////////////////////////////////
		// CREATE - The create command in the CLI util
		// This creates the projects on Heroku from the templates.
		////////////////////////////////////////////////////////////

		///////////////////////////////////////////
		// LOAD THE CONFIG FILE
		//////////////////////////////////////////
		file := os.Args[1]

		// Add the .json extension to the file name if it was not included.
		if strings.HasSuffix(file, ".json") != true {
			file += ".json"
		}

		homeDir, err := homedir.Dir()
		config, err := ioutil.ReadFile(homeDir + "/.flashpoint/projects/" + file)

		// If the file does not exist, end the script
		checkWithMessage(err, "That configuration file does not exist.")

		configAsBytes := []byte(config)

		// Have to use var because we have to initial value for this
		var configStruct Config

		// Get up in that JSON file and convert it into a Map
		if err := json.Unmarshal(configAsBytes, &configStruct); err != nil {
			log.Fatal("Your configuration file is incorrectly formed.")
		}

		// The number of apps
		numOfApps := len(configStruct.Apps)

		////////////////////////////////////////
		// DETERMINE AND STORE THE BRANCHES
		////////////////////////////////////////

		// Create a place to store the branches
		branches := make([]string, numOfApps)

		boldBlue.Println("\nCHOOSE YOUR BRANCHES")
		fmt.Println("=================================")
		for index, app := range configStruct.Apps {
			boldWhite.Print("\nAPP: ")
			fmt.Println(app.Name)
			reader := bufio.NewReader(os.Stdin)

			// Make it store the current branch by default. The following is
			// the long and convoluted way
			// of piping commands
			os.Chdir(app.Path)
			c1 := exec.Command("git", "symbolic-ref", "HEAD")
			c2 := exec.Command("sed", "s!refs\\/heads/!!")

			r, w := io.Pipe()
			c1.Stdout = w
			c2.Stdin = r

			var b2 bytes.Buffer
			c2.Stdout = &b2

			c1.Start()
			c2.Start()
			c1.Wait()
			w.Close()
			c2.Wait()
			x, err := ioutil.ReadAll(&b2)
			check(err)
			branches[index] = strings.Replace(string(x), "\n", "", -1)

			// Allow the user to specify a branch and if the user
			// specified, overide the default
			boldWhite.Print("\nWhich branch would you like to use?")
			fmt.Printf(" (Press enter to use '%s'):", branches[index])
			text, _ := reader.ReadString('\n')

			if len(text) > 1 {
				branches[index] = text

				//TODO: Check branch for validity
			}
		}

		/////////////////////////////////////////////////
		// START TICKER. (Output dots)
		////////////////////////////////////////////////
		ticker := time.NewTicker(time.Millisecond * 500)
		go func() {
			for range ticker.C {
				fmt.Printf(".")
			}
		}()

		/////////////////////////////////////////////////
		// SET UP REVIEW APP NAMES
		////////////////////////////////////////////////
		// We are setting this up first to make sure that every app has access
		// to all the review app names and URLS they need.
		reviewAppNames := make([]string, numOfApps)
		for index, app := range configStruct.Apps {
			// save the app name. Must start with a letter
			reviewAppName := "z" + uniqueString() + "-" + app.ParentAppName

			// We'll need to limit the number of chars in the name (Heroku limit)
			if len(reviewAppName) > 30 {
				reviewAppName = reviewAppName[0:30]
			}

			// Save the name for later use
			reviewAppNames[index] = reviewAppName
		}

		/////////////////////////////////////////////////
		// START WORKING ON THE APPS
		////////////////////////////////////////////////
		for index, app := range configStruct.Apps {
			/////////////////////////////////////////////////
			// FORK THE PARENT APPS
			////////////////////////////////////////////////
			boldGreen.Print("\n[" + app.Name + "] ")
			boldBlue.Println("FORKING YOUR HEROKU APP")
			fmt.Println("=================================")
			os.Chdir(app.Path)
			out, err := exec.Command("heroku", "fork", "--from", app.ParentAppName, "--to", reviewAppNames[index]).CombinedOutput()
			fmt.Println(string(out))
			check(err)

			/////////////////////////////////////////////////
			// ADD THE SAME COLLABORATORS
			////////////////////////////////////////////////
			out8, err8 := exec.Command("heroku", "access", "--app", app.ParentAppName).CombinedOutput()
			check(err8)
			for _, collaboratorStr := range strings.Split(string(out8), "\n") {
				collaborator := strings.Split(collaboratorStr, " ")[0]
				if strings.Contains(collaborator, "@") {
					_, err := exec.Command("heroku", "access:add", collaborator, "--app", reviewAppNames[index]).CombinedOutput()

					checkWithMessage(err, "There was an issue adding the collaborator "+collaborator)
				}
			}

			/////////////////////////////////////////////////
			// ADD THE GIT REMOTES
			////////////////////////////////////////////////
			// Check if the file already contains the include path if it doesn't then add it to the file
			readOnlyConfig, err := ioutil.ReadFile(app.Path + "/.git/config")
			check(err)

			if strings.Contains(string(readOnlyConfig), "\n[include]\n  path = ./flashpointrepos") == false {
				// Open the file
				gitConfig, err := os.OpenFile(app.Path+"/.git/config", os.O_APPEND|os.O_WRONLY, 0664)
				check(err)

				defer gitConfig.Close()

				if _, err := gitConfig.WriteString("\n[include]\n  path = ./flashpointrepos"); err != nil {
					check(err)
				}
			}

			// Now add the remote manually by adding it to our file
			flashpointRepos, err := os.OpenFile(app.Path+"/.git/flashpointrepos", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0664)
			check(err)

			defer flashpointRepos.Close()

			if _, err := flashpointRepos.WriteString(fmt.Sprintf("\n#DO NOT MANUALLY EDIT THESE IN ANY WAY\n[remote \"%s\"]\n  url = https://git.heroku.com/%s.git\n  fetch = +refs/heads/*:refs/remotes/%s/*", reviewAppNames[index], reviewAppNames[index], reviewAppNames[index])); err != nil {
				check(err)
			}

			//////////////////////////////////////////
			// SET ENV VARIABLES
			//////////////////////////////////////////
			boldGreen.Print("\n[" + app.Name + "] ")
			boldBlue.Println("SETTING YOUR ENVIRONMENT VARIABLES")
			fmt.Println("=================================")
			if app.Env != nil {
				args := []string{"config:set", "--app", reviewAppNames[index]}
				for key, value := range app.Env {
					// evaluate the variables in the string
					value = evaluateVars(value, reviewAppNames)

					// Add to the args
					args = append(args, key+"="+value)
				}

				out, err := exec.Command("heroku", args...).CombinedOutput()
				fmt.Println(string(out))
				check(err)
			}

			//////////////////////////////////////////////
			// PUSH LOCAL CHANGES TO THE NEW APP
			///////////////////////////////////////////////
			boldGreen.Print("\n[" + app.Name + "] ")
			boldBlue.Println("PUSHING YOUR BRANCH TO HEROKU.")
			fmt.Println("=================================")
			out2, err2 := exec.Command("git", "push", "-f", "--no-verify", reviewAppNames[index], branches[index]+":master").CombinedOutput()
			fmt.Println(string(out2))
			check(err2)

			////////////////////////////////////
			// RUN SCRIPTS
			////////////////////////////////////
			boldGreen.Print("\n[" + app.Name + "] ")
			boldBlue.Println("RUNNING YOUR SCRIPTS")
			fmt.Println("=================================")
			if app.Scripts != nil {
				for env, script := range app.Scripts {
					// evaluate the variables in the string
					script = evaluateVars(script, reviewAppNames)

					if env == "remote" {
						out, err := exec.Command("heroku", "run", "--app", reviewAppNames[index], script).CombinedOutput()
						fmt.Println(string(out))
						check(err)
					} else {
						out, err := exec.Command("bash", "-c", script).CombinedOutput()
						fmt.Println(string(out))
						check(err)
					}
				}
			}
		}

		/////////////////////////////////////////////////
		// STOP TICKER.
		////////////////////////////////////////////////
		ticker.Stop()

		///////////////////////////////////////////////////////
		// USER REPORT
		//////////////////////////////////////////////////////
		boldBlue.Println("\nALL DONE")
		fmt.Println("=================================")
		for index, app := range configStruct.Apps {
			boldGreen.Println(app.Name)
			boldWhite.Print("URL: ")
			fmt.Printf("%s.herokuapp.com\n", reviewAppNames[index])
			boldWhite.Print("Branch: ")
			fmt.Printf("%s", branches[index])
			boldWhite.Print("\nUpdate command: ")
			fmt.Printf("git push -f %s %s:master\n\n", reviewAppNames[index], branches[index])
		}
		boldGreen.Println("Remember to clean up after yourself every now and then dear by running the following:")
		fmt.Println("flashpoint clean")
		boldRed.Println("I wont clean up for you! >:( ")
	} else if len(os.Args) >= 2 && os.Args[1] == "clean" {
		/////////////////////////////////////////////////////////////
		// CLEAN - The clean command for the cli util. Deletes all apps
		// created by this tool that haven't been accessed in the last
		// five days
		////////////////////////////////////////////////////////////

		homeDir, err := homedir.Dir()

		check(err)

		filepath.Walk(homeDir+"/.flashpoint/projects/", destroyOldApps)

	} else {
		log.Fatal("Not a valid command")
	}

	/////////////////////////////////////////////////
	// SHOW TIME ELAPSED
	////////////////////////////////////////////////

	elapsed := time.Since(start)
	boldGreen.Printf("\n\nFlashpoint took %f seconds\n", elapsed.Seconds())
}
