package handlers

import (
	"context"
	"fmt"

	"github.com/danielgtaylor/huma/v2"

	"github.com/syhily/bookworm/handlers/register"
)

func init() {
	huma.Get(register.API, "/greeting/{name}", greeting)
}

// GreetingOutput represents the greeting operation response.
type GreetingOutput struct {
	Body struct {
		Message string `json:"message" example:"Hello, world!" doc:"Greeting message"`
	}
}

func greeting(ctx context.Context, input *struct {
	Name string `path:"name" maxLength:"30" example:"world" doc:"Name to greet"`
}) (*GreetingOutput, error) {
	resp := &GreetingOutput{}
	resp.Body.Message = fmt.Sprintf("Hello, %s!", input.Name)
	return resp, nil
}
