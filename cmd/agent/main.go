package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"quickce/pkg/agent"
)

func main() {
	server := flag.String("server", "http://localhost:8080", "Central server URL")
	name := flag.String("name", "", "Agent name (e.g. beijing-01)")
	location := flag.String("location", "", "Agent location (e.g. Beijing)")
	agentID := flag.String("agent-id", "", "Agent ID (omit for auto-generate)")
	flag.Parse()

	if *name == "" {
		*name = os.Getenv("AGENT_NAME")
	}
	if *location == "" {
		*location = os.Getenv("AGENT_LOCATION")
	}
	if *agentID == "" {
		*agentID = os.Getenv("AGENT_ID")
	}
	if *name == "" {
		log.Fatal("agent name is required (use -name flag or AGENT_NAME env)")
	}
	if *location == "" {
		log.Fatal("agent location is required (use -location flag or AGENT_LOCATION env)")
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	client := agent.New(*server, *name, *location, *agentID)

	if err := client.Register(ctx); err != nil {
		log.Fatalf("Agent registration failed: %v", err)
	}

	log.Printf("Agent %s (%s) started, ID=%s, polling %s", *name, *location, client.ID, *server)
	client.Run(ctx)
	log.Println("Agent stopped")
}
