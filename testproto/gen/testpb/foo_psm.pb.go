// Code generated by protoc-gen-go-psm. DO NOT EDIT.

package testpb

import (
	context "context"
	fmt "fmt"
	psm "github.com/pentops/protostate/psm"
	proto "google.golang.org/protobuf/proto"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

// StateObjectOptions: FooPSM
type FooPSMEventer = psm.Eventer[
	*FooState,
	FooStatus,
	*FooEvent,
	FooPSMEvent,
]

type FooPSM = psm.StateMachine[
	*FooState,
	FooStatus,
	*FooEvent,
	FooPSMEvent,
]

func DefaultFooPSMConfig() *psm.StateMachineConfig[
	*FooState,
	FooStatus,
	*FooEvent,
	FooPSMEvent,
] {
	return psm.NewStateMachineConfig[
		*FooState,
		FooStatus,
		*FooEvent,
		FooPSMEvent,
	](FooPSMConverter{}, DefaultFooPSMTableSpec)
}

func NewFooPSM(db psm.Transactor, config *psm.StateMachineConfig[
	*FooState,
	FooStatus,
	*FooEvent,
	FooPSMEvent,
]) (*FooPSM, error) {
	return psm.NewStateMachine[
		*FooState,
		FooStatus,
		*FooEvent,
		FooPSMEvent,
	](db, config)
}

type FooPSMTableSpec = psm.TableSpec[
	*FooState,
	FooStatus,
	*FooEvent,
	FooPSMEvent,
]

var DefaultFooPSMTableSpec = FooPSMTableSpec{
	StateTable: "foo",
	EventTable: "foo_event",
	PrimaryKey: func(event *FooEvent) (map[string]interface{}, error) {
		return map[string]interface{}{
			"id": event.FooId,
		}, nil
	},
	StateColumns: func(state *FooState) (map[string]interface{}, error) {
		return map[string]interface{}{
			"tenant_id": state.TenantId,
		}, nil
	},
	EventColumns: func(event *FooEvent) (map[string]interface{}, error) {
		metadata := event.Metadata
		return map[string]interface{}{
			"id":        metadata.EventId,
			"timestamp": metadata.Timestamp,
			"actor":     metadata.Actor,
			"foo_id":    event.FooId,
			"tenant_id": event.TenantId,
		}, nil
	},
	EventPrimaryKeyFieldPaths: []string{
		"metadata.event_id",
	},
	StatePrimaryKeyFieldPaths: []string{
		"foo_id",
	},
}

type FooPSMTransitionBaton = psm.TransitionBaton[*FooEvent, FooPSMEvent]

func FooPSMFunc[SE FooPSMEvent](cb func(context.Context, FooPSMTransitionBaton, *FooState, SE) error) psm.TransitionFunc[
	*FooState,
	FooStatus,
	*FooEvent,
	FooPSMEvent,
	SE,
] {
	return psm.TransitionFunc[
		*FooState,
		FooStatus,
		*FooEvent,
		FooPSMEvent,
		SE,
	](cb)
}

type FooPSMEventKey string

const (
	FooPSMEventNil     FooPSMEventKey = "<nil>"
	FooPSMEventCreated FooPSMEventKey = "created"
	FooPSMEventUpdated FooPSMEventKey = "updated"
	FooPSMEventDeleted FooPSMEventKey = "deleted"
)

type FooPSMEvent interface {
	proto.Message
	PSMEventKey() FooPSMEventKey
}
type FooPSMConverter struct{}

func (c FooPSMConverter) Unwrap(e *FooEvent) FooPSMEvent {
	return e.UnwrapPSMEvent()
}

func (c FooPSMConverter) StateLabel(s *FooState) string {
	return s.Status.String()
}

func (c FooPSMConverter) EventLabel(e FooPSMEvent) string {
	return string(e.PSMEventKey())
}

func (c FooPSMConverter) EmptyState(e *FooEvent) *FooState {
	return &FooState{
		FooId:    e.FooId,
		TenantId: e.TenantId,
	}
}

func (c FooPSMConverter) DeriveChainEvent(e *FooEvent, systemActor psm.SystemActor, eventKey string) *FooEvent {
	metadata := &Metadata{
		EventId:   systemActor.NewEventID(e.Metadata.EventId, eventKey),
		Timestamp: timestamppb.Now(),
	}
	actorProto := systemActor.ActorProto()
	refl := metadata.ProtoReflect()
	refl.Set(refl.Descriptor().Fields().ByName("actor"), actorProto)
	return &FooEvent{
		Metadata: metadata,
		FooId:    e.FooId,
		TenantId: e.TenantId,
	}
}

func (c FooPSMConverter) CheckStateKeys(s *FooState, e *FooEvent) error {
	if s.FooId != e.FooId {
		return fmt.Errorf("state field 'FooId' %q does not match event field %q", s.FooId, e.FooId)
	}
	if s.TenantId == nil {
		if e.TenantId != nil {
			return fmt.Errorf("state field 'TenantId' is nil, but event field is not (%q)", *e.TenantId)
		}
	} else if e.TenantId == nil {
		return fmt.Errorf("state field 'TenantId' is not nil (%q), but event field is", *e.TenantId)
	} else if *s.TenantId != *e.TenantId {
		return fmt.Errorf("state field 'TenantId' %q does not match event field %q", *s.TenantId, *e.TenantId)
	}
	return nil
}

func (ee *FooEventType) UnwrapPSMEvent() FooPSMEvent {
	if ee == nil {
		return nil
	}
	switch v := ee.Type.(type) {
	case *FooEventType_Created_:
		return v.Created
	case *FooEventType_Updated_:
		return v.Updated
	case *FooEventType_Deleted_:
		return v.Deleted
	default:
		return nil
	}
}
func (ee *FooEventType) PSMEventKey() FooPSMEventKey {
	tt := ee.UnwrapPSMEvent()
	if tt == nil {
		return FooPSMEventNil
	}
	return tt.PSMEventKey()
}
func (ee *FooEvent) PSMEventKey() FooPSMEventKey {
	return ee.Event.PSMEventKey()
}
func (ee *FooEvent) UnwrapPSMEvent() FooPSMEvent {
	return ee.Event.UnwrapPSMEvent()
}
func (ee *FooEvent) SetPSMEvent(inner FooPSMEvent) {
	if ee.Event == nil {
		ee.Event = &FooEventType{}
	}
	switch v := inner.(type) {
	case *FooEventType_Created:
		ee.Event.Type = &FooEventType_Created_{Created: v}
	case *FooEventType_Updated:
		ee.Event.Type = &FooEventType_Updated_{Updated: v}
	case *FooEventType_Deleted:
		ee.Event.Type = &FooEventType_Deleted_{Deleted: v}
	default:
		panic("invalid type")
	}
}
func (*FooEventType_Created) PSMEventKey() FooPSMEventKey {
	return FooPSMEventCreated
}
func (*FooEventType_Updated) PSMEventKey() FooPSMEventKey {
	return FooPSMEventUpdated
}
func (*FooEventType_Deleted) PSMEventKey() FooPSMEventKey {
	return FooPSMEventDeleted
}

// State Query Service for %sfoo
// QuerySet is the query set for the Foo service.

type FooPSMQuerySet = psm.StateQuerySet[
	*GetFooRequest,
	*GetFooResponse,
	*ListFoosRequest,
	*ListFoosResponse,
	*ListFooEventsRequest,
	*ListFooEventsResponse,
]

func NewFooPSMQuerySet(
	smSpec psm.QuerySpec[
		*GetFooRequest,
		*GetFooResponse,
		*ListFoosRequest,
		*ListFoosResponse,
		*ListFooEventsRequest,
		*ListFooEventsResponse,
	],
	options psm.StateQueryOptions,
) (*FooPSMQuerySet, error) {
	return psm.BuildStateQuerySet[
		*GetFooRequest,
		*GetFooResponse,
		*ListFoosRequest,
		*ListFoosResponse,
		*ListFooEventsRequest,
		*ListFooEventsResponse,
	](smSpec, options)
}

type FooPSMQuerySpec = psm.QuerySpec[
	*GetFooRequest,
	*GetFooResponse,
	*ListFoosRequest,
	*ListFoosResponse,
	*ListFooEventsRequest,
	*ListFooEventsResponse,
]

func DefaultFooPSMQuerySpec(tableSpec psm.StateTableSpec) FooPSMQuerySpec {
	return psm.QuerySpec[
		*GetFooRequest,
		*GetFooResponse,
		*ListFoosRequest,
		*ListFoosResponse,
		*ListFooEventsRequest,
		*ListFooEventsResponse,
	]{
		StateTableSpec: tableSpec,
		ListRequestFilter: func(req *ListFoosRequest) (map[string]interface{}, error) {
			filter := map[string]interface{}{}
			if req.TenantId != nil {
				filter["tenant_id"] = *req.TenantId
			}
			return filter, nil
		},
		ListEventsRequestFilter: func(req *ListFooEventsRequest) (map[string]interface{}, error) {
			filter := map[string]interface{}{}
			filter["foo_id"] = req.FooId
			return filter, nil
		},
	}
}
