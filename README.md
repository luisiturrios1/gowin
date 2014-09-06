gowin
=====

Provide simple Windows OS interface to manipulate windows registry, environment variables, default paths and windows services from Golang lenguaje

###How to use it
```
	import "github.com/luisiturrios/gowin"
```


###Example Read windows registry key 
```
	val, err := gowin.GetRegKey("HKLM", `Software\Microsoft\Windows\CurrentVersion\Explorer\Shell Folders`, "Common AppData")
	if err != nil {
		log.Println(err)
	}
	fmt.Printf(val)
```
###Example Read windows ShellFolders
```
	folders := gowin.ShellFolders{gowin.ALL}
	//Or default context is USER
	folder := new(gowin.ShellFolders)

	//Read ProgramFiles
	fmt.Println(folders.ProgramFiles())
	
	//Read all user AppData
	folders.Context = gowin.ALL
	fmt.Println(folders.AppData())
	
	//Read Current user AppData
	folders.Context = gowin.USER
	fmt.Println(folders.AppData())
```

###Donation

If you appreciate the work in this repo and like the continue development donate to this paypal account 
luisiturrios@me.com
Thanks
