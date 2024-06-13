package psm

import (
	"context"
	"fmt"
	"strings"

	"github.com/pentops/log.go/log"
	"github.com/pentops/o5-messaging/o5msg"
	"github.com/pentops/protostate/gen/state/v1/psm_pb"
	"github.com/pentops/sqrlx.go/sqrlx"
)

type hookBaton[
	K IKeyset,
	S IState[K, ST, SD],
	ST IStatusEnum,
	SD IStateData,
	E IEvent[K, S, ST, SD, IE],
	IE IInnerEvent,
] struct {
	sideEffects []o5msg.Message
	chainEvents []IE
	causedBy    E
}

func (td *hookBaton[K, S, ST, SD, E, IE]) ChainEvent(inner IE) {
	td.chainEvents = append(td.chainEvents, inner)
}

func (td *hookBaton[K, S, ST, SD, E, IE]) SideEffect(msg o5msg.Message) {
	td.sideEffects = append(td.sideEffects, msg)
}

func (td *hookBaton[K, S, ST, SD, E, IE]) FullCause() E {
	return td.causedBy
}

func (td *hookBaton[K, S, ST, SD, E, IE]) AsCause() *psm_pb.Cause {
	causeMetadata := td.causedBy.PSMMetadata()
	return &psm_pb.Cause{
		Type: &psm_pb.Cause_PsmEvent{
			PsmEvent: &psm_pb.PSMEventCause{
				EventId:      causeMetadata.EventId,
				StateMachine: td.causedBy.PSMKeys().PSMFullName(),
			},
		},
	}
}

type transitionSpec[
	K IKeyset,
	S IState[K, ST, SD],
	ST IStatusEnum,
	SD IStateData,
	E IEvent[K, S, ST, SD, IE],
	IE IInnerEvent,
] struct {
	fromStatus []ST
	toStatus   ST
	noop       bool
	eventType  string
	mutations  []transitionMutationWrapper[K, S, ST, SD, E, IE]
	logic      []transitionLogicHookWrapper[K, S, ST, SD, E, IE]
	data       []transitionDataHookWrapper[K, S, ST, SD, E, IE]
}

func (hs *transitionSpec[K, S, ST, SD, E, IE]) runTransitionMutations(
	ctx context.Context,
	state S,
	event E,
) error {

	log.Debug(ctx, "running transition mutations")
	sd := state.PSMData()
	innerEvent := event.UnwrapPSMEvent()

	if hs.noop {
		return nil
	}

	if hs.toStatus != 0 {
		log.WithField(ctx, "toStatus", hs.toStatus).Debug("setting status")
		state.SetStatus(hs.toStatus)
	}

	for _, mutation := range hs.mutations {
		log.WithField(ctx, "mutation", mutation).Debug("running mutation")
		err := mutation.TransitionMutation(sd, innerEvent)
		if err != nil {
			return err
		}
	}
	log.WithField(ctx, "mutationCount", hs.mutations).Debug("mutations complete")
	return nil
}

func (hs *transitionSpec[K, S, ST, SD, E, IE]) runTransitionHooks(
	ctx context.Context,
	tx sqrlx.Transaction,
	baton HookBaton[K, S, ST, SD, E, IE],
	state S,
	event E,
) error {
	innerEvent := event.UnwrapPSMEvent()

	log.Debug(ctx, "running transition mutations")
	for _, logic := range hs.logic {
		log.WithField(ctx, "logic", logic).Debug("running logic hook")
		err := logic.TransitionLogicHook(ctx, baton, state, innerEvent)
		if err != nil {
			return err
		}
	}

	for _, data := range hs.data {
		log.WithField(ctx, "data", data).Debug("running data hook")
		err := data.TransitionDataHook(ctx, tx, state, innerEvent)
		if err != nil {
			return err
		}
	}

	log.WithFields(ctx, map[string]interface{}{
		"logicCount": len(hs.logic),
		"dataCount":  len(hs.data),
	}).Debug("transition hooks complete")

	return nil
}

type transitionSet[
	K IKeyset,
	S IState[K, ST, SD],
	ST IStatusEnum,
	SD IStateData,
	E IEvent[K, S, ST, SD, IE],
	IE IInnerEvent,
] struct {
	logic     []generalLogicHookWrapper[K, S, ST, SD, E, IE]
	stateData []generalStateDataWrapper[K, S, ST, SD, E, IE]
	eventData []generalEventDataHookWrapper[K, S, ST, SD, E, IE]

	transitions []*transitionSpec[K, S, ST, SD, E, IE]
}

func (hs *transitionSet[K, S, ST, SD, E, IE]) LogicHook(hook GeneralLogicHook[K, S, ST, SD, E, IE]) {
	hs.logic = append(hs.logic, hook)
}

func (hs *transitionSet[K, S, ST, SD, E, IE]) StateDataHook(hook GeneralStateDataHook[K, S, ST, SD, E, IE]) {
	hs.stateData = append(hs.stateData, hook)
}

func (hs *transitionSet[K, S, ST, SD, E, IE]) EventDataHook(hook GeneralEventDataHook[K, S, ST, SD, E, IE]) {
	hs.eventData = append(hs.eventData, hook)
}

func (hs *transitionSet[K, S, ST, SD, E, IE]) runGlobalTransitionHooks(
	ctx context.Context,
	tx sqrlx.Transaction,
	baton HookBaton[K, S, ST, SD, E, IE],
	state S,
	event E,
) error {
	for _, hook := range hs.logic {
		err := hook.GeneralLogicHook(ctx, baton, state, event)
		if err != nil {
			return err
		}
	}

	for _, hook := range hs.stateData {
		err := hook.GeneralStateDataHook(ctx, tx, state)
		if err != nil {
			return err
		}
	}

	for _, hook := range hs.eventData {
		err := hook.GeneralEventDataHook(ctx, tx, state, event)
		if err != nil {
			return err
		}
	}

	return nil
}

func (hs transitionSet[K, S, ST, SD, E, IE]) findTransitions(status ST, wantType string) (*transitionSpec[K, S, ST, SD, E, IE], error) {
	hooks := make([]*transitionSpec[K, S, ST, SD, E, IE], 0, 1)

	for _, search := range hs.transitions {
		if search.eventType != wantType {
			continue
		}
		if len(search.fromStatus) == 0 {
			hooks = append(hooks, search)
			continue
		}
		for _, fromStatus := range search.fromStatus {
			if fromStatus == status {
				hooks = append(hooks, search)
				break
			}
		}
	}

	if len(hooks) == 0 {
		return nil, fmt.Errorf("no transition found for %s on %s", status, wantType)
	}

	if len(hooks) == 1 {
		return hooks[0], nil
	}

	merged, err := hs.mergeHooks(status, wantType, hooks)
	if err != nil {
		return nil, err
	}

	return merged, nil

}

func (hs *transitionSet[K, S, ST, SD, E, IE]) mergeHooks(status ST, eventType string, hooks []*transitionSpec[K, S, ST, SD, E, IE]) (*transitionSpec[K, S, ST, SD, E, IE], error) {
	merged := &transitionSpec[K, S, ST, SD, E, IE]{
		fromStatus: []ST{status},
		eventType:  eventType,
	}

	for _, hook := range hooks {
		merged.mutations = append(merged.mutations, hook.mutations...)
		merged.logic = append(merged.logic, hook.logic...)
		merged.data = append(merged.data, hook.data...)

		if hook.toStatus != 0 {
			if merged.toStatus == 0 {
				merged.toStatus = hook.toStatus
			} else if merged.toStatus != hook.toStatus {
				return nil, fmt.Errorf("conflicting toStatus transitions for fromStatus %q event %q", status.ShortString(), eventType)
			}
		}
	}

	return merged, nil
}

func (hs transitionSet[K, S, ST, SD, E, IE]) PrintMermaid() (string, error) {

	lines := []string{
		"stateDiagram-v2",
	}

	type specSet struct {
		status      ST
		event       string
		transitions []*transitionSpec[K, S, ST, SD, E, IE]
	}
	byKey := map[string]specSet{}

	for _, transition := range hs.transitions {
		var key string
		if len(transition.fromStatus) == 0 {
			continue // TODO: Handle GLobal transitions
		}
		for _, from := range transition.fromStatus {
			key = fmt.Sprintf("%s-%s", from.ShortString(), transition.eventType)
			entry, ok := byKey[key]
			if !ok {
				entry = specSet{
					status: from,
					event:  transition.eventType,
				}
			}
			entry.transitions = append(entry.transitions, transition)
			byKey[key] = entry
		}
	}

	for _, spec := range byKey {
		merged, err := hs.mergeHooks(spec.status, spec.event, spec.transitions)
		if err != nil {
			return "", err
		}
		if merged.toStatus == 0 {
			merged.toStatus = merged.fromStatus[0]
		}

		var fromString string
		if merged.fromStatus[0] == 0 {
			fromString = "[*]"
		} else {
			fromString = merged.fromStatus[0].ShortString()
		}
		lines = append(lines, fmt.Sprintf("%s --> %s : %s", fromString, merged.toStatus.ShortString(), merged.eventType))
	}

	return strings.Join(lines, "\n"), nil
}
