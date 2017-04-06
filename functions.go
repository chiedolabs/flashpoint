package main

import (
	"encoding/json"
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Checks if there's an error and if there is outputs
// that error in a pretty little format.
func check(e error) {
	checkWithMessage(e, "Something went wrong. Please try again.")
}

// Outputs a custom error message in a pretty format
func checkWithMessage(e error, message string) {
	red := color.New(color.FgRed)
	if e != nil {
		red.Print("Error: ")
		fmt.Println(message)
		os.Exit(0)
	}
}

// Destroys old apps created with flashpoint
func destroyOldApps(path string, f os.FileInfo, err error) error {
	boldGreen := color.New(color.FgGreen, color.Bold)

	if strings.HasSuffix(path, ".json") {
		config, err := ioutil.ReadFile(path)

		check(err)

		configAsBytes := []byte(config)

		// Have to use var because we have to initial value for this
		var configStruct Config

		// Get up in that JSON file and convert it into a Map
		if err := json.Unmarshal(configAsBytes, &configStruct); err != nil {
			log.Fatal("Your configuration file is incorrectly formed.")
		}

		// Loop through each app so we can open each flashpoint git remote
		// file and check each of the apps
		for _, app := range configStruct.Apps {
			readOnlyFlashpointRepos, err := ioutil.ReadFile(app.Path + "/.git/flashpointrepos")
			if err != nil {
				// This means the file doesn't exist so continue
				continue
			}

			// Get a list of apps from the flashpoint git remote file
			r, _ := regexp.Compile("remote \"(.*)\"")
			remoteList := r.FindAllStringSubmatch(string(readOnlyFlashpointRepos), -1)
			if len(remoteList) < 1 {
				boldGreen.Println(configStruct.Project + "[" + app.Name + "] has no running apps.")
				continue
			}

			// Get a slice of app names in the format we want it in
			appNames := []string{}
			for index, _ := range remoteList {
				appNames = append(appNames, remoteList[index][1])
			}

			// Using Heroku's API, check and see if any of the apps have been inactive
			// for a certain amount of days

			// Get the HEROKU token as we will need that for the CURL requests
			out, err := exec.Command("heroku", "auth:token").CombinedOutput()
			check(err)
			herokuToken := strings.Replace(string(out), "\n", "", -1)

			// We will store the list of repos to delete here
			toDelete := []string{}

			for _, name := range appNames {
				req, err := http.NewRequest("GET", "https://api.heroku.com/apps/"+name+"/dynos", nil)
				check(err)
				req.Header.Set("Accept", "application/vnd.heroku+json; version=3")
				req.Header.Set("Authorization", "Bearer "+herokuToken)

				resp, err := http.DefaultClient.Do(req)
				check(err)
				defer resp.Body.Close()

				bodyBytes, err2 := ioutil.ReadAll(resp.Body)
				check(err2)
				bodyString := string(bodyBytes)

				if resp.StatusCode == 200 { // OK
					// Do stuff with API response
					r, _ := regexp.Compile("\"updated_at\":\"(.*)\"")

					updatedAt := strings.Split(r.FindStringSubmatch(bodyString)[1], "\n")[0]

					// Weird how golang creates format strings but whatev
					timeFormat := "2006-01-02T15:04"

					// Need to remove the Z stuff so it can be parsed
					updatedAt = updatedAt[0 : len(updatedAt)-4]
					then, err := time.Parse(timeFormat, updatedAt)
					check(err)
					hoursInactive := time.Since(then).Hours()

					// If an app hasn't been active for more than this time,
					// delete it
					hoursInactiveBeforeDeletion := 5 * 24
					if (len(os.Args) >= 3 && os.Args[2] == "destroy_all") || int(hoursInactive) > hoursInactiveBeforeDeletion {
						toDelete = append(toDelete, name)
					}

				} else {
					if strings.Contains(bodyString, "not_found") {
						// if the app wasn't found, set it for deletion
						toDelete = append(toDelete, name)
					}
					fmt.Println(bodyString)
				}
			}

			if len(toDelete) == 0 {
				boldGreen.Println(configStruct.Project + " [" + app.Name + "] has no running apps that are expired.")
			} else {

				// Get the old content of the file so we can overwrite it
				fileContentArr := strings.Split(string(readOnlyFlashpointRepos), "\n")

				// We need to delete each app from Heroku and then update the
				// repos config
				for _, appToDelete := range toDelete {
					// Regexes weren't working for me so I did this hackish thing...
					for index, _ := range fileContentArr {
						if strings.Contains(fileContentArr[index], "remote \""+appToDelete+"\"") {
							// This means we've found a section we need to delete...
							fileContentArr[index-1] = "" // The comment line
							fileContentArr[index] = ""   // The remote line
							fileContentArr[index+1] = "" // And so on
							fileContentArr[index+2] = ""
						}
					}

					newFileContents := ""
					// Recreate new file contents
					for _, value := range fileContentArr {
						if value != "" {
							newFileContents += value + "\n"
						}
					}

					// Now that we're ready, destroy the app on Heroku
					out, _ := exec.Command("heroku", "apps:destroy", "--app", appToDelete, "--confirm", appToDelete).CombinedOutput()
					fmt.Println(string(out))

					// Save the changes
					err = ioutil.WriteFile(app.Path+"/.git/flashpointrepos", []byte(newFileContents), 0664)
					check(err)

				}
				boldGreen.Println(configStruct.Project + " [" + app.Name + "] SUCCESS: " + strconv.Itoa((len(toDelete))) + " apps deleted.\n\n")
			}
		}
	}

	return nil
}

// Replaces custom flashpoint variables in a string with
// the appropriate values
func evaluateVars(value string, reviewAppNames []string) string {
	for iii, _ := range reviewAppNames {
		varStr := ""

		// Make REVIEW_APP_NAMES a config variable
		varStr = "$REVIEW_APP_NAMES[" + strconv.Itoa(iii) + "]"
		value = strings.Replace(value, varStr, reviewAppNames[iii], -1)

		// Make REVIEW_APP_URLS a config variable
		varStr = "$REVIEW_APP_URLS[" + strconv.Itoa(iii) + "]"
		urlStr := "https://" + reviewAppNames[iii] + ".herokuapp.com"
		value = strings.Replace(value, varStr, urlStr, -1)

	}

	return value
}

// Returns a unique alpha string
func uniqueString() string {
	const chars = "abcdefghijklmnopqrstuvwxyz0123456789"
	strlen := 6
	result := ""
	for i := 0; i < strlen; i++ {
		index := r.Intn(len(chars))
		result += chars[index : index+1]
	}
	return result
}
