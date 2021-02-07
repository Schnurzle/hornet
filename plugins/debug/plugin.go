package debug

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	"go.uber.org/dig"

	"github.com/gohornet/hornet/pkg/model/storage"
	"github.com/gohornet/hornet/pkg/model/utxo"
	"github.com/gohornet/hornet/pkg/node"
	"github.com/gohornet/hornet/pkg/protocol/gossip"
	"github.com/gohornet/hornet/pkg/restapi"
	"github.com/gohornet/hornet/pkg/tangle"
)

const (
	// ParameterMessageID is used to identify a message by it's ID.
	ParameterMessageID = "messageID"

	// ParameterMilestoneIndex is used to identify a milestone.
	ParameterMilestoneIndex = "milestoneIndex"
)

const (
	// RouteDebugSolidifier is the debug route to manually trigger the solidifier.
	// GET triggers the solidifier.
	RouteDebugSolidifier = "/solidifier"

	// RouteDebugOutputs is the debug route for getting all output IDs.
	// GET returns the outputIDs for all outputs.
	RouteDebugOutputs = "/outputs"

	// RouteDebugOutputsUnspent is the debug route for getting all unspent output IDs.
	// GET returns the outputIDs for all unspent outputs.
	RouteDebugOutputsUnspent = "/outputs/unspent"

	// RouteDebugOutputsSpent is the debug route for getting all spent output IDs.
	// GET returns the outputIDs for all spent outputs.
	RouteDebugOutputsSpent = "/outputs/spent"

	// RouteDebugAddresses is the debug route for getting all known addresses.
	// GET returns all known addresses encoded in hex.
	RouteDebugAddresses = "/addresses"

	// RouteDebugAddressesEd25519 is the debug route for getting all known ed25519 addresses.
	// GET returns all known ed25519 addresses encoded in hex.
	RouteDebugAddressesEd25519 = "/addresses/ed25519"

	// RouteDebugMilestoneDiffs is the debug route for getting a milestone diff by it's milestoneIndex.
	// GET returns the utxo diff (new outputs & spents) for the milestone index.
	RouteDebugMilestoneDiffs = "/ms-diff/:" + ParameterMilestoneIndex

	// RouteDebugRequests is the debug route for getting all pending requests.
	// GET returns a list of all pending requests.
	RouteDebugRequests = "/requests"

	// RouteDebugMessageCone is the debug route for traversing a cone of a message.
	// it traverses the parents of a message until they reference an older milestone than the start message.
	// GET returns the path of this traversal and the "entry points".
	RouteDebugMessageCone = "/message-cones/:" + ParameterMessageID
)

func init() {
	Plugin = &node.Plugin{
		Status: node.Disabled,
		Pluggable: node.Pluggable{
			Name:      "Debug",
			DepsFunc:  func(cDeps dependencies) { deps = cDeps },
			Configure: configure,
		},
	}
}

var (
	Plugin *node.Plugin

	// ErrNodeNotSync is returned when the node was not synced.
	ErrNodeNotSync = errors.New("node not synced")

	deps dependencies
)

type dependencies struct {
	dig.In
	Storage      *storage.Storage
	Tangle       *tangle.Tangle
	RequestQueue gossip.RequestQueue
	UTXO         *utxo.Manager
	Echo         *echo.Echo
}

func configure() {
	routeGroup := deps.Echo.Group("/api/plugins/debug")

	routeGroup.GET(RouteDebugSolidifier, func(c echo.Context) error {
		deps.Tangle.TriggerSolidifier()

		return restapi.JSONResponse(c, http.StatusOK, "solidifier triggered")
	})

	routeGroup.GET(RouteDebugOutputs, func(c echo.Context) error {
		resp, err := debugOutputsIDs(c)
		if err != nil {
			return err
		}

		return restapi.JSONResponse(c, http.StatusOK, resp)
	})

	routeGroup.GET(RouteDebugOutputsUnspent, func(c echo.Context) error {
		resp, err := debugUnspentOutputsIDs(c)
		if err != nil {
			return err
		}

		return restapi.JSONResponse(c, http.StatusOK, resp)
	})

	routeGroup.GET(RouteDebugOutputsSpent, func(c echo.Context) error {
		resp, err := debugSpentOutputsIDs(c)
		if err != nil {
			return err
		}

		return restapi.JSONResponse(c, http.StatusOK, resp)
	})

	routeGroup.GET(RouteDebugAddresses, func(c echo.Context) error {
		resp, err := debugAddresses(c)
		if err != nil {
			return err
		}

		return restapi.JSONResponse(c, http.StatusOK, resp)
	})

	routeGroup.GET(RouteDebugAddressesEd25519, func(c echo.Context) error {
		resp, err := debugAddressesEd25519(c)
		if err != nil {
			return err
		}

		return restapi.JSONResponse(c, http.StatusOK, resp)
	})

	routeGroup.GET(RouteDebugMilestoneDiffs, func(c echo.Context) error {
		resp, err := debugMilestoneDiff(c)
		if err != nil {
			return err
		}

		return restapi.JSONResponse(c, http.StatusOK, resp)
	})

	routeGroup.GET(RouteDebugRequests, func(c echo.Context) error {
		resp, err := debugRequests(c)
		if err != nil {
			return err
		}

		return restapi.JSONResponse(c, http.StatusOK, resp)
	})

	routeGroup.GET(RouteDebugMessageCone, func(c echo.Context) error {
		resp, err := debugMessageCone(c)
		if err != nil {
			return err
		}

		return restapi.JSONResponse(c, http.StatusOK, resp)
	})
}
