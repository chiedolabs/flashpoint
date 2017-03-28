# Flashpoint

An experience similar to 'Heroku review apps' with the focus being API-driven applications that consist of multiple git repositories.

Instead of only having one staging version of your site on Heroku, you now have the option of creating a separate staging version for each branch that needs to be reviewed.

## Installation and Upgrading

```
wget -O /usr/local/bin/flashpoint https://github.com/chiedolabs/flashpoint/raw/master/flashpoint?date=$(date +%s) && chmod +x /usr/local/bin/flashpoint
```

## Getting started

Create a directory in your root folder named .flashpoint

```
mkdir ~/.flashpoint
```

Create a project

- You can create a project by creating an ~/.flashpoint/<PROJECTNAME>.json file (eg. ~/.flashpoint/example.json)
- Use the content in [the example config](./example-config.json) as a starting point.

## Deploy your project

```
flashpoint <PROJECT_FILE_NAME> create
# An example would be. Notice that we're using the file name
# flashpoint example.json create
```

## Understanding the json config file format

- **project** - The name of your project. You can make it whatever you want
- **apps** - Each of your apps included in this project. This will be an array as there can be many.
    - **name** - The name of the app
    - **parent_app_name** - The name of the app on Heroku that will be used as the template (leave out the herokuapp.com portion)
    - **path** - The absolute system path of this app's git repository on your machine.
    - **scripts** - An array of scripts to run on the heroku app after you deploy it. This is useful for doing migrations, seeds, etc.

## Development

- Make your changes
- Run `go build`

## Gotchas

- Review apps are created on your personal Heroku account.
- Review apps created with this tool don't have anything to do with Heroku Pipelines and don't automatically get deleted when you close a pull request.
- If the app you are creating a review app version of has any addons, you'll need to make sure you have a credit card on file even if there won't be any charges.
- We are using Heroku forking to copy the apps so be sure to read [this doc](https://devcenter.heroku.com/articles/fork-app) to be aware of the implications.
- Changes you make will not automatically be deployed to the review app. You will need to manually push any changes you make after the review app creation to that review app.
- It is still up to you to push your changes to github, etc.
- Git remotes are created with each review app so you may want to remove the oldgit remotes on occassions.

## TODO

- Add to homebrew for better installation.
- Create an automation script for creating new projects.
- Better output using something like [this](https://github.com/fatih/color)