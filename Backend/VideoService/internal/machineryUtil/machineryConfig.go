package machineryutil

import (
    "fmt"
    "net/url"

    "github.com/RichardKnop/machinery/v2"
    "github.com/RichardKnop/machinery/v2/config"
    backendsiface "github.com/RichardKnop/machinery/v2/backends/iface"
    brokersiface "github.com/RichardKnop/machinery/v2/brokers/iface"
    locksiface "github.com/RichardKnop/machinery/v2/locks/iface"
    "github.com/RichardKnop/machinery/v2/brokers/amqp"
    redisbackend "github.com/RichardKnop/machinery/v2/backends/redis"
)

// LoadMachineryConfigFromYaml loads config from machinery.yaml
func LoadMachineryConfigFromYaml() (*config.Config, error) {
    cnf, err := config.NewFromYaml("machinery.yaml", false)
    if err != nil {
        return nil, fmt.Errorf("failed to load machinery.yaml: %w", err)
    }
    return cnf, nil
}

// buildBroker returns a RabbitMQ (AMQP) broker implementation.
func buildBroker(cfg *config.Config) (brokersiface.Broker, error) {
    return amqp.New(cfg), nil
}

// buildBackend returns a Redis‐based result backend.
func buildBackend(cfg *config.Config) (backendsiface.Backend, error) {
    u, err := url.Parse(cfg.ResultBackend)
    if err != nil {
        return nil, fmt.Errorf("invalid ResultBackend URI %q: %w", cfg.ResultBackend, err)
    }
    if u.Scheme != "redis" {
        return nil, fmt.Errorf("ResultBackend URI must use redis:// scheme, got %q", u.Scheme)
    }
    host := u.Host
    var password string
    if u.User != nil {
        if pw, ok := u.User.Password(); ok {
            password = pw
        }
    }
    var db = 0
    socketPath := ""

    return redisbackend.New(cfg, host, password, socketPath, db), nil
}

// buildLock returns nil because we’re not using a separate distributed lock.
func buildLock() (locksiface.Lock, error) {
    return nil, nil
}

// CreateMachineryServer wires everything together and returns *machinery.Server.
func CreateMachineryServer() (*machinery.Server, error) {
    // 1) Load config from YAML
    cfg, err := LoadMachineryConfigFromYaml()
    if err != nil {
        return nil, err
    }

    // 2) Instantiate the broker (RabbitMQ)
    broker, err := buildBroker(cfg)
    if err != nil {
        return nil, fmt.Errorf("cannot create AMQP broker: %w", err)
    }

    // 3) Instantiate the Redis backend
    backend, err := buildBackend(cfg)
    if err != nil {
        return nil, fmt.Errorf("cannot create Redis backend: %w", err)
    }

    // 4) No lock backend in this example
    lock, err := buildLock()
    if err != nil {
        return nil, fmt.Errorf("cannot create lock: %w", err)
    }

    server := machinery.NewServer(cfg, broker, backend, lock)
    return server, nil
}