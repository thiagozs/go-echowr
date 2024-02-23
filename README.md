# Echowr - Wrapper for Echo server

```go
// healthcheck
 sys := server.NewRouters()
 sys.AddRouter("/healthcheck", 
     server.Methods{http.MethodGet: func(c server.Context) error {
         return c.String(http.StatusOK, "test passed")
     },
 })
 
 // register the routes
 _ = a.RegisterRouters(server.ROOT, sys)

  sigChan := make(chan os.Signal, 1)
  signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

 // start the server
 a.Start()

 sig := <-sigChan
 fmt.Printf("\nSignal %s received, shutting down\n", sig.String())

// shutdown
 _ = a.GracefulShutdown()

 os.Exit(0)
```

-----

## Versioning and license

Our version numbers follow the [semantic versioning specification](http://semver.org/). You can see the available versions by checking the [tags on this repository](https://github.com/thiagozs/go-echowr/tags). For more details about our license model, please take a look at the [LICENSE](LICENSE.md) file.

**2024**, thiagozs.
