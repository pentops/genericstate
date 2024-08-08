// Code generated by protoc-gen-go-psm. DO NOT EDIT.

package test_pb

import (
	context "context"
	fmt "fmt"
	psm_pb "github.com/pentops/j5/gen/psm/state/v1/psm_pb"
	psm "github.com/pentops/protostate/psm"
	sqrlx "github.com/pentops/sqrlx.go/sqrlx"
)

// PSM FooPSM

type FooPSM = psm.StateMachine[
	*FooKeys,      // implements psm.IKeyset
	*FooState,     // implements psm.IState
	FooStatus,     // implements psm.IStatusEnum
	*FooStateData, // implements psm.IStateData
	*FooEvent,     // implements psm.IEvent
	FooPSMEvent,   // implements psm.IInnerEvent
]

type FooPSMDB = psm.DBStateMachine[
	*FooKeys,      // implements psm.IKeyset
	*FooState,     // implements psm.IState
	FooStatus,     // implements psm.IStatusEnum
	*FooStateData, // implements psm.IStateData
	*FooEvent,     // implements psm.IEvent
	FooPSMEvent,   // implements psm.IInnerEvent
]

type FooPSMEventSpec = psm.EventSpec[
	*FooKeys,      // implements psm.IKeyset
	*FooState,     // implements psm.IState
	FooStatus,     // implements psm.IStatusEnum
	*FooStateData, // implements psm.IStateData
	*FooEvent,     // implements psm.IEvent
	FooPSMEvent,   // implements psm.IInnerEvent
]

type FooPSMEventKey = string

const (
	FooPSMEventNil     FooPSMEventKey = "<nil>"
	FooPSMEventCreated FooPSMEventKey = "created"
	FooPSMEventUpdated FooPSMEventKey = "updated"
	FooPSMEventDeleted FooPSMEventKey = "deleted"
)

// EXTEND FooKeys with the psm.IKeyset interface

// PSMIsSet is a helper for != nil, which does not work with generic parameters
func (msg *FooKeys) PSMIsSet() bool {
	return msg != nil
}

// PSMFullName returns the full name of state machine with package prefix
func (msg *FooKeys) PSMFullName() string {
	return "test.v1.foo"
}
func (msg *FooKeys) PSMKeyValues() (map[string]string, error) {
	keyset := map[string]string{
		"foo_id":         msg.FooId,
		"meta_tenant_id": msg.MetaTenantId,
	}
	if msg.TenantId != nil {
		keyset["tenant_id"] = *msg.TenantId
	}
	return keyset, nil
}

// EXTEND FooState with the psm.IState interface

// PSMIsSet is a helper for != nil, which does not work with generic parameters
func (msg *FooState) PSMIsSet() bool {
	return msg != nil
}

func (msg *FooState) PSMMetadata() *psm_pb.StateMetadata {
	if msg.Metadata == nil {
		msg.Metadata = &psm_pb.StateMetadata{}
	}
	return msg.Metadata
}

func (msg *FooState) PSMKeys() *FooKeys {
	return msg.Keys
}

func (msg *FooState) SetStatus(status FooStatus) {
	msg.Status = status
}

func (msg *FooState) SetPSMKeys(inner *FooKeys) {
	msg.Keys = inner
}

func (msg *FooState) PSMData() *FooStateData {
	if msg.Data == nil {
		msg.Data = &FooStateData{}
	}
	return msg.Data
}

// EXTEND FooStateData with the psm.IStateData interface

// PSMIsSet is a helper for != nil, which does not work with generic parameters
func (msg *FooStateData) PSMIsSet() bool {
	return msg != nil
}

// EXTEND FooEvent with the psm.IEvent interface

// PSMIsSet is a helper for != nil, which does not work with generic parameters
func (msg *FooEvent) PSMIsSet() bool {
	return msg != nil
}

func (msg *FooEvent) PSMMetadata() *psm_pb.EventMetadata {
	if msg.Metadata == nil {
		msg.Metadata = &psm_pb.EventMetadata{}
	}
	return msg.Metadata
}

func (msg *FooEvent) PSMKeys() *FooKeys {
	return msg.Keys
}

func (msg *FooEvent) SetPSMKeys(inner *FooKeys) {
	msg.Keys = inner
}

// PSMEventKey returns the FooPSMEventPSMEventKey for the event, implementing psm.IEvent
func (msg *FooEvent) PSMEventKey() FooPSMEventKey {
	tt := msg.UnwrapPSMEvent()
	if tt == nil {
		return FooPSMEventNil
	}
	return tt.PSMEventKey()
}

// UnwrapPSMEvent implements psm.IEvent, returning the inner event message
func (msg *FooEvent) UnwrapPSMEvent() FooPSMEvent {
	if msg == nil {
		return nil
	}
	if msg.Event == nil {
		return nil
	}
	switch v := msg.Event.Type.(type) {
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

// SetPSMEvent sets the inner event message from a concrete type, implementing psm.IEvent
func (msg *FooEvent) SetPSMEvent(inner FooPSMEvent) error {
	if msg.Event == nil {
		msg.Event = &FooEventType{}
	}
	switch v := inner.(type) {
	case *FooEventType_Created:
		msg.Event.Type = &FooEventType_Created_{Created: v}
	case *FooEventType_Updated:
		msg.Event.Type = &FooEventType_Updated_{Updated: v}
	case *FooEventType_Deleted:
		msg.Event.Type = &FooEventType_Deleted_{Deleted: v}
	default:
		return fmt.Errorf("invalid type %T for FooEventType", v)
	}
	return nil
}

type FooPSMEvent interface {
	psm.IInnerEvent
	PSMEventKey() FooPSMEventKey
}

// EXTEND FooEventType_Created with the FooPSMEvent interface

// PSMIsSet is a helper for != nil, which does not work with generic parameters
func (msg *FooEventType_Created) PSMIsSet() bool {
	return msg != nil
}

func (*FooEventType_Created) PSMEventKey() FooPSMEventKey {
	return FooPSMEventCreated
}

// EXTEND FooEventType_Updated with the FooPSMEvent interface

// PSMIsSet is a helper for != nil, which does not work with generic parameters
func (msg *FooEventType_Updated) PSMIsSet() bool {
	return msg != nil
}

func (*FooEventType_Updated) PSMEventKey() FooPSMEventKey {
	return FooPSMEventUpdated
}

// EXTEND FooEventType_Deleted with the FooPSMEvent interface

// PSMIsSet is a helper for != nil, which does not work with generic parameters
func (msg *FooEventType_Deleted) PSMIsSet() bool {
	return msg != nil
}

func (*FooEventType_Deleted) PSMEventKey() FooPSMEventKey {
	return FooPSMEventDeleted
}

func FooPSMBuilder() *psm.StateMachineConfig[
	*FooKeys,      // implements psm.IKeyset
	*FooState,     // implements psm.IState
	FooStatus,     // implements psm.IStatusEnum
	*FooStateData, // implements psm.IStateData
	*FooEvent,     // implements psm.IEvent
	FooPSMEvent,   // implements psm.IInnerEvent
] {
	return &psm.StateMachineConfig[
		*FooKeys,      // implements psm.IKeyset
		*FooState,     // implements psm.IState
		FooStatus,     // implements psm.IStatusEnum
		*FooStateData, // implements psm.IStateData
		*FooEvent,     // implements psm.IEvent
		FooPSMEvent,   // implements psm.IInnerEvent
	]{}
}

func FooPSMMutation[SE FooPSMEvent](cb func(*FooStateData, SE) error) psm.TransitionMutation[
	*FooKeys,      // implements psm.IKeyset
	*FooState,     // implements psm.IState
	FooStatus,     // implements psm.IStatusEnum
	*FooStateData, // implements psm.IStateData
	*FooEvent,     // implements psm.IEvent
	FooPSMEvent,   // implements psm.IInnerEvent
	SE,            // Specific event type for the transition
] {
	return psm.TransitionMutation[
		*FooKeys,      // implements psm.IKeyset
		*FooState,     // implements psm.IState
		FooStatus,     // implements psm.IStatusEnum
		*FooStateData, // implements psm.IStateData
		*FooEvent,     // implements psm.IEvent
		FooPSMEvent,   // implements psm.IInnerEvent
		SE,            // Specific event type for the transition
	](cb)
}

type FooPSMHookBaton = psm.HookBaton[
	*FooKeys,      // implements psm.IKeyset
	*FooState,     // implements psm.IState
	FooStatus,     // implements psm.IStatusEnum
	*FooStateData, // implements psm.IStateData
	*FooEvent,     // implements psm.IEvent
	FooPSMEvent,   // implements psm.IInnerEvent
]

func FooPSMLogicHook[SE FooPSMEvent](cb func(context.Context, FooPSMHookBaton, *FooState, SE) error) psm.TransitionLogicHook[
	*FooKeys,      // implements psm.IKeyset
	*FooState,     // implements psm.IState
	FooStatus,     // implements psm.IStatusEnum
	*FooStateData, // implements psm.IStateData
	*FooEvent,     // implements psm.IEvent
	FooPSMEvent,   // implements psm.IInnerEvent
	SE,            // Specific event type for the transition
] {
	return psm.TransitionLogicHook[
		*FooKeys,      // implements psm.IKeyset
		*FooState,     // implements psm.IState
		FooStatus,     // implements psm.IStatusEnum
		*FooStateData, // implements psm.IStateData
		*FooEvent,     // implements psm.IEvent
		FooPSMEvent,   // implements psm.IInnerEvent
		SE,            // Specific event type for the transition
	](cb)
}
func FooPSMDataHook[SE FooPSMEvent](cb func(context.Context, sqrlx.Transaction, *FooState, SE) error) psm.TransitionDataHook[
	*FooKeys,      // implements psm.IKeyset
	*FooState,     // implements psm.IState
	FooStatus,     // implements psm.IStatusEnum
	*FooStateData, // implements psm.IStateData
	*FooEvent,     // implements psm.IEvent
	FooPSMEvent,   // implements psm.IInnerEvent
	SE,            // Specific event type for the transition
] {
	return psm.TransitionDataHook[
		*FooKeys,      // implements psm.IKeyset
		*FooState,     // implements psm.IState
		FooStatus,     // implements psm.IStatusEnum
		*FooStateData, // implements psm.IStateData
		*FooEvent,     // implements psm.IEvent
		FooPSMEvent,   // implements psm.IInnerEvent
		SE,            // Specific event type for the transition
	](cb)
}
func FooPSMLinkHook[SE FooPSMEvent, DK psm.IKeyset, DIE psm.IInnerEvent](
	linkDestination psm.LinkDestination[DK, DIE],
	cb func(context.Context, *FooState, SE, func(DK, DIE)) error,
) psm.LinkHook[
	*FooKeys,      // implements psm.IKeyset
	*FooState,     // implements psm.IState
	FooStatus,     // implements psm.IStatusEnum
	*FooStateData, // implements psm.IStateData
	*FooEvent,     // implements psm.IEvent
	FooPSMEvent,   // implements psm.IInnerEvent
	SE,            // Specific event type for the transition
	DK,            // Destination Keys
	DIE,           // Destination Inner Event
] {
	return psm.LinkHook[
		*FooKeys,      // implements psm.IKeyset
		*FooState,     // implements psm.IState
		FooStatus,     // implements psm.IStatusEnum
		*FooStateData, // implements psm.IStateData
		*FooEvent,     // implements psm.IEvent
		FooPSMEvent,   // implements psm.IInnerEvent
		SE,            // Specific event type for the transition
		DK,            // Destination Keys
		DIE,           // Destination Inner Event
	]{
		Derive:      cb,
		Destination: linkDestination,
	}
}
func FooPSMGeneralLogicHook(cb func(context.Context, FooPSMHookBaton, *FooState, *FooEvent) error) psm.GeneralLogicHook[
	*FooKeys,      // implements psm.IKeyset
	*FooState,     // implements psm.IState
	FooStatus,     // implements psm.IStatusEnum
	*FooStateData, // implements psm.IStateData
	*FooEvent,     // implements psm.IEvent
	FooPSMEvent,   // implements psm.IInnerEvent
] {
	return psm.GeneralLogicHook[
		*FooKeys,      // implements psm.IKeyset
		*FooState,     // implements psm.IState
		FooStatus,     // implements psm.IStatusEnum
		*FooStateData, // implements psm.IStateData
		*FooEvent,     // implements psm.IEvent
		FooPSMEvent,   // implements psm.IInnerEvent
	](cb)
}
func FooPSMGeneralStateDataHook(cb func(context.Context, sqrlx.Transaction, *FooState) error) psm.GeneralStateDataHook[
	*FooKeys,      // implements psm.IKeyset
	*FooState,     // implements psm.IState
	FooStatus,     // implements psm.IStatusEnum
	*FooStateData, // implements psm.IStateData
	*FooEvent,     // implements psm.IEvent
	FooPSMEvent,   // implements psm.IInnerEvent
] {
	return psm.GeneralStateDataHook[
		*FooKeys,      // implements psm.IKeyset
		*FooState,     // implements psm.IState
		FooStatus,     // implements psm.IStatusEnum
		*FooStateData, // implements psm.IStateData
		*FooEvent,     // implements psm.IEvent
		FooPSMEvent,   // implements psm.IInnerEvent
	](cb)
}
func FooPSMGeneralEventDataHook(cb func(context.Context, sqrlx.Transaction, *FooState, *FooEvent) error) psm.GeneralEventDataHook[
	*FooKeys,      // implements psm.IKeyset
	*FooState,     // implements psm.IState
	FooStatus,     // implements psm.IStatusEnum
	*FooStateData, // implements psm.IStateData
	*FooEvent,     // implements psm.IEvent
	FooPSMEvent,   // implements psm.IInnerEvent
] {
	return psm.GeneralEventDataHook[
		*FooKeys,      // implements psm.IKeyset
		*FooState,     // implements psm.IState
		FooStatus,     // implements psm.IStatusEnum
		*FooStateData, // implements psm.IStateData
		*FooEvent,     // implements psm.IEvent
		FooPSMEvent,   // implements psm.IInnerEvent
	](cb)
}
