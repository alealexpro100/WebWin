# WebWin

## What is it?

This is Web Admin Panel for Windows Server Core.
The aim of this project is to create lightweight Web Admin Control and make it easy-to-use.
If you don't have a real lack of RAM (only 10MB!), please use Microsoft's Windows Admin Center.

## Architecture

* There are two parts: client and server.
* Client is written in TypeScript using Materialize and JQuery.
* Server is written in GoLang using Echo and PowerShell.
* There is "plugins" architecture, that allows easily add more functionality.

## Notes

* This is proof-of-concept project. It tries to demonstrate that go + powershell are effective combination for web applications.
* Please, do not use it in production. It is education project.
* For now, it is _NOT_ intended to be used in real projects.
