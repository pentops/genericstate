package psm

import (
	"github.com/pentops/protostate/gen/state/v1/psm_pb"
	"google.golang.org/protobuf/proto"
)

/*
# Generic Type Parameter Sets

Two sets of generic type sets exist:

`K S ST SD E IE`
`K S ST SD E IE SE`

Both share the same types, as follows, and defined below

### `K IKeyset`
### `S IState[K, ST, SD]`
### `ST IStatusEnum`
### `SD IStateData`
### `SD IStateData,
E IEvent[K, S, ST, SD, IE]`
### `IE IInnerEvent`
### `SE IInnerEvent`

The Specific single typed event *struct* which is the specific event for the transition.
SE implements the same interface of IE.
e.g. *testpb.FooPSMEvent_Created, the concrete proto message which implements testpb.FooPSMEvent


The state machine deals with the first shorter chain, as it deals with all events.
Transitions deal with a single specific event type, so have the extra SE parameter.

K, S, ST, E, and IE are set to one single type for the entire state machine
SE is set to a single type for each transition.
*/

// IGenericProtoMessage is the base extensions shared by all message entities in the PSM generated code
type IPSMMessage interface {
	proto.Message
	PSMIsSet() bool
}

// IStatusEnum is enum representing the named state of the entity.
// e.g. *testpb.FooStatus (int32)
type IStatusEnum interface {
	~int32
	ShortString() string
	String() string
}

type IKeyset interface {
	IPSMMessage
	PSMFullName() string
	PSMKeyValues() (map[string]string, error)
}

// IState[K, ST, SD]is the main State Entity e.g. *testpb.FooState
type IState[K IKeyset, ST IStatusEnum, SD IStateData] interface {
	IPSMMessage
	GetStatus() ST
	SetStatus(ST)
	PSMMetadata() *psm_pb.StateMetadata
	PSMKeys() K
	SetPSMKeys(K)
	PSMData() SD
}

// IStateData is the Data Entity e.g. *testpb.FooStateData
type IStateData interface {
	IPSMMessage
}

// IEvent is the Event Wrapper, the top level which has metadata, foreign keys to the state, and the event itself.
// e.g. *testpb.FooEvent, the concrete proto message
type IEvent[
	K IKeyset,
	S IState[K, ST, SD],
	ST IStatusEnum,
	SD IStateData,
	Inner any,
] interface {
	proto.Message
	UnwrapPSMEvent() Inner
	SetPSMEvent(Inner) error
	PSMKeys() K
	SetPSMKeys(K)
	PSMMetadata() *psm_pb.EventMetadata
	PSMIsSet() bool
}

// IInnerEvent is the typed event *interface* which is the set of all possible events for the state machine
// e.g. testpb.FooPSMEvent interface - this is generated by the protoc plugin in _psm.pb.go
// It is set at compile time specifically to the interface type.
type IInnerEvent interface {
	IPSMMessage
	PSMEventKey() string
}
