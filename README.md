#### Usage
```
	

    logkit.Init(FIlE, "test", logkit.LevelDebug)
    defer logkit.Exit()
    
    logkit.Info("this is a log")
    
    logkit.Infof("this is log %s", "arg")


```
