[![wercker status](https://app.wercker.com/status/d25f1a79c7ec302ba836fde68da07e21/m/ "wercker status")](https://app.wercker.com/project/byKey/d25f1a79c7ec302ba836fde68da07e21)

# What's this?

This is the skeleton of a Web API written in [Golang](https://golang.org). In order to use it, you have to clone it, rename it as you want (also rename all the imports from 'gost' to your app's name) and then start coding over this template.

This template contains basic endpoints for Users (+ login system) and Transactions (payments made between users). Both the endpoints are fully working ones, however the user is free to modify/delete them as they will. 

# Configuration steps for the API

1. Install Go and set up your [GOPATH](http://golang.org/doc/code.html#GOPATH). Starting with version Go1.4, you also need to set the *GOROOT_BOOTSTRAP* variable, to the same path as your *GOROOT*.

2. Install [MongoDb](https://scotch.io/tutorials/an-introduction-to-mongodb#installation-and-running-mongodb)

3. Create a database named __serverName_db__ and then create an user for the database using the following command in **mongodb shell**:
	`db.createUser( { user: "serverNameAdmin", pwd: "serverNamePass", roles: [ { role: "readWrite", db: "serverName_db" } ] } )`

4. Install all the necessary dependencies using the following command in the **cmd/console/terminal**:
	`go get -v`

5. For testing purposes, create another database named __serverName_db_test__, but don't create a user for it like for the main database.
In order for the tests to run, you need to set the following environment variables correctly:

If you don't want to use the terminal for creating the databases, you can always use [Robomongo](http://robomongo.org), a very easy to use GUI for mongodb.

	**MONGODB_URL** = connection_string_for_mongodb (i.e. 'mongodb://localhost:27017')
	**GOST_TESTAPP_DB_NAME** = serverName_test_app_db_name
	**GOST_TESTAPP_DB_CONN** = $MONGODB_URL/$GOST_TESTAPP_DB_NAME
	**GOST_TESTAPP_NAME** = serverName_test_app_name
	**GOST_TESTAPP_INSTANCE** = /gost-test/ (access path, such as: **/some_domain/gost-test/**some_link_path)
	**GOST_TESTAPP_HTTP** = serverName_testapp_http_server (i.e. :7500 for localhost:7500/; use 0.0.0.0:7500 for access from outside the local domain)

### Unix bash example
	export MONGODB_URL="localhost"
    export GOST_TESTAPP_DB_NAME="gost_test_app_db"
    export GOST_TESTAPP_DB_CONN="$MONGODB_URL/$GOST_TESTAPP_DB_NAME"
    export GOST_TESTAPP_NAME="gost_test_app"
    export GOST_TESTAPP_INSTANCE="/gost-test/"
    export GOST_TESTAPP_HTTP="localhost:7500"

<br/>

### Windows Powershell example
	[Environment]::SetEnvironmentVariable("MONGODB_URL", "localhost", "User")
	[Environment]::SetEnvironmentVariable("GOST_TESTAPP_DB_NAME", "gost_test_app_db", "User")
	[Environment]::SetEnvironmentVariable("GOST_TESTAPP_DB_CONN", ("{0}/{1}" -f [Environment]::GetEnvironmentVariable("MONGODB_URL", "User"),[Environment]::GetEnvironmentVariable	("GOST_TESTAPP_DB_NAME", "User")), "User")
	[Environment]::SetEnvironmentVariable("GOST_TESTAPP_NAME", "gost_test_app", "User")
	[Environment]::SetEnvironmentVariable("GOST_TESTAPP_INSTANCE", "/gost-test/", "User")
	[Environment]::SetEnvironmentVariable("GOST_TESTAPP_HTTP ", "localhost:7500", "User")

# Easily renaming the template

To rename this template from __gost__ to your new application name that will be built over this, use: 
- For **UNIX**: the bash script named __app-rename__ to easily do all the renaming (folders, files and imports). Use the --help flag for usage details"
- For **Windows**: the executable application named __AppRename.exe__
