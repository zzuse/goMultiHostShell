# golang Multi-Host-Shell
#What
If you want to connect to multiple hosts and run the same shell command or maybe you want to run simultaneously, then this tool will suit you.
#Why
My work place have many linux hosts what we always login to, we want to monitor processes or run same command, but then I found login to multiple hosts and start a login shell is very handy.
And sometimes I want to run some commands at almost the same time on multiple machines, such as doing stress tests. Then I created this tool.
#How this work
This tool is a golang program. It simulates a ssh client and login into other hosts. If on other hosts, the specified command or shell exists, then it will run that command and return the output to your stdout.
Cause the command running is simultaneously, so that your stdout maybe full of mixed output which is from different hosts and the running sequence is unpredictable.
#How to build
```shell
#on mymac
go build multiHostShell.go
#or cross compile 
GOOS=linux GOARCH=amd64 go build -o multiHostShell multiHostShell.go
GOOS=windows GOARCH=amd64 go build -o multiHostShell.exe multiHostShell.go  
GOOS=windows GOARCH=386 go build -o multiHostShell.exe multiHostShell.go 
```
#How to run
First you need to edit the specified configuration file ---  hosts.json, you can see the example hosts.json.
In the hosts.json you need to input **hostip,user,pass**. Now this version do not support public key accessing and configuration yet, maybe later.
```python
multissh -c "/home/work/run.sh" -u ./hosts.json      
```
#Note
The shell here ---- **run.sh**, must include the environment variable setup part,  or fail.
If you have any questions, please send messages me
#License
BSD




