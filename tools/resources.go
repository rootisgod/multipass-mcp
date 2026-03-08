package tools

import (
	"context"
	"encoding/json"
	"regexp"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func RegisterResources(s *server.MCPServer) {
	// Static resources
	s.AddResource(
		mcp.NewResource("multipass://instances", "Instances",
			mcp.WithResourceDescription("List all Multipass instances with state, IPs, image, and resource usage."),
			mcp.WithMIMEType("application/json"),
		),
		handleListInstances,
	)

	s.AddResource(
		mcp.NewResource("multipass://images", "Images",
			mcp.WithResourceDescription("List available Multipass images that can be launched."),
			mcp.WithMIMEType("application/json"),
		),
		handleListImages,
	)

	s.AddResource(
		mcp.NewResource("multipass://networks", "Networks",
			mcp.WithResourceDescription("List host network devices available for Multipass instances."),
			mcp.WithMIMEType("application/json"),
		),
		handleListNetworks,
	)

	s.AddResource(
		mcp.NewResource("multipass://version", "Version",
			mcp.WithResourceDescription("Get Multipass version information."),
			mcp.WithMIMEType("application/json"),
		),
		handleGetVersion,
	)

	s.AddResource(
		mcp.NewResource("multipass://aliases", "Aliases",
			mcp.WithResourceDescription("List configured Multipass command aliases."),
			mcp.WithMIMEType("application/json"),
		),
		handleListAliases,
	)

	// Template resource
	s.AddResourceTemplate(
		mcp.NewResourceTemplate("multipass://instance/{name}", "Instance Details",
			mcp.WithTemplateDescription("Get detailed information about a specific Multipass instance."),
			mcp.WithTemplateMIMEType("application/json"),
		),
		handleGetInstance,
	)
}

func jsonResource(ctx context.Context, uri string, command string) ([]mcp.ResourceContents, error) {
	data, err := runMultipassJSON(ctx, defaultTimeout, command)
	if err != nil {
		return nil, err
	}
	pretty, _ := json.MarshalIndent(json.RawMessage(data), "", "  ")
	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      uri,
			MIMEType: "application/json",
			Text:     string(pretty),
		},
	}, nil
}

func handleListInstances(ctx context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return jsonResource(ctx, req.Params.URI, "list")
}

func handleListImages(ctx context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return jsonResource(ctx, req.Params.URI, "find")
}

func handleListNetworks(ctx context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return jsonResource(ctx, req.Params.URI, "networks")
}

func handleGetVersion(ctx context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return jsonResource(ctx, req.Params.URI, "version")
}

func handleListAliases(ctx context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	return jsonResource(ctx, req.Params.URI, "aliases")
}

func handleGetInstance(ctx context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
	re := regexp.MustCompile(`multipass://instance/(.+)`)
	matches := re.FindStringSubmatch(req.Params.URI)
	if len(matches) < 2 {
		return nil, nil
	}
	name := matches[1]

	data, err := runMultipassJSON(ctx, defaultTimeout, "info", name)
	if err != nil {
		return nil, err
	}
	pretty, _ := json.MarshalIndent(json.RawMessage(data), "", "  ")
	return []mcp.ResourceContents{
		mcp.TextResourceContents{
			URI:      req.Params.URI,
			MIMEType: "application/json",
			Text:     string(pretty),
		},
	}, nil
}
