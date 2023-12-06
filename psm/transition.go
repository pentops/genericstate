package psm

import (
	"context"
	"fmt"

	"google.golang.org/protobuf/proto"
)

/*
# Generic Type Parameter Sets

Two sets of generic type sets exist:

`S ST E IE`
`S ST E IE SE`

Both share the same types, as follows

### `S IState[ST]`

The State Entity, implements IState[ST], e.g. *testpb.FooState

### `ST IStatusEnum`

The Status Enum of the state entity implements IStatusEnum,
e.g. *testpb.FooStatus (int32)

### `E IEvent[IE]`

The Event Wrapper, the top level which has metadata, foreign keys to the state, and the event itself.
e.g. *testpb.FooEvent, the concrete proto message

### `IE IInnerEvent`

The Inner Event, the typed event *interface* which is the set of all possible events for the state machine
e.g. testpb.FooPSMEvent interface - this is generated by the protoc plugin in _psm.pb.go
It is set at compile time specifically to the interface type.

### `SE IInnerEvent`

The Specific single typed event *struct* which is the specific event for the transition.
SE implements the same interface of IE.
e.g. *testpb.FooPSMEvent_Created, the concrete proto message which implements testpb.FooPSMEvent


The state machine deals with the first shorter chain, as it deals with all events.
Transitions deal with a single specific event type, so have the extra SE parameter.

S, ST, E, and IE are set to one single type for the entire state machine
SE is set to a single type for each transition.
*/

type IStatusEnum interface {
	~int32
	ShortString() string
}

type IState[Status IStatusEnum] interface {
	proto.Message
	GetStatus() Status
}

type IEvent[Inner any] interface {
	proto.Message
	UnwrapPSMEvent() Inner
}

type IInnerEvent interface {
	proto.Message
}

type TypedTransitionHandler[
	S proto.Message,
	ST IStatusEnum,
	E IEvent[IE],
	IE IInnerEvent,
] interface {
	RunTransition(context.Context, TransitionBaton[E, IE], S, IE) error
	HandlesEvent(E) bool
}

type TransitionFunc[
	S proto.Message,
	ST IStatusEnum,
	E IEvent[IE],
	IE IInnerEvent,
	SE IInnerEvent,
] func(context.Context, TransitionBaton[E, IE], S, SE) error

// RunTransition implements TransitionHandler, where SE is the specific event
// cast from the interface IE provided in the call
func (f TransitionFunc[S, ST, E, IE, SE]) RunTransition(
	ctx context.Context,
	tb TransitionBaton[E, IE],
	state S,
	event IE,
) error {
	// Cast the interface ET IInnerEvent to the specific type of event which
	// this func handles
	asType, ok := any(event).(SE)
	if !ok {
		return fmt.Errorf("unexpected event type: %T", event)
	}
	return f(ctx, tb, state, asType)
}

func (f TransitionFunc[S, ST, E, IE, SE]) HandlesEvent(outerEvent E) bool {
	// Check if the parameter passed as ET (IInnerEvent) is the specific type
	// (IE, also IInnerEvent, but typed) which this transition handles
	event := outerEvent.UnwrapPSMEvent()
	_, ok := any(event).(SE)
	return ok
}

type TransitionFilter[
	S IState[ST],
	ST IStatusEnum,
	E IEvent[IE],
	IE IInnerEvent,
] interface {
	Matches(S, IE) bool
}

type CastTransitionFilter[
	S IState[ST],
	ST IStatusEnum,
	E IEvent[IE],
	IE IInnerEvent,
] struct {
}

type TypedTransition[
	S IState[ST],
	ST IStatusEnum,
	E IEvent[IE],
	IE IInnerEvent,
] struct {
	handler     TypedTransitionHandler[S, ST, E, IE]
	fromStatus  []ST
	eventFilter func(E) bool
}

func (f TypedTransition[S, ST, E, IE]) Matches(state S, outerEvent E) bool {
	if !f.handler.HandlesEvent(outerEvent) {
		return false
	}

	if f.fromStatus != nil {
		didMatch := false
		currentStatus := state.GetStatus()
		for _, fromStatus := range f.fromStatus {
			if fromStatus == currentStatus {
				didMatch = true
				break
			}
		}
		if !didMatch {
			return false
		}
	}

	if f.eventFilter != nil && !f.eventFilter(outerEvent) {
		return false
	}
	return true
}

func (f TypedTransition[S, ST, E, IE]) RunTransition(
	ctx context.Context,
	tb TransitionBaton[E, IE],
	state S,
	event IE,
) error {
	return f.handler.RunTransition(ctx, tb, state, event)
}
