# Flashpoint

<div style="text-align:center;">
<img src="logo/flashpoint.png" width="400" />
</div>

**Beta notice:** This CLI utility is still in beta and is rapidly changing. Many changes are breaking changes and require manual intervention.

A CLI utility for creating temporary Heroku deployments from a starter Heroku app and a git branch. The experience is similar to '[Heroku Review Apps](https://devcenter.heroku.com/articles/github-integration-review-apps)' but focuses on API-driven projects that consist of multiple git repositories.

# Table of Contents
1. [Requirements](#aa)
1. [Installation and Upgrading](#a)
1. [Pros and cons](#aa)
1. [Getting Started](#b)
1. [Using the tool](#c)
1. [The config file format](#d)
1. [Development](#e)
1. [Gotchas](#f)
1. [Todos](#g)

## <div id="aa">Requirements</div>

- [Heroku cli version 5.8.1+](https://devcenter.heroku.com/articles/heroku-cli)
- [wget](https://www.gnu.org/software/wget/)
  - OSX: `brew install wget`
  - Ubuntu: `sudo apt-get install wget`

## <div id="a">Installation and Upgrading</div>

```
# The following will always be updated with the latest 'stable' version
# via the version number. You can replace the version number below with any version
# you want to install. When you run this, your current version will be overwritten

wget -O /usr/local/bin/flashpoint https://github.com/chiedolabs/flashpoint/blob/0.0.4/flashpoint?raw=true \
&& chmod +x /usr/local/bin/flashpoint
```

## <div id="aa">Pros and Cons </div>

#### Flashpoint Pros:

- Your apps aren't automatically deleted after 5 days of inactivity, you can run the `clean` command to clean all apps that have been inactive for five days or more but at your own leisure.
- Trivial to start using. You don't have to change anything on Heroku to do so.
- You can use Flashpoint for projects that consist of multiple git repos
- You can automate environment variables based on the newly created app names from other apps in your project.

#### Flashpoint Cons:

- Your apps aren't automatically deleted after 5 days of inactivity. You have to run the `clean` command on occasions or run with a cronjob.
- All collaborators on the parent project get an invitation to collaborate on each newly created "review app".
- Apps aren't automatically created with pull requests nor are they automatically deleted.
- You have to manually push to github and your review app remote when you make changes

## <div id="b">Getting started</div>

Create a directory in your root folder named `.flashpoint` along with a `projects` sub-directory.

```
mkdir -p ~/.flashpoint/projects
```

Create a project

- The first step is that you create template heroku apps for each of your apps to inherit from. You don't want to use the production heroku app as the parent app (template).
- You can then create a project by creating an ~/.flashpoint/projects/<PROJECTNAME>.json file (eg. ~/.flashpoint/projects/example.json)
- Use the content in [the example config](./example-config.json) as a starting point.

## <div id="c">Using the tool</div>

#### Deploying a project from your current branch(s)

```
flashpoint <PROJECT_FILE_NAME> create
```
> **Note** Each time you run this (whether you've created a flashpoint deployment from the same branches already or not), a new group of Heroku apps are deployed (or just one Heroku app if that's all your project contains). See [Updating](#updating), to update an existing deployment.

See a fleshed out example (Notice that we're using the file name but not the absolute path).

```
flashpoint example.json create
```
#### <a name="updating"></a>Updating Exsisting Projects
To Update an existing project run:

```
git push -f <HEROKU_APP_NAME> <LOCALBRANCH>:master
```
The update command will also be given to you after creating your project.

#### Cleaning up old apps

To delete apps created with this tool that have been inactive for over 5 days, you can run the script below.

```
flashpoint clean
```

Or to destroy all apps created with this tool run

```
flashpoint clean --destroy_all
```

If you don't want to worry about running that command manually, you can schedule it with a cronjob by entering `crontab -e` in the terminal and then adding the following content and then saving the file.

```
# This will run every day at noon. Notice that if your computer is shut down or sleeping at the time, the script will not run until the next day at the specified time.
0 12 * * *  flashpoint clean
```

## <div id="d">Understanding the json config file format</div>

- **project** - The name of your project. You can make it whatever you want
- **apps** - Each of your apps included in this project. This will be an array as there can be many.
    - **name** - The name of the app
    - **parent_app_name** - The name of the app on Heroku that will be used as the template (leave out the herokuapp.com portion)
    - **path** - The absolute system path of this app's git repository on your machine.
    - **scripts** - An array of scripts to run on the heroku app after you deploy it. This is useful for doing migrations, seeds, etc.

#### Using alternative heroku mysql databases

If you aren't using the Heroku default postgresql setup and instead are using ClearDB for example, you'll need to make this your first item under `scripts` for the backend so that the `DATABASE_URL` is set to the correct value after the new app is created. It's hackish but it works. It's sets the `DATABASE_URL` on the newly created heroku app to the `CLEARDB_DATABASE_URL` which contains the newly created database's URL. Without this, `CLEARDB_DATABASE_URL` will be correct on the new app but `DATABASE_URL` will not be.

```
"local": "heroku config:set --app $REVIEW_APP_NAMES[0] DATABASE_URL=$(echo $(heroku config:get --app $REVIEW_APP_NAMES[0] CLEARDB_DATABASE_URL) | sed 's/mysql/mysql2/g')",

```

## <div id="e">Development</div>

- Install the packages `go get ./...`
- Make your changes
- Test with `go run *.go`
- Run `go build` to build the binary

## <div id="f">Gotchas</div>

- The app are created on your personal Heroku account.
- The apps created with this tool don't have anything to do with Heroku Pipelines and don't automatically get deleted when you close a pull request.
- If the parent app has any addons, you'll need to make sure you have a credit card on file on your personal account even if there won't be any charges.
- We are using Heroku forking to copy the apps so be sure to read [this doc](https://devcenter.heroku.com/articles/fork-app) to be aware of the implications.
- It is still up to you to push your changes to github, etc.

## <div id="g">To-dos</div>

- Create an automation script for creating new projects.
- Complete beta and move to version 1.0

## Legal Stuff

- [LICENSE](./LICENSE)
- Logo <a href='http://www.freepik.com/free-vector/lightning-logo_719539.htm'>Designed by Freepik</a>
