// Code generated by protoc-gen-go-psm. DO NOT EDIT.

package testpb

import (
	context "context"
	fmt "fmt"
	psm_pb "github.com/pentops/protostate/gen/state/v1/psm_pb"
	pgstore "github.com/pentops/protostate/pgstore"
	psm "github.com/pentops/protostate/psm"
	sqrlx "github.com/pentops/sqrlx.go/sqrlx"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
)

// PSM BarPSM

type BarPSM = psm.StateMachine[
	*BarKeys,      // implements psm.IKeyset
	*BarState,     // implements psm.IState
	BarStatus,     // implements psm.IStatusEnum
	*BarStateData, // implements psm.IStateData
	*BarEvent,     // implements psm.IEvent
	BarPSMEvent,   // implements psm.IInnerEvent
]

type BarPSMDB = psm.DBStateMachine[
	*BarKeys,      // implements psm.IKeyset
	*BarState,     // implements psm.IState
	BarStatus,     // implements psm.IStatusEnum
	*BarStateData, // implements psm.IStateData
	*BarEvent,     // implements psm.IEvent
	BarPSMEvent,   // implements psm.IInnerEvent
]

type BarPSMEventer = psm.Eventer[
	*BarKeys,      // implements psm.IKeyset
	*BarState,     // implements psm.IState
	BarStatus,     // implements psm.IStatusEnum
	*BarStateData, // implements psm.IStateData
	*BarEvent,     // implements psm.IEvent
	BarPSMEvent,   // implements psm.IInnerEvent
]

type BarPSMEventSpec = psm.EventSpec[
	*BarKeys,      // implements psm.IKeyset
	*BarState,     // implements psm.IState
	BarStatus,     // implements psm.IStatusEnum
	*BarStateData, // implements psm.IStateData
	*BarEvent,     // implements psm.IEvent
	BarPSMEvent,   // implements psm.IInnerEvent
]

type BarPSMEventKey = string

const (
	BarPSMEventNil     BarPSMEventKey = "<nil>"
	BarPSMEventCreated BarPSMEventKey = "created"
	BarPSMEventUpdated BarPSMEventKey = "updated"
	BarPSMEventDeleted BarPSMEventKey = "deleted"
)

// EXTEND BarKeys with the psm.IKeyset interface

// PSMIsSet is a helper for != nil, which does not work with generic parameters
func (msg *BarKeys) PSMIsSet() bool {
	return msg != nil
}

// PSMFullName returns the full name of state machine with package prefix
func (msg *BarKeys) PSMFullName() string {
	return "test.v1.bar"
}

// EXTEND BarState with the psm.IState interface

// PSMIsSet is a helper for != nil, which does not work with generic parameters
func (msg *BarState) PSMIsSet() bool {
	return msg != nil
}

func (msg *BarState) PSMMetadata() *psm_pb.StateMetadata {
	if msg.Metadata == nil {
		msg.Metadata = &psm_pb.StateMetadata{}
	}
	return msg.Metadata
}

func (msg *BarState) PSMKeys() *BarKeys {
	return msg.Keys
}

func (msg *BarState) SetStatus(status BarStatus) {
	msg.Status = status
}

func (msg *BarState) SetPSMKeys(inner *BarKeys) {
	msg.Keys = inner
}

func (msg *BarState) PSMData() *BarStateData {
	if msg.Data == nil {
		msg.Data = &BarStateData{}
	}
	return msg.Data
}

// EXTEND BarStateData with the psm.IStateData interface

// PSMIsSet is a helper for != nil, which does not work with generic parameters
func (msg *BarStateData) PSMIsSet() bool {
	return msg != nil
}

// EXTEND BarEvent with the psm.IEvent interface

// PSMIsSet is a helper for != nil, which does not work with generic parameters
func (msg *BarEvent) PSMIsSet() bool {
	return msg != nil
}

func (msg *BarEvent) PSMMetadata() *psm_pb.EventMetadata {
	if msg.Metadata == nil {
		msg.Metadata = &psm_pb.EventMetadata{}
	}
	return msg.Metadata
}

func (msg *BarEvent) PSMKeys() *BarKeys {
	return msg.Keys
}

func (msg *BarEvent) SetPSMKeys(inner *BarKeys) {
	msg.Keys = inner
}

// PSMEventKey returns the BarPSMEventPSMEventKey for the event, implementing psm.IEvent
func (msg *BarEvent) PSMEventKey() BarPSMEventKey {
	tt := msg.UnwrapPSMEvent()
	if tt == nil {
		return BarPSMEventNil
	}
	return tt.PSMEventKey()
}

// UnwrapPSMEvent implements psm.IEvent, returning the inner event message
func (msg *BarEvent) UnwrapPSMEvent() BarPSMEvent {
	if msg == nil {
		return nil
	}
	if msg.Event == nil {
		return nil
	}
	switch v := msg.Event.Type.(type) {
	case *BarEventType_Created_:
		return v.Created
	case *BarEventType_Updated_:
		return v.Updated
	case *BarEventType_Deleted_:
		return v.Deleted
	default:
		return nil
	}
}

// SetPSMEvent sets the inner event message from a concrete type, implementing psm.IEvent
func (msg *BarEvent) SetPSMEvent(inner BarPSMEvent) error {
	if msg.Event == nil {
		msg.Event = &BarEventType{}
	}
	switch v := inner.(type) {
	case *BarEventType_Created:
		msg.Event.Type = &BarEventType_Created_{Created: v}
	case *BarEventType_Updated:
		msg.Event.Type = &BarEventType_Updated_{Updated: v}
	case *BarEventType_Deleted:
		msg.Event.Type = &BarEventType_Deleted_{Deleted: v}
	default:
		return fmt.Errorf("invalid type %T for BarEventType", v)
	}
	return nil
}

type BarPSMEvent interface {
	psm.IInnerEvent
	PSMEventKey() BarPSMEventKey
}

// EXTEND BarEventType_Created with the BarPSMEvent interface

// PSMIsSet is a helper for != nil, which does not work with generic parameters
func (msg *BarEventType_Created) PSMIsSet() bool {
	return msg != nil
}

func (*BarEventType_Created) PSMEventKey() BarPSMEventKey {
	return BarPSMEventCreated
}

// EXTEND BarEventType_Updated with the BarPSMEvent interface

// PSMIsSet is a helper for != nil, which does not work with generic parameters
func (msg *BarEventType_Updated) PSMIsSet() bool {
	return msg != nil
}

func (*BarEventType_Updated) PSMEventKey() BarPSMEventKey {
	return BarPSMEventUpdated
}

// EXTEND BarEventType_Deleted with the BarPSMEvent interface

// PSMIsSet is a helper for != nil, which does not work with generic parameters
func (msg *BarEventType_Deleted) PSMIsSet() bool {
	return msg != nil
}

func (*BarEventType_Deleted) PSMEventKey() BarPSMEventKey {
	return BarPSMEventDeleted
}

type BarPSMTableSpec = psm.PSMTableSpec[
	*BarKeys,      // implements psm.IKeyset
	*BarState,     // implements psm.IState
	BarStatus,     // implements psm.IStatusEnum
	*BarStateData, // implements psm.IStateData
	*BarEvent,     // implements psm.IEvent
	BarPSMEvent,   // implements psm.IInnerEvent
]

var DefaultBarPSMTableSpec = BarPSMTableSpec{
	TableMap: psm.TableMap{
		State: psm.StateTableSpec{
			TableName: "bar",
			Root:      &pgstore.ProtoFieldSpec{ColumnName: "state", PathFromRoot: pgstore.ProtoPathSpec{}},
		},
		Event: psm.EventTableSpec{
			TableName:     "bar_event",
			Root:          &pgstore.ProtoFieldSpec{ColumnName: "data", PathFromRoot: pgstore.ProtoPathSpec{}},
			ID:            &pgstore.ProtoFieldSpec{ColumnName: "id", PathFromRoot: pgstore.ProtoPathSpec{"metadata", "event_id"}},
			Timestamp:     &pgstore.ProtoFieldSpec{ColumnName: "timestamp", PathFromRoot: pgstore.ProtoPathSpec{"metadata"}},
			Sequence:      &pgstore.ProtoFieldSpec{ColumnName: "sequence", PathFromRoot: pgstore.ProtoPathSpec{"metadata"}},
			StateSnapshot: &pgstore.ProtoFieldSpec{ColumnName: "state", PathFromRoot: pgstore.ProtoPathSpec{"keys"}},
		},
		KeyColumns: []psm.KeyColumn{{
			ColumnName: "bar_id",
			ProtoName:  protoreflect.Name("bar_id"),
			Primary:    true,
			Required:   true,
		}, {
			ColumnName: "bar_other_id",
			ProtoName:  protoreflect.Name("bar_other_id"),
			Primary:    true,
			Required:   true,
		}},
	},
	KeyValues: func(keys *BarKeys) (map[string]string, error) {
		keyset := map[string]string{
			"bar_id":       keys.BarId,
			"bar_other_id": keys.BarOtherId,
		}
		return keyset, nil
	},
}

func DefaultBarPSMConfig() *psm.StateMachineConfig[
	*BarKeys,      // implements psm.IKeyset
	*BarState,     // implements psm.IState
	BarStatus,     // implements psm.IStatusEnum
	*BarStateData, // implements psm.IStateData
	*BarEvent,     // implements psm.IEvent
	BarPSMEvent,   // implements psm.IInnerEvent
] {
	return psm.NewStateMachineConfig[
		*BarKeys,      // implements psm.IKeyset
		*BarState,     // implements psm.IState
		BarStatus,     // implements psm.IStatusEnum
		*BarStateData, // implements psm.IStateData
		*BarEvent,     // implements psm.IEvent
		BarPSMEvent,   // implements psm.IInnerEvent
	](DefaultBarPSMTableSpec)
}

func NewBarPSM(config *psm.StateMachineConfig[
	*BarKeys,      // implements psm.IKeyset
	*BarState,     // implements psm.IState
	BarStatus,     // implements psm.IStatusEnum
	*BarStateData, // implements psm.IStateData
	*BarEvent,     // implements psm.IEvent
	BarPSMEvent,   // implements psm.IInnerEvent
]) (*BarPSM, error) {
	return psm.NewStateMachine[
		*BarKeys,      // implements psm.IKeyset
		*BarState,     // implements psm.IState
		BarStatus,     // implements psm.IStatusEnum
		*BarStateData, // implements psm.IStateData
		*BarEvent,     // implements psm.IEvent
		BarPSMEvent,   // implements psm.IInnerEvent
	](config)
}

func BarPSMMutation[SE BarPSMEvent](cb func(*BarStateData, SE) error) psm.PSMMutationFunc[
	*BarKeys,      // implements psm.IKeyset
	*BarState,     // implements psm.IState
	BarStatus,     // implements psm.IStatusEnum
	*BarStateData, // implements psm.IStateData
	*BarEvent,     // implements psm.IEvent
	BarPSMEvent,   // implements psm.IInnerEvent
	SE,            // Specific event type for the transition
] {
	return psm.PSMMutationFunc[
		*BarKeys,      // implements psm.IKeyset
		*BarState,     // implements psm.IState
		BarStatus,     // implements psm.IStatusEnum
		*BarStateData, // implements psm.IStateData
		*BarEvent,     // implements psm.IEvent
		BarPSMEvent,   // implements psm.IInnerEvent
		SE,            // Specific event type for the transition
	](cb)
}

type BarPSMHookBaton = psm.HookBaton[
	*BarKeys,      // implements psm.IKeyset
	*BarState,     // implements psm.IState
	BarStatus,     // implements psm.IStatusEnum
	*BarStateData, // implements psm.IStateData
	*BarEvent,     // implements psm.IEvent
	BarPSMEvent,   // implements psm.IInnerEvent
]

func BarPSMHook[SE BarPSMEvent](cb func(context.Context, sqrlx.Transaction, BarPSMHookBaton, *BarState, SE) error) psm.PSMHookFunc[
	*BarKeys,      // implements psm.IKeyset
	*BarState,     // implements psm.IState
	BarStatus,     // implements psm.IStatusEnum
	*BarStateData, // implements psm.IStateData
	*BarEvent,     // implements psm.IEvent
	BarPSMEvent,   // implements psm.IInnerEvent
	SE,            // Specific event type for the transition
] {
	return psm.PSMHookFunc[
		*BarKeys,      // implements psm.IKeyset
		*BarState,     // implements psm.IState
		BarStatus,     // implements psm.IStatusEnum
		*BarStateData, // implements psm.IStateData
		*BarEvent,     // implements psm.IEvent
		BarPSMEvent,   // implements psm.IInnerEvent
		SE,            // Specific event type for the transition
	](cb)
}
func BarPSMGeneralHook(cb func(context.Context, sqrlx.Transaction, BarPSMHookBaton, *BarState, *BarEvent) error) psm.GeneralStateHook[
	*BarKeys,      // implements psm.IKeyset
	*BarState,     // implements psm.IState
	BarStatus,     // implements psm.IStatusEnum
	*BarStateData, // implements psm.IStateData
	*BarEvent,     // implements psm.IEvent
	BarPSMEvent,   // implements psm.IInnerEvent
] {
	return psm.GeneralStateHook[
		*BarKeys,      // implements psm.IKeyset
		*BarState,     // implements psm.IState
		BarStatus,     // implements psm.IStatusEnum
		*BarStateData, // implements psm.IStateData
		*BarEvent,     // implements psm.IEvent
		BarPSMEvent,   // implements psm.IInnerEvent
	](cb)
}
