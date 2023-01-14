package eshttp

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/gosom/kit/es"
	"github.com/gosom/kit/lib"
	"github.com/gosom/kit/web"
)

func RegisterDomainRoutes(domain string, mux web.Router, store es.EventStore, registry *es.Registry, aggFactory es.AggregateFactory) {
	handler := NewDomainHandler(domain, store, registry, aggFactory)
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
}

func NewDomainHandler(domain string, store es.EventStore, registry *es.Registry, aggFactory es.AggregateFactory) *DomainHandler {
	return &DomainHandler{
		domain:     domain,
		store:      store,
		registry:   registry,
		aggFactory: aggFactory,
	}
}

type PostCommandResponse struct {
	ID string `json:"id"`
}

func (a *DomainHandler) PostCommand(w http.ResponseWriter, r *http.Request) {
	command, err := es.ParseCommandRequest(a.registry, io.Reader(r.Body))
	if err != nil {
		if errors.Is(err, es.ErrInvalidCommand) {
			web.JSONError(w, r, lib.ErrBadRequest)
			return
		}
		web.JSONError(w, r, err)
		return
	}
	cr, err := es.CommandToCommandRecord(a.domain, command)
	if err != nil {
		web.JSONError(w, r, err)
		return
	}
	commandID, err := a.store.SaveCommandRecords(r.Context(), cr)
	if err != nil {
		web.JSONError(w, r, err)
		return
	}
	if len(commandID) == 0 {
		web.JSONError(w, r, lib.ErrInternal)
		return
	}
	web.JSON(w, r, http.StatusOK, PostCommandResponse{ID: commandID[0]})
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
		web.JSONError(w, r, lib.ErrBadRequest)
		return
	}
	command, err := a.store.GetCommand(r.Context(), commandId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			web.JSONError(w, r, lib.ErrNotFound)
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
		web.JSONError(w, r, lib.ErrBadRequest)
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
		web.JSONError(w, r, lib.ErrBadRequest)
		return
	}
	records, err := a.store.LoadEvents(r.Context(), aggregateId)
	if err != nil {
		web.JSONError(w, r, err)
		return
	}
	if len(records) == 0 {
		web.JSONError(w, r, lib.ErrNotFound)
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
