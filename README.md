# WebWin

## What is it?

This is Web Admin Panel for Windows Server Core.
The aim of this project is to create lightweight Web Admin Control and make it easy-to-use.
If you don't have a real lack of RAM (only 10MB!), please use Microsoft's Windows Admin Center.

## Architecture

* There are two parts: client and server.
* Client is written using [Materialize](https://github.com/materializecss/materialize), [Material icons](https://fonts.google.com/icons?selected=Material+Icons) and JQuery.
* Server is written using Golang and [Echo](https://echo.labstack.com).
* There is "plugins" architecture, that allows easily add more functionality.

## Plugin

* To see example of building plugin see `001_test` plugin in `web_root\plugins` directory.
* Plugins use PowerShell to execute server side's jobs. Also, you can use only `.bat` file without PowerShell.

## API

* `/api/plugins?action=list` - Get json file with list of plugins. NOTE: Plugin list is updated only on startup.
* `/api/plugins?action=get_status&id=your_id` - Get status of job with `your_id` id.
* `/api/plugins?action=jobs_clear` - Clear completed jobs. Sends `bad requested` if any job is not completed.
* `/api/plugins/your_plugin?action=add&param=your_param` - Start `your_plugin` with `your_param` parameter.
* `/api/internal?action=set_auth&user=your_user&pass_hash=your_pass_hash` - Set defined user and pass_hash.


## Notes

* This is proof-of-concept project. It tries to demonstrate that Go + PowerShell are effective combination for web applications.
* Please, do _not_ use it in production. It is education project.
