package hostredirect

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	"github.com/robinbraemer/event"
	"github.com/spf13/viper"
	"go.minekube.com/common/minecraft/component"
	"go.minekube.com/gate/pkg/edition/java/proxy"
)

// Config holds the server mappings
type Config struct {
	ServerMappings map[string]string `mapstructure:"serverMappings"`
}

// Plugin is a demo plugin that redirects players based on the host they connect with.
var Plugin = proxy.Plugin{
	Name: "HostRedirect",
	Init: func(ctx context.Context, p *proxy.Proxy) error {
		// Get the logger for this plugin.
		log := logr.FromContextOrDiscard(ctx)
		log.Info("Hello from HostRedirect plugin!")

		// Initialize Viper configuration
		viper.SetConfigName("mapping")
		viper.SetConfigType("yaml")
		viper.AddConfigPath("./")

		// Set default values
		viper.SetDefault("serverMappings", map[string]string{
			"examplehost1.com": "server1",
			"examplehost2.com": "server2",
		})

		// Read the config file
		if err := viper.ReadInConfig(); err != nil {
			// Create a config file if it doesn't exist
			log.Info("Couldn't find mapping.yml, creating a new config file")
			err = viper.WriteConfigAs("./mapping.yml")
			if err != nil {
				return fmt.Errorf("error creating config file: %w", err)
			}
		}

		// Manually unmarshal the config
		var config Config
		if err := viper.Unmarshal(&config); err != nil {
			return fmt.Errorf("failed to unmarshal config: %w", err)
		}

		// Log the server mappings for debugging purposes
		log.Info("Server Mappings", "serverMappings", config.ServerMappings)

		// Register some event handlers.
		event.Subscribe(p.Event(), 0, onPlayerChooseInitialServer(p, log, config.ServerMappings))

		return nil
	},
}

// onPlayerChooseInitialServer handles the PlayerChooseInitialServerEvent to redirect players.
func onPlayerChooseInitialServer(p *proxy.Proxy, log logr.Logger, serverMappings map[string]string) func(*proxy.PlayerChooseInitialServerEvent) {
	return func(e *proxy.PlayerChooseInitialServerEvent) {
		// Get the player's connecting host.
		conn := e.Player().VirtualHost()
		host := conn.String()
		host = strings.Split(host, ":")[0]

		// Check if we have a server mapping for this host.
		serverName, ok := serverMappings[host]
		if !ok {
			msg := fmt.Sprintf("No server mapping for host: %s", host)
			log.Info(msg)
			return
		}

		// Get the registered server by name.
		var server proxy.RegisteredServer
		for _, s := range p.Servers() {
			if s.ServerInfo().Name() == serverName {
				server = s
				break
			}
		}

		if server == nil {
			msg := fmt.Sprintf("Server not found: %s", serverName)
			log.Info(msg)
			e.Player().Disconnect(&component.Text{Content: msg})
			return
		}

		// Redirect the player to the appropriate server.
		e.SetInitialServer(server)
		log.Info("Redirecting player", "username", e.Player().Username(), "server", serverName)
	}
}
