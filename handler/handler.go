//go:generate go run middleware HandlerMiddleware:Stack
package handler

import (
	"fmt"
	"net/http"

	middleware2 "github.plaid.com/plaid/typed-middleware"
	"github.plaid.com/plaid/typed-middleware/appmiddleware"
)


// app programmer defines their needs
type HandlerMiddleware interface {
	// references each middleware they require here
	appmiddleware.ClientIDFromRequest
}

type getUserHandler struct {
	stack Stack
}

func NewUserHandler(
	authMiddleware appmiddleware.AuthenticationFromRequestMiddleware,
) getUserHandler {
	// This part will be dependency injected as per normal - currently via app context
	// Since clientID middleware has no constructor
	h := getUserHandler{}
	h.stack = NewStack(
		// currently push instantiation up here, leaving the generator
		// free of that complexity
		appmiddleware.ClientIDFromRequestMiddleware{},
		authMiddleware,
	)
	return h
}

func (h *getUserHandler) Handle(res http.ResponseWriter, req http.Request) {
	result, override := h.stack.Run(req)
	// the stack value could come from ctx for now, or be replaced by a mock
	if override != nil {
		// or can explicitly check out what's happened: an error, or a result spec
		middleware2.Respond(override, res)
		return
	}

	cid := result.ClientID()
	fmt.Fprintf(res, "Got client ID %d", cid)
}


