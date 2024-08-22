package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/mroth/weightedrand"
)

type runResponse struct {
	Success     bool     `json:"success"`
	TracedNodes []string `json:"traced_nodes,omitempty"`
}

func handleRunRequest(c *fiber.Ctx) error {
	ctx, span := tracer.Start(c.UserContext(), "Handle run request")
	defer span.End()

	if len(NEXT_NODE) < 1 {
		c.JSON(fiber.Map{
			"success": true,
			"traced_nodes": []string{
				NODE_NAME,
			},
		})
		return nil
	}

	chosenNextNode := chooseNextNode(ctx)
	nodeResponse := requestToNode(ctx, chosenNextNode)
	parsedResponse := parseNodeResponse(ctx, nodeResponse)

	if !parsedResponse.Success || parsedResponse.TracedNodes == nil {
		c.JSON(fiber.Map{
			"success": false,
		})
		return nil
	}

	parsedResponse.TracedNodes = append(
		parsedResponse.TracedNodes,
		NODE_NAME,
	)

	c.JSON(parsedResponse)
	return nil
}

func chooseNextNode(ctx context.Context) string {
	_, span := tracer.Start(ctx, "Choose next node from list - "+NODE_NAME)
	defer span.End()

	nextNodes := strings.Split(NEXT_NODE, ",")
	choices := []weightedrand.Choice{}

	for index, nextNode := range nextNodes {
		choices = append(choices, weightedrand.Choice{
			Item:   nextNode,
			Weight: uint(len(nextNodes) - index),
		})
	}

	chooser, _ := weightedrand.NewChooser(choices...)
	return chooser.Pick().(string)
}

func requestToNode(ctx context.Context, node string) []byte {
	c, span := tracer.Start(ctx, "HTTP request to next node")
	defer span.End()

	nodeUrl := fmt.Sprintf("http://%s/api/v1/run", node)
	req, err := http.NewRequestWithContext(c, "GET", nodeUrl, nil)
	if err != nil {
		return []byte("{\"success\":false}")
	}

	resp, err := client.Do(req)
	if err != nil {
		return []byte("{\"success\":false}")
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []byte("{\"success\":false}")
	}

	return body[:]
}

func parseNodeResponse(ctx context.Context, nodeResponse []byte) (response runResponse) {
	_, span := tracer.Start(ctx, "Parsing response from node")
	defer span.End()

	if err := json.Unmarshal([]byte(nodeResponse), &response); err != nil {
		return runResponse{
			Success: false,
		}
	}

	return
}
