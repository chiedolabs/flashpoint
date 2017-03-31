# Flashpoint

<div style="text-align:center;">
<img src="logo/flashpoint.png" width="400" />
</div>

**Beta notice:** This CLI utility is still in beta and is rapidly changing. Many changes are breaking changes and require manual intervention.

A CLI utility for creating temporary Heroku deployments from a starter Heroku app and a git branch. The experience is similar to '[Heroku Review Apps](https://devcenter.heroku.com/articles/github-integration-review-apps)' but focuses on API-driven projects that consist of multiple git repositories.

# Table of Contents
1. [Installation and Upgrading](#a)
2. [Getting Started](#b)
3. [Using the tool](#c)
4. [The config file format](#d)
5. [Development](#e)
6. [Gotchas](#f)
7. [Todos](#g)

## <div id="a">Installation and Upgrading</div>

```
# The following will always be updated with the latest 'stable' version
# via the version number. You can replace the version number below with any version
# you want to install. When you run this, your current version will be overwritten

wget -O /usr/local/bin/flashpoint https://github.com/chiedolabs/flashpoint/blob/0.0.1/flashpoint?raw=true \
&& chmod +x /usr/local/bin/flashpoint
```

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

See a fleshed out example (Notice that we're using the file name but not the absolute path).

```
flashpoint example.json create
```

#### Cleaning up old apps

To delete apps created with this tool that have been inactive for over 5 days, you can run the script below.

```
flashpoint clean
```

Or to destroy all apps created with this tool run

```
flashpoint clean destroy_all
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

## <div id="e">Development</div>

- Install the packages `go get ./...`
- Make your changes
- Run `go build`

## <div id="f">Gotchas</div>

- The app are created on your personal Heroku account.
- The apps created with this tool don't have anything to do with Heroku Pipelines and don't automatically get deleted when you close a pull request.
- If the parent app has any addons, you'll need to make sure you have a credit card on file on your personal account even if there won't be any charges.
- We are using Heroku forking to copy the apps so be sure to read [this doc](https://devcenter.heroku.com/articles/fork-app) to be aware of the implications.
- Changes you make will not automatically be deployed to your app. You will need to manually push any changes you make after the initial deployment. This can be done with `git push -f <HEROKU_APP_NAME> <YOUR_BRANCHNAME>:master`
- It is still up to you to push your changes to github, etc.

## <div id="g">To-dos</div>

- Comment my functions
- Better docs
- Create an automation script for creating new projects.
- Complete beta and move to version 1.0

## Legal Stuff

- [LICENSE](./LICENSE)
- Logo <a href='http://www.freepik.com/free-vector/lightning-logo_719539.htm'>Designed by Freepik</a>
