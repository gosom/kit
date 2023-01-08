package eshttp

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gosom/kit/core"
	"github.com/gosom/kit/es"
	"github.com/gosom/kit/web"
)

func RegisterDomainRoutes(domain string, mux web.Router, store es.EventStore, registry *es.Registry, aggFactory es.AggregateFactory, dispatcher es.CommandDispatcher) {
	handler := NewDomainHandler(domain, store, registry, aggFactory, dispatcher)
	mux.MethodFunc(http.MethodGet, fmt.Sprintf("/%s/commands/{commandId}", handler.domain), handler.GetCommand)
	mux.MethodFunc(http.MethodPost, fmt.Sprintf("/%s/commands", handler.domain), handler.PostCommand)
	mux.MethodFunc(http.MethodGet, fmt.Sprintf("/%s/events/{aggregateId}", handler.domain), handler.GetEvents)
	mux.MethodFunc(http.MethodGet, fmt.Sprintf("/%s/aggregates/{aggregateId}", handler.domain), handler.GetAggregate)
}

type DomainHandler struct {
	domain     string
	store      es.EventStore
	registry   *es.Registry
	aggFactory es.AggregateFactory
	dispatcher es.CommandDispatcher
}

func NewDomainHandler(domain string, store es.EventStore, registry *es.Registry, aggFactory es.AggregateFactory,
	dispatcher es.CommandDispatcher) *DomainHandler {
	return &DomainHandler{
		domain:     domain,
		store:      store,
		registry:   registry,
		aggFactory: aggFactory,
		dispatcher: dispatcher,
	}
}

type PostCommandRequest struct {
	Name    string          `json:"name"`
	Payload json.RawMessage `json:"payload"`
}

type PostCommandResponse struct {
	ID string `json:"id"`
}

func (u *DomainHandler) PostCommand(w http.ResponseWriter, r *http.Request) {
	var req PostCommandRequest
	if err := web.DecodeBody(r, &req); err != nil {
		web.JSONError(w, r, err)
		return
	}
	conv, ok := u.registry.GetCommand(req.Name)
	if !ok {
		fmt.Println("command not found")
		web.JSONError(w, r, core.ErrBadRequest)
		return
	}
	cr := es.CommandRecord{
		RecordBase: es.RecordBase{
			Data:      req.Payload,
			EventType: req.Name,
		},
	}
	cmd, err := conv(cr)
	if err != nil {
		web.JSONError(w, r, err)
		return
	}
	commandID, err := u.dispatcher.Dispatch(r.Context(), cmd)
	if err != nil {
		if errors.Is(err, es.ErrInvalidCommand) {
			web.JSONError(w, r, core.ErrBadRequest)
			return
		}
		web.JSONError(w, r, err)
		return
	}
	web.JSON(w, r, http.StatusOK, PostCommandResponse{ID: commandID})
}

type GetCommandResponse es.CommandRecord

func (u GetCommandResponse) MarshalJSON() ([]byte, error) {
	type Alias GetCommandResponse
	m := make(map[string]interface{})
	if err := json.Unmarshal(u.Data, &m); err != nil {
		return nil, err
	}

	return json.Marshal(&struct {
		Alias
		Data map[string]interface{} `json:"Data"`
	}{
		Data:  m,
		Alias: (Alias)(u),
	})
}

func (a *DomainHandler) GetCommand(w http.ResponseWriter, r *http.Request) {
	commandId := web.StringURLParam(r, "commandId")
	if len(commandId) == 0 {
		web.JSONError(w, r, core.ErrBadRequest)
		return
	}
	command, err := a.store.GetCommand(r.Context(), commandId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			web.JSONError(w, r, core.ErrNotFound)
			return
		}
		web.JSONError(w, r, err)
		return
	}
	web.JSON(w, r, http.StatusOK, GetCommandResponse(command))
}

type GetEventResponse es.EventRecord

func (u GetEventResponse) MarshalJSON() ([]byte, error) {
	type Alias GetEventResponse
	m := make(map[string]interface{})
	if err := json.Unmarshal(u.Data, &m); err != nil {
		return nil, err
	}

	return json.Marshal(&struct {
		Alias
		Data map[string]interface{} `json:"Data"`
	}{
		Data:  m,
		Alias: (Alias)(u),
	})
}

func (a *DomainHandler) GetEvents(w http.ResponseWriter, r *http.Request) {
	aggregateId := web.StringURLParam(r, "aggregateId")
	if len(aggregateId) == 0 {
		web.JSONError(w, r, core.ErrBadRequest)
		return
	}
	events, err := a.store.LoadEvents(r.Context(), aggregateId)
	if err != nil {
		web.JSONError(w, r, err)
		return
	}
	items := make([]GetEventResponse, len(events))
	for i := range events {
		items[i] = GetEventResponse(events[i])
	}
	web.JSON(w, r, http.StatusOK, items)
}

func (a *DomainHandler) GetAggregate(w http.ResponseWriter, r *http.Request) {
	aggregateId := web.StringURLParam(r, "aggregateId")
	if len(aggregateId) == 0 {
		web.JSONError(w, r, core.ErrBadRequest)
		return
	}
	records, err := a.store.LoadEvents(r.Context(), aggregateId)
	if err != nil {
		web.JSONError(w, r, err)
		return
	}
	if len(records) == 0 {
		web.JSONError(w, r, core.ErrNotFound)
		return
	}
	events, err := es.EventRecordsToEvents(a.registry, records)
	if err != nil {
		web.JSONError(w, r, err)
		return
	}
	agg, err := a.aggFactory()
	if err != nil {
		web.JSONError(w, r, err)
		return
	}
	if err := es.Load(agg, events); err != nil {
		web.JSONError(w, r, err)
		return
	}
	web.JSON(w, r, http.StatusOK, agg)
}
