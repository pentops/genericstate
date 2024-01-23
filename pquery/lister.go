package pquery

import (
	"context"
	"database/sql"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	"github.com/bufbuild/protovalidate-go"
	"github.com/elgris/sqrl"
	sq "github.com/elgris/sqrl"
	"github.com/pentops/log.go/log"
	"github.com/pentops/protostate/dbconvert"
	"github.com/pentops/protostate/gen/list/v1/psml_pb"
	"github.com/pentops/sqrlx.go/sqrlx"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/dynamicpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ListRequest interface {
	proto.Message
}

type ListResponse interface {
	proto.Message
}

type ListSpec[
	REQ ListRequest,
	RES ListResponse,
] struct {
	TableName  string
	DataColumn string

	Auth     AuthProvider
	AuthJoin []*LeftJoin

	RequestFilter func(REQ) (map[string]interface{}, error)
}

type sortSpec struct {
	nestedField
	desc bool
}

type ListerOption func(*listerOptions)

type listerOptions struct {
	tieBreakerFields []string
}

// TieBreakerFields is a list of fields to use as a tie breaker when the
// list request message does not specify these fields.
func WithTieBreakerFields(fields ...string) ListerOption {
	return func(lo *listerOptions) {
		lo.tieBreakerFields = fields
	}
}

type ListReflectionSet struct {
	defaultPageSize uint64

	arrayField        protoreflect.FieldDescriptor
	pageResponseField protoreflect.FieldDescriptor
	pageRequestField  protoreflect.FieldDescriptor
	queryRequestField protoreflect.FieldDescriptor

	defaultSortFields []sortSpec
	tieBreakerFields  []sortSpec

	// fields in the list request object which become mandatory filters
	RequestFilterFields []protoreflect.FieldDescriptor
}

func resolveListerOptions(options []ListerOption) listerOptions {
	optionsStruct := listerOptions{}
	for _, option := range options {
		option(&optionsStruct)
	}
	return optionsStruct
}

func ValidateListMethod(req protoreflect.MessageDescriptor, res protoreflect.MessageDescriptor, options ...ListerOption) error {
	_, err := buildListReflection(req, res, resolveListerOptions(options))
	return err
}

func BuildListReflection(req protoreflect.MessageDescriptor, res protoreflect.MessageDescriptor, options ...ListerOption) (*ListReflectionSet, error) {
	return buildListReflection(req, res, resolveListerOptions(options))
}

func buildListReflection(req protoreflect.MessageDescriptor, res protoreflect.MessageDescriptor, options listerOptions) (*ListReflectionSet, error) {
	var err error
	ll := ListReflectionSet{
		defaultPageSize: uint64(20),
	}
	fields := res.Fields()

	for i := 0; i < fields.Len(); i++ {
		field := fields.Get(i)
		msg := field.Message()
		if msg == nil {
			return nil, fmt.Errorf("field %s is a '%s', but should be a message", field.Name(), field.Kind())
		}

		if msg.FullName() == "psm.list.v1.PageResponse" {
			ll.pageResponseField = field
			continue
		}

		if field.Cardinality() == protoreflect.Repeated {
			if ll.arrayField != nil {
				return nil, fmt.Errorf("multiple repeated fields (%s and %s)", ll.arrayField.Name(), field.Name())
			}

			ll.arrayField = field
			continue
		}
		return nil, fmt.Errorf("unknown field in response: '%s' of type %s", field.Name(), field.Kind())
	}

	if ll.arrayField == nil {
		return nil, fmt.Errorf("no repeated field in response, %s must have a repeated message", res.FullName())
	}

	if ll.pageResponseField == nil {
		return nil, fmt.Errorf("no page field in response, %s must have a psm.list.v1.PageResponse", res.FullName())
	}

	ll.defaultSortFields = buildDefaultSorts(ll.arrayField.Message().Fields())

	ll.tieBreakerFields, err = buildTieBreakerFields(req, ll.arrayField.Message(), options.tieBreakerFields)
	if err != nil {
		return nil, err
	}

	if len(ll.defaultSortFields) == 0 && len(ll.tieBreakerFields) == 0 {
		return nil, fmt.Errorf("no default sort field found, %s must have at least one field annotated as default sort, or specify a tie breaker in %s", ll.arrayField.Message().FullName(), req.FullName())
	}

	requestFields := req.Fields()
	for i := 0; i < requestFields.Len(); i++ {
		field := requestFields.Get(i)
		msg := field.Message()
		if msg != nil {
			switch msg.FullName() {
			case "psm.list.v1.PageRequest":
				ll.pageRequestField = field
				continue
			case "psm.list.v1.QueryRequest":
				ll.queryRequestField = field
				continue
			default:
				return nil, fmt.Errorf("unknown field in request: '%s' of type %s", field.Name(), field.Kind())
			}
		}

		// Assume this is a filter field
		switch field.Kind() {
		case protoreflect.StringKind:
			ll.RequestFilterFields = append(ll.RequestFilterFields, field)
		case protoreflect.BoolKind:
			ll.RequestFilterFields = append(ll.RequestFilterFields, field)
		default:
			return nil, fmt.Errorf("unsupported filter field in request: '%s' of type %s", field.Name(), field.Kind())
		}

	}

	if ll.pageRequestField == nil {
		return nil, fmt.Errorf("no page field in request, %s must have a psm.list.v1.PageRequest", req.FullName())
	}

	if ll.queryRequestField == nil {
		return nil, fmt.Errorf("no query field in request, %s must have a psm.list.v1.QueryRequest", req.FullName())
	}

	arrayFieldOpt := ll.arrayField.Options().(*descriptorpb.FieldOptions)
	validateOpt := proto.GetExtension(arrayFieldOpt, validate.E_Field).(*validate.FieldConstraints)
	if repeated := validateOpt.GetRepeated(); repeated != nil {
		if repeated.MaxItems != nil {
			ll.defaultPageSize = *repeated.MaxItems
		}
	}

	return &ll, nil
}

type Lister[
	REQ ListRequest,
	RES ListResponse,
] struct {
	ListReflectionSet

	tableName  string
	dataColumn string
	auth       AuthProvider
	authJoin   []*LeftJoin

	requestFilter func(REQ) (map[string]interface{}, error)

	validator *protovalidate.Validator
}

func NewLister[
	REQ ListRequest,
	RES ListResponse,
](spec ListSpec[REQ, RES], options ...ListerOption) (*Lister[REQ, RES], error) {
	ll := &Lister[REQ, RES]{
		tableName:  spec.TableName,
		dataColumn: spec.DataColumn,
		auth:       spec.Auth,
		authJoin:   spec.AuthJoin,
	}

	descriptors := newMethodDescriptor[REQ, RES]()

	optionsStruct := resolveListerOptions(options)

	listFields, err := buildListReflection(descriptors.request, descriptors.response, optionsStruct)
	if err != nil {
		return nil, err
	}
	ll.ListReflectionSet = *listFields
	ll.requestFilter = spec.RequestFilter

	ll.validator, err = protovalidate.New()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize validator: %w", err)
	}

	return ll, nil
}

func buildTieBreakerFields(req protoreflect.MessageDescriptor, arrayField protoreflect.MessageDescriptor, fallback []string) ([]sortSpec, error) {
	listRequestAnnotation, ok := proto.GetExtension(req.Options().(*descriptorpb.MessageOptions), psml_pb.E_ListRequest).(*psml_pb.ListRequestMessage)
	if ok && listRequestAnnotation != nil && len(listRequestAnnotation.SortTiebreaker) > 0 {
		tieBreakerFields := make([]sortSpec, 0, len(listRequestAnnotation.SortTiebreaker))
		for _, tieBreaker := range listRequestAnnotation.SortTiebreaker {
			field, err := findField(arrayField, tieBreaker)
			if err != nil {
				return nil, err
			}
			tieBreakerFields = append(tieBreakerFields, sortSpec{
				nestedField: *field,
				desc:        false,
			})
		}
		return tieBreakerFields, nil
	}

	if len(fallback) == 0 {
		return []sortSpec{}, nil
	}

	tieBreakerFields := make([]sortSpec, 0, len(fallback))
	for _, tieBreaker := range fallback {
		field, err := findField(arrayField, tieBreaker)
		if err != nil {
			return nil, err
		}
		tieBreakerFields = append(tieBreakerFields, sortSpec{
			nestedField: *field,
			desc:        false,
		})
	}
	return tieBreakerFields, nil

}

func buildDefaultSorts(messageFields protoreflect.FieldDescriptors) []sortSpec {
	var defaultSortFields []sortSpec

	for i := 0; i < messageFields.Len(); i++ {
		field := messageFields.Get(i)
		fieldOpts := proto.GetExtension(field.Options().(*descriptorpb.FieldOptions), psml_pb.E_Field).(*psml_pb.FieldConstraint)

		if fieldOpts != nil {
			isDefaultSort := false

			switch fieldOps := fieldOpts.Type.(type) {
			case *psml_pb.FieldConstraint_Double:
				if fieldOps.Double.Sorting != nil && fieldOps.Double.Sorting.DefaultSort {
					isDefaultSort = true
				}
			case *psml_pb.FieldConstraint_Fixed32:
				if fieldOps.Fixed32.Sorting != nil && fieldOps.Fixed32.Sorting.DefaultSort {
					isDefaultSort = true
				}
			case *psml_pb.FieldConstraint_Fixed64:
				if fieldOps.Fixed64.Sorting != nil && fieldOps.Fixed64.Sorting.DefaultSort {
					isDefaultSort = true
				}
			case *psml_pb.FieldConstraint_Float:
				if fieldOps.Float.Sorting != nil && fieldOps.Float.Sorting.DefaultSort {
					isDefaultSort = true
				}
			case *psml_pb.FieldConstraint_Int32:
				if fieldOps.Int32.Sorting != nil && fieldOps.Int32.Sorting.DefaultSort {
					isDefaultSort = true
				}
			case *psml_pb.FieldConstraint_Int64:
				if fieldOps.Int64.Sorting != nil && fieldOps.Int64.Sorting.DefaultSort {
					isDefaultSort = true
				}
			case *psml_pb.FieldConstraint_Sfixed32:
				if fieldOps.Sfixed32.Sorting != nil && fieldOps.Sfixed32.Sorting.DefaultSort {
					isDefaultSort = true
				}
			case *psml_pb.FieldConstraint_Sfixed64:
				if fieldOps.Sfixed64.Sorting != nil && fieldOps.Sfixed64.Sorting.DefaultSort {
					isDefaultSort = true
				}
			case *psml_pb.FieldConstraint_Sint32:
				if fieldOps.Sint32.Sorting != nil && fieldOps.Sint32.Sorting.DefaultSort {
					isDefaultSort = true
				}
			case *psml_pb.FieldConstraint_Sint64:
				if fieldOps.Sint64.Sorting != nil && fieldOps.Sint64.Sorting.DefaultSort {
					isDefaultSort = true
				}
			case *psml_pb.FieldConstraint_Uint32:
				if fieldOps.Uint32.Sorting != nil && fieldOps.Uint32.Sorting.DefaultSort {
					isDefaultSort = true
				}
			case *psml_pb.FieldConstraint_Uint64:
				if fieldOps.Uint64.Sorting != nil && fieldOps.Uint64.Sorting.DefaultSort {
					isDefaultSort = true
				}
			case *psml_pb.FieldConstraint_Timestamp:
				if fieldOps.Timestamp.Sorting != nil && fieldOps.Timestamp.Sorting.DefaultSort {
					isDefaultSort = true
				}
			}
			if isDefaultSort {
				defaultSortFields = append(defaultSortFields, sortSpec{
					nestedField: nestedField{
						field:     field,
						fieldPath: []protoreflect.FieldDescriptor{field},
						jsonPath:  []string{field.JSONName()},
					},
					desc: true,
				})
			}
		} else if field.Kind() == protoreflect.MessageKind {
			subSort := buildDefaultSorts(field.Message().Fields())
			for idx, subSortField := range subSort {
				subSortField.jsonPath = append([]string{field.JSONName()}, subSortField.jsonPath...)
				subSortField.fieldPath = append([]protoreflect.FieldDescriptor{field}, subSortField.fieldPath...)
				subSort[idx] = subSortField
			}
			defaultSortFields = append(defaultSortFields, subSort...)
		}
	}

	return defaultSortFields
}

func (ll *Lister[REQ, RES]) getPageSize(req protoreflect.Message) (uint64, error) {
	pageSize := ll.defaultPageSize

	pageReq, ok := req.Get(ll.pageRequestField).Message().Interface().(*psml_pb.PageRequest)
	if ok && pageReq != nil && pageReq.PageSize != nil {
		pageSize = uint64(*pageReq.PageSize)

		if pageSize > ll.defaultPageSize {
			return 0, fmt.Errorf("page size exceeds the maximum allowed size of %d", ll.defaultPageSize)
		}
	}

	return pageSize, nil
}

func (ll *Lister[REQ, RES]) List(ctx context.Context, db Transactor, reqMsg proto.Message, resMsg proto.Message) error {
	if err := ll.validator.Validate(reqMsg); err != nil {
		return fmt.Errorf("validating request %s: %w", reqMsg.ProtoReflect().Descriptor().FullName(), err)
	}

	res := resMsg.ProtoReflect()
	req := reqMsg.ProtoReflect()

	pageSize, err := ll.getPageSize(req)
	if err != nil {
		return err
	}

	selectQuery, err := ll.BuildQuery(ctx, req, res)
	if err != nil {
		return err
	}

	var jsonRows = make([][]byte, 0, pageSize)
	if err := db.Transact(ctx, &sqrlx.TxOptions{
		ReadOnly:  true,
		Retryable: true,
		Isolation: sql.LevelReadCommitted,
	}, func(ctx context.Context, tx sqrlx.Transaction) error {
		rows, err := tx.Query(ctx, selectQuery)
		if err != nil {
			return fmt.Errorf("run select: %w", err)
		}
		defer rows.Close()

		for rows.Next() {
			var json []byte
			if err := rows.Scan(&json); err != nil {
				return err
			}
			jsonRows = append(jsonRows, json)
		}
		return rows.Err()
	}); err != nil {
		stmt, _, _ := selectQuery.ToSql()
		log.WithField(ctx, "query", stmt).Error("list query")
		return fmt.Errorf("list query: %w", err)
	}

	list := res.Mutable(ll.arrayField).List()
	res.Set(ll.arrayField, protoreflect.ValueOf(list))

	var nextToken string
	for idx, rowBytes := range jsonRows {
		rowMessage := list.NewElement().Message()
		if err := protojson.Unmarshal(rowBytes, rowMessage.Interface()); err != nil {
			return fmt.Errorf("unmarshal into %s from %s: %w", rowMessage.Descriptor().FullName(), string(rowBytes), err)
		}
		if idx >= int(pageSize) {
			// TODO: This works but the token is huge.
			// The eventual solution will need to look at
			// the sorting and filtering of the query and either encode them
			// directly, or encode a subset of the message as required.
			lastBytes, err := proto.Marshal(rowMessage.Interface())
			if err != nil {
				return fmt.Errorf("marshalling final row: %w", err)
			}
			nextToken = base64.StdEncoding.EncodeToString(lastBytes)
			break
		}
		list.Append(protoreflect.ValueOf(rowMessage))
	}

	pageResponse := &psml_pb.PageResponse{}
	if nextToken != "" {
		pageResponse.NextToken = &nextToken
	}

	res.Set(ll.pageResponseField, protoreflect.ValueOf(pageResponse.ProtoReflect()))

	return nil
}

func (ll *Lister[REQ, RES]) BuildQuery(ctx context.Context, req protoreflect.Message, res protoreflect.Message) (*sqrl.SelectBuilder, error) {
	as := newAliasSet()
	tableAlias := as.Next()

	selectQuery := sq.
		Select(fmt.Sprintf("%s.%s", tableAlias, ll.dataColumn)).
		From(fmt.Sprintf("%s AS %s", ll.tableName, tableAlias))

	sortFields := ll.defaultSortFields
	// TODO: Dynamic Sorts

	sortFields = append(sortFields, ll.tieBreakerFields...)

	for _, sortField := range sortFields {
		direction := "ASC"
		if sortField.desc {
			direction = "DESC"
		}
		selectQuery.OrderBy(fmt.Sprintf("%s.%s%s %s", tableAlias, ll.dataColumn, sortField.jsonbPath(), direction))
	}

	if ll.requestFilter != nil {
		filter, err := ll.requestFilter(req.Interface().(REQ))
		if err != nil {
			return nil, err
		}

		if len(filter) > 0 {

			filterMapped, err := dbconvert.FieldsToEqMap(tableAlias, filter)
			if err != nil {
				return nil, err
			}

			selectQuery.Where(filterMapped)
		}
	}

	if ll.auth != nil {
		authFilter, err := ll.auth.AuthFilter(ctx)
		if err != nil {
			return nil, err
		}
		authAlias := tableAlias

		for _, join := range ll.authJoin {
			priorAlias := authAlias
			authAlias = as.Next()
			selectQuery = selectQuery.LeftJoin(fmt.Sprintf(
				"%s AS %s ON %s",
				join.TableName,
				authAlias,
				join.On.SQL(priorAlias, authAlias),
			))
		}

		authFilterMapped, err := dbconvert.FieldsToEqMap(authAlias, authFilter)
		if err != nil {
			return nil, err
		}

		selectQuery = selectQuery.Where(authFilterMapped)
	}

	pageSize, err := ll.getPageSize(req)
	if err != nil {
		return nil, err
	}

	selectQuery.Limit(pageSize + 1)

	// TODO: Request Filters req := reqMsg.ProtoReflect()

	reqPage, ok := req.Get(ll.pageRequestField).Message().Interface().(*psml_pb.PageRequest)
	if ok && reqPage != nil && reqPage.GetToken() != "" {
		rowMessage := dynamicpb.NewMessage(ll.arrayField.Message())

		rowBytes, err := base64.StdEncoding.DecodeString(reqPage.GetToken())
		if err != nil {
			return nil, fmt.Errorf("decode token: %w", err)
		}

		if err := proto.Unmarshal(rowBytes, rowMessage.Interface()); err != nil {
			return nil, fmt.Errorf("unmarshal into %s from %s: %w", rowMessage.Descriptor().FullName(), string(rowBytes), err)
		}

		lhsFields := make([]string, 0, len(sortFields))
		rhsValues := make([]interface{}, 0, len(sortFields))
		rhsPlaceholders := make([]string, 0, len(sortFields))

		for _, sortField := range sortFields {

			rowSelecter := fmt.Sprintf("%s.%s%s",
				tableAlias,
				ll.dataColumn,
				sortField.jsonbPath(),
			)
			valuePlaceholder := "?"

			fieldVal, err := walkProtoValue(rowMessage, sortField.fieldPath)
			if err != nil {
				return nil, fmt.Errorf("sort field %s: %w", strings.Join(sortField.jsonPath, "."), err)
			}
			dbVal := fieldVal.Interface()
			switch subType := dbVal.(type) {
			case *dynamicpb.Message:
				name := subType.Descriptor().FullName()
				msgBytes, err := proto.Marshal(subType)
				if err != nil {
					return nil, fmt.Errorf("marshal %s: %w", name, err)
				}

				switch name {
				case "google.protobuf.Timestamp":
					ts := timestamppb.Timestamp{}
					if err := proto.Unmarshal(msgBytes, &ts); err != nil {
						return nil, fmt.Errorf("unmarshal %s: %w", name, err)
					}
					intVal := ts.AsTime().Round(time.Microsecond).UnixMicro()
					// Go rounds half-up.
					// Postgres is undocumented, but can only handle
					// microseconds.
					rowSelecter = fmt.Sprintf("(EXTRACT(epoch FROM (%s)::timestamp) * 1000000)::bigint", rowSelecter)
					if sortField.desc {
						intVal = intVal * -1
						rowSelecter = fmt.Sprintf("-1 * %s", rowSelecter)
					}
					dbVal = intVal
				default:
					return nil, fmt.Errorf("sort field %s is a message of type %s", sortField.field.Name(), name)
				}

			case string:
				dbVal = subType
				if sortField.desc {
					// String fields aren't valid for sorting in listify, they
					// can only be used for the tie-breaker so the order itself
					// is not important, only that it is consistent
					return nil, fmt.Errorf("sort field %s is a string, strings cannot be sorted DESC", sortField.field.Name())
				}

			default:
				return nil, fmt.Errorf("unknown sort field %s, type %T", strings.Join(sortField.jsonPath, "."), dbVal)
			}

			lhsFields = append(lhsFields, rowSelecter)
			rhsValues = append(rhsValues, dbVal)
			rhsPlaceholders = append(rhsPlaceholders, valuePlaceholder)

		}

		// From https://www.postgresql.org/docs/current/functions-comparisons.html#ROW-WISE-COMPARISON
		// >> for the <, <=, > and >= cases, the row elements are compared left-to-right, stopping as soon
		// >> as an unequal or null pair of elements is found. If either of this pair of elements is null,
		// >> the result of the row comparison is unknown (null); otherwise comparison of this pair of elements
		// >> determines the result. For example, ROW(1,2,NULL) < ROW(1,3,0) yields true, not null, because the
		// >> third pair of elements are not considered.
		//
		// This means that we can use the row comparison with the same fields as
		// the sort fields to exclude the exact record we want, rather than the
		// filter being applied to all fields equally which takes out valid
		// records.
		// `(1, 30) >= (1, 20)` is true, so is `1 >= 1 AND 30 >= 20`
		// `(2, 10) >= (1, 20)` is also true, but `2 >= 1 AND 10 >= 20` is false
		// Since the tuple comparrison starts from the left and stops at the first term.
		//
		// The downside is that we have to negate the values to sort in reverse
		// order, as we don't get an operator per term. This gets strange for
		// some data types and will create some crazy indexes.
		//
		// TODO: Optimise the cases when the order is ASC and therefore we don't
		// need to flip, but also the cases where we can just reverse the whole
		// comparrison and reverse all flips to simplify, noting again that it
		// does not actually matter in which order the string field is sorted...
		// or don't because indexes.

		selectQuery = selectQuery.Where(
			fmt.Sprintf("(%s) >= (%s)",
				strings.Join(lhsFields, ","),
				strings.Join(rhsPlaceholders, ","),
			), rhsValues...)
	}

	return selectQuery, nil
}

type nestedField struct {
	jsonPath  []string
	fieldPath []protoreflect.FieldDescriptor
	field     protoreflect.FieldDescriptor
}

func walkProtoValue(msg protoreflect.Message, path []protoreflect.FieldDescriptor) (protoreflect.Value, error) {
	if len(path) == 0 {
		return protoreflect.Value{}, fmt.Errorf("empty path")
	}
	var val protoreflect.Value
	field := path[0]

	if !msg.Has(field) {
		return protoreflect.Value{}, fmt.Errorf("field %s not found in %s", field.Name(), msg.Descriptor().FullName())
	}
	val = msg.Get(field)
	if len(path) == 1 {
		return val, nil
	}
	if field.Kind() != protoreflect.MessageKind {
		return protoreflect.Value{}, fmt.Errorf("field %s is not a message", field.Name())
	}

	return walkProtoValue(val.Message(), path[1:])
}

func (nf nestedField) jsonbPath() string {
	out := strings.Builder{}
	last := len(nf.jsonPath) - 1
	for idx, part := range nf.jsonPath {
		// last part gets a double >
		if idx == last {
			out.WriteString("->>")
		} else {
			out.WriteString("->")
		}
		out.WriteString(fmt.Sprintf("'%s'", part))
	}

	return out.String()
}

func findField(message protoreflect.MessageDescriptor, path string) (*nestedField, error) {
	var fieldName protoreflect.Name
	var furtherPath string
	parts := strings.SplitN(path, ".", 2)
	if len(parts) == 2 {
		fieldName = protoreflect.Name(parts[0])
		furtherPath = parts[1]
	} else {
		fieldName = protoreflect.Name(path)
	}

	field := message.Fields().ByName(fieldName)
	if field == nil {
		return nil, fmt.Errorf("no field named '%s' in %s", fieldName, message.FullName())
	}

	if furtherPath == "" {
		return &nestedField{
			field:     field,
			fieldPath: []protoreflect.FieldDescriptor{field},
			jsonPath:  []string{field.JSONName()},
		}, nil
	}

	if field.Kind() != protoreflect.MessageKind {
		return nil, fmt.Errorf("field %s is not a message", fieldName)
	}
	spec, err := findField(field.Message(), furtherPath)
	if err != nil {
		return nil, fmt.Errorf("field %s: %w", parts[0], err)
	}
	return &nestedField{
		field:     spec.field,
		fieldPath: append([]protoreflect.FieldDescriptor{field}, spec.fieldPath...),
		jsonPath:  append([]string{field.JSONName()}, spec.jsonPath...),
	}, nil

}
