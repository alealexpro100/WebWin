WebWin
======

### What is it?
This is admin panel for windows server core. It doesn't consume much RAM and CPU (Only 10 MB!). If you don't have a real lack of RAM, please you Microsoft's Windows Admin Center.

### Architecture:
* There are two parts: client and server. 
* Client is written in TypeScript using React and AntD.
* Server is written in GoLang using Gin and powershell.
* Server uses "plugins" as its main functionality.

### Notes
* This is proof-of-concept project. It tries to demonstrate that go + powershell are effective combination for web applications.
* Please, do not use it in production. It is education project.
* For now, it is _NOT_ intened to be used in real projects.