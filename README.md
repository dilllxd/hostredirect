# Host Redirect for Gate

**Host Redirect** is a plugin for the [Gate](https://gate.minekube.com/) proxy that allows you to redirect users based on the host they connect to.

## Getting Started

### 1. Add the Package

Add the `hostredirect` package to your Gate project:

```bash
go get github.com/dilllxd/hostredirect
```

### 2. Register the Plugin

Include the plugin in your `main()` function:

```go
func main() {
    proxy.Plugins = append(proxy.Plugins,
        // your plugins
        hostredirect.Plugin,
    )
    gate.Execute()
}
```

### 3. Configure the Plugin

After starting your server, a new file named `mapping.yml` will be created. Configure it with your host-to-server mappings:

```yaml
servermappings:
    examplehost1.com: server1
    examplehost2.com: server2
```

In this example, connections to `examplehost1.com` will be redirected to `server1`, and connections to `examplehost2.com` will be redirected to `server2`.

---

Feel free to reach out if you have any questions or need further assistance with setting up the Host Redirect plugin for Gate.
