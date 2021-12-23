gowin
=====

Provide a simple Windows OS interface to manipulate Windows registry, environment variables, default paths and Windows services with Golang. 

### How to use it
```go
go get github.com/luisiturrios/gowin
```

```go
import "github.com/luisiturrios/gowin"
```

### Example write, read & removing windows registry keys
```go
//Write string on the registry require admin privileges
err = gowin.WriteStringReg("HKLM",`Software\iturrios\gowin`,"value","Hello world")
if err != nil {
  log.Println(err)
} else {
  fmt.Println("Key inserted")
}
```
```go
//Write uint32 on the registry require admin privileges
err = gowin.WriteDwordReg("HKLM",`Software\iturrios\gowin`,"value2", 4294967295)
if err != nil {
  log.Println(err)
} else{
  fmt.Println("Key inserted")
}
```
```go
//Get reg
val, err := gowin.GetReg("HKLM", `Software\Microsoft\Windows\CurrentVersion\Explorer\Shell Folders`, "Common AppData")
if err != nil {
  log.Println(err)
}
fmt.Printf(val)
```
```go
//Remove key
err = gowin.DeleteKey("HKLM",`Software\iturrios`,"gowin")
if err != nil {
  log.Println(err)
} else {
  fmt.Println("Key Removed")
}
```

### Example reading windows ShellFolders
```go
folders := gowin.ShellFolders{Context: gowin.ALL}
//Or 
folder := new(gowin.ShellFolders)

//Read ProgramFiles
fmt.Println(folders.ProgramFiles())

//Read all user AppData
folders.Context = gowin.ALL
fmt.Println(folders.AppData())

//Read Current user AppData
folders.Context = gowin.USER
fmt.Println(folders.AppData())

//Functions
folders.ProgramFiles()
folders.AppData()
folders.Desktop()
folders.Documents()
folders.StartMenu()
folders.StartMenuPrograms()
```

### Example reading windows environment variables
```go
//Get environment var
goroot := gowin.GetEnvVar("GOROOT")
fmt.Printf("GORROT: %s\n", goroot)
```
```go
//Write environment var
err := gowin.WriteEnvVar("TVAR","hello word")
if err != nil {
  log.Println(err)
}
```

### Donation
If you appreciate the work in this repo and like the continue development donate to this paypal account 
luisiturrios@me.com
Thanks
