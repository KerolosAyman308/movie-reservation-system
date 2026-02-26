to enable debug i needed to 
``go env -w GOARCH=amd64`` to change the arch and install dlve
go to command platte in vs code and search for go install/update and choose 


in the lunch json 
"program": "${workspaceFolder}/cmd",
            "cwd": "${workspaceFolder}"


in setting.json

"go.alternateTools": {
        "dlv": "C:\\Users\\User\\go\\bin\\windows_amd64\\dlv.exe"
    }

-----------------------
to use hot reload integrate dlv as in the .air.toml
ant edit the lunch.json:
```
{
    "name": "Launch Package",
    "type": "go",
    "request": "attach",
    "mode": "remote",
    "debugAdapter": "dlv-dap",
    "remotePath": "${workspaceFolder}",
    "port": 2345,
    "host": "127.0.0.1"
}
```

f5 to attach the debugger in vscode t