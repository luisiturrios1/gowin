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