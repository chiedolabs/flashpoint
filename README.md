# Flashpoint

**Beta notice:** This CLI utility is still in beta and is rapidly changing. Many changes are breaking changes and require manual intervention. 

An experience similar to 'Heroku review apps' with the focus being API-driven applications that consist of multiple git repositories.

Instead of only having one staging version of your site on Heroku, you now have the option of creating a separate staging version for each branch that needs to be reviewed.

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
wget -O /usr/local/bin/flashpoint https://github.com/chiedolabs/flashpoint/raw/master/flashpoint?date=$(date +%s) && chmod +x /usr/local/bin/flashpoint
```

## <div id="b">Getting started</div>

Create a directory in your root folder named .flashpoint

```
mkdir ~/.flashpoint
```

Create a project

- You can create a project by creating an ~/.flashpoint/<PROJECTNAME>.json file (eg. ~/.flashpoint/example.json)
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

- Add clean up script for apps. The user should be asked if he would like to clean up upon running the create script.
- Add to homebrew for better installation.
- Create an automation script for creating new projects.