// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"database/sql/driver"
	"fmt"
	"math"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/safedep/vet/ent/codesourcefile"
	"github.com/safedep/vet/ent/depsusageevidence"
	"github.com/safedep/vet/ent/predicate"
)

// CodeSourceFileQuery is the builder for querying CodeSourceFile entities.
type CodeSourceFileQuery struct {
	config
	ctx                    *QueryContext
	order                  []codesourcefile.OrderOption
	inters                 []Interceptor
	predicates             []predicate.CodeSourceFile
	withDepsUsageEvidences *DepsUsageEvidenceQuery
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the CodeSourceFileQuery builder.
func (csfq *CodeSourceFileQuery) Where(ps ...predicate.CodeSourceFile) *CodeSourceFileQuery {
	csfq.predicates = append(csfq.predicates, ps...)
	return csfq
}

// Limit the number of records to be returned by this query.
func (csfq *CodeSourceFileQuery) Limit(limit int) *CodeSourceFileQuery {
	csfq.ctx.Limit = &limit
	return csfq
}

// Offset to start from.
func (csfq *CodeSourceFileQuery) Offset(offset int) *CodeSourceFileQuery {
	csfq.ctx.Offset = &offset
	return csfq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (csfq *CodeSourceFileQuery) Unique(unique bool) *CodeSourceFileQuery {
	csfq.ctx.Unique = &unique
	return csfq
}

// Order specifies how the records should be ordered.
func (csfq *CodeSourceFileQuery) Order(o ...codesourcefile.OrderOption) *CodeSourceFileQuery {
	csfq.order = append(csfq.order, o...)
	return csfq
}

// QueryDepsUsageEvidences chains the current query on the "deps_usage_evidences" edge.
func (csfq *CodeSourceFileQuery) QueryDepsUsageEvidences() *DepsUsageEvidenceQuery {
	query := (&DepsUsageEvidenceClient{config: csfq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := csfq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := csfq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(codesourcefile.Table, codesourcefile.FieldID, selector),
			sqlgraph.To(depsusageevidence.Table, depsusageevidence.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, true, codesourcefile.DepsUsageEvidencesTable, codesourcefile.DepsUsageEvidencesColumn),
		)
		fromU = sqlgraph.SetNeighbors(csfq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first CodeSourceFile entity from the query.
// Returns a *NotFoundError when no CodeSourceFile was found.
func (csfq *CodeSourceFileQuery) First(ctx context.Context) (*CodeSourceFile, error) {
	nodes, err := csfq.Limit(1).All(setContextOp(ctx, csfq.ctx, ent.OpQueryFirst))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{codesourcefile.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (csfq *CodeSourceFileQuery) FirstX(ctx context.Context) *CodeSourceFile {
	node, err := csfq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first CodeSourceFile ID from the query.
// Returns a *NotFoundError when no CodeSourceFile ID was found.
func (csfq *CodeSourceFileQuery) FirstID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = csfq.Limit(1).IDs(setContextOp(ctx, csfq.ctx, ent.OpQueryFirstID)); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{codesourcefile.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (csfq *CodeSourceFileQuery) FirstIDX(ctx context.Context) int {
	id, err := csfq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single CodeSourceFile entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one CodeSourceFile entity is found.
// Returns a *NotFoundError when no CodeSourceFile entities are found.
func (csfq *CodeSourceFileQuery) Only(ctx context.Context) (*CodeSourceFile, error) {
	nodes, err := csfq.Limit(2).All(setContextOp(ctx, csfq.ctx, ent.OpQueryOnly))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{codesourcefile.Label}
	default:
		return nil, &NotSingularError{codesourcefile.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (csfq *CodeSourceFileQuery) OnlyX(ctx context.Context) *CodeSourceFile {
	node, err := csfq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only CodeSourceFile ID in the query.
// Returns a *NotSingularError when more than one CodeSourceFile ID is found.
// Returns a *NotFoundError when no entities are found.
func (csfq *CodeSourceFileQuery) OnlyID(ctx context.Context) (id int, err error) {
	var ids []int
	if ids, err = csfq.Limit(2).IDs(setContextOp(ctx, csfq.ctx, ent.OpQueryOnlyID)); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{codesourcefile.Label}
	default:
		err = &NotSingularError{codesourcefile.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (csfq *CodeSourceFileQuery) OnlyIDX(ctx context.Context) int {
	id, err := csfq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of CodeSourceFiles.
func (csfq *CodeSourceFileQuery) All(ctx context.Context) ([]*CodeSourceFile, error) {
	ctx = setContextOp(ctx, csfq.ctx, ent.OpQueryAll)
	if err := csfq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*CodeSourceFile, *CodeSourceFileQuery]()
	return withInterceptors[[]*CodeSourceFile](ctx, csfq, qr, csfq.inters)
}

// AllX is like All, but panics if an error occurs.
func (csfq *CodeSourceFileQuery) AllX(ctx context.Context) []*CodeSourceFile {
	nodes, err := csfq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of CodeSourceFile IDs.
func (csfq *CodeSourceFileQuery) IDs(ctx context.Context) (ids []int, err error) {
	if csfq.ctx.Unique == nil && csfq.path != nil {
		csfq.Unique(true)
	}
	ctx = setContextOp(ctx, csfq.ctx, ent.OpQueryIDs)
	if err = csfq.Select(codesourcefile.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (csfq *CodeSourceFileQuery) IDsX(ctx context.Context) []int {
	ids, err := csfq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (csfq *CodeSourceFileQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, csfq.ctx, ent.OpQueryCount)
	if err := csfq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, csfq, querierCount[*CodeSourceFileQuery](), csfq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (csfq *CodeSourceFileQuery) CountX(ctx context.Context) int {
	count, err := csfq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (csfq *CodeSourceFileQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, csfq.ctx, ent.OpQueryExist)
	switch _, err := csfq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (csfq *CodeSourceFileQuery) ExistX(ctx context.Context) bool {
	exist, err := csfq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the CodeSourceFileQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (csfq *CodeSourceFileQuery) Clone() *CodeSourceFileQuery {
	if csfq == nil {
		return nil
	}
	return &CodeSourceFileQuery{
		config:                 csfq.config,
		ctx:                    csfq.ctx.Clone(),
		order:                  append([]codesourcefile.OrderOption{}, csfq.order...),
		inters:                 append([]Interceptor{}, csfq.inters...),
		predicates:             append([]predicate.CodeSourceFile{}, csfq.predicates...),
		withDepsUsageEvidences: csfq.withDepsUsageEvidences.Clone(),
		// clone intermediate query.
		sql:  csfq.sql.Clone(),
		path: csfq.path,
	}
}

// WithDepsUsageEvidences tells the query-builder to eager-load the nodes that are connected to
// the "deps_usage_evidences" edge. The optional arguments are used to configure the query builder of the edge.
func (csfq *CodeSourceFileQuery) WithDepsUsageEvidences(opts ...func(*DepsUsageEvidenceQuery)) *CodeSourceFileQuery {
	query := (&DepsUsageEvidenceClient{config: csfq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	csfq.withDepsUsageEvidences = query
	return csfq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		Path string `json:"path,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.CodeSourceFile.Query().
//		GroupBy(codesourcefile.FieldPath).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (csfq *CodeSourceFileQuery) GroupBy(field string, fields ...string) *CodeSourceFileGroupBy {
	csfq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &CodeSourceFileGroupBy{build: csfq}
	grbuild.flds = &csfq.ctx.Fields
	grbuild.label = codesourcefile.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		Path string `json:"path,omitempty"`
//	}
//
//	client.CodeSourceFile.Query().
//		Select(codesourcefile.FieldPath).
//		Scan(ctx, &v)
func (csfq *CodeSourceFileQuery) Select(fields ...string) *CodeSourceFileSelect {
	csfq.ctx.Fields = append(csfq.ctx.Fields, fields...)
	sbuild := &CodeSourceFileSelect{CodeSourceFileQuery: csfq}
	sbuild.label = codesourcefile.Label
	sbuild.flds, sbuild.scan = &csfq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a CodeSourceFileSelect configured with the given aggregations.
func (csfq *CodeSourceFileQuery) Aggregate(fns ...AggregateFunc) *CodeSourceFileSelect {
	return csfq.Select().Aggregate(fns...)
}

func (csfq *CodeSourceFileQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range csfq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, csfq); err != nil {
				return err
			}
		}
	}
	for _, f := range csfq.ctx.Fields {
		if !codesourcefile.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if csfq.path != nil {
		prev, err := csfq.path(ctx)
		if err != nil {
			return err
		}
		csfq.sql = prev
	}
	return nil
}

func (csfq *CodeSourceFileQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*CodeSourceFile, error) {
	var (
		nodes       = []*CodeSourceFile{}
		_spec       = csfq.querySpec()
		loadedTypes = [1]bool{
			csfq.withDepsUsageEvidences != nil,
		}
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*CodeSourceFile).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &CodeSourceFile{config: csfq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, csfq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := csfq.withDepsUsageEvidences; query != nil {
		if err := csfq.loadDepsUsageEvidences(ctx, query, nodes,
			func(n *CodeSourceFile) { n.Edges.DepsUsageEvidences = []*DepsUsageEvidence{} },
			func(n *CodeSourceFile, e *DepsUsageEvidence) {
				n.Edges.DepsUsageEvidences = append(n.Edges.DepsUsageEvidences, e)
			}); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (csfq *CodeSourceFileQuery) loadDepsUsageEvidences(ctx context.Context, query *DepsUsageEvidenceQuery, nodes []*CodeSourceFile, init func(*CodeSourceFile), assign func(*CodeSourceFile, *DepsUsageEvidence)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[int]*CodeSourceFile)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
		if init != nil {
			init(nodes[i])
		}
	}
	query.withFKs = true
	query.Where(predicate.DepsUsageEvidence(func(s *sql.Selector) {
		s.Where(sql.InValues(s.C(codesourcefile.DepsUsageEvidencesColumn), fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.deps_usage_evidence_used_in
		if fk == nil {
			return fmt.Errorf(`foreign-key "deps_usage_evidence_used_in" is nil for node %v`, n.ID)
		}
		node, ok := nodeids[*fk]
		if !ok {
			return fmt.Errorf(`unexpected referenced foreign-key "deps_usage_evidence_used_in" returned %v for node %v`, *fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}

func (csfq *CodeSourceFileQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := csfq.querySpec()
	_spec.Node.Columns = csfq.ctx.Fields
	if len(csfq.ctx.Fields) > 0 {
		_spec.Unique = csfq.ctx.Unique != nil && *csfq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, csfq.driver, _spec)
}

func (csfq *CodeSourceFileQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(codesourcefile.Table, codesourcefile.Columns, sqlgraph.NewFieldSpec(codesourcefile.FieldID, field.TypeInt))
	_spec.From = csfq.sql
	if unique := csfq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if csfq.path != nil {
		_spec.Unique = true
	}
	if fields := csfq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, codesourcefile.FieldID)
		for i := range fields {
			if fields[i] != codesourcefile.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := csfq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := csfq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := csfq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := csfq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (csfq *CodeSourceFileQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(csfq.driver.Dialect())
	t1 := builder.Table(codesourcefile.Table)
	columns := csfq.ctx.Fields
	if len(columns) == 0 {
		columns = codesourcefile.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if csfq.sql != nil {
		selector = csfq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if csfq.ctx.Unique != nil && *csfq.ctx.Unique {
		selector.Distinct()
	}
	for _, p := range csfq.predicates {
		p(selector)
	}
	for _, p := range csfq.order {
		p(selector)
	}
	if offset := csfq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := csfq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// CodeSourceFileGroupBy is the group-by builder for CodeSourceFile entities.
type CodeSourceFileGroupBy struct {
	selector
	build *CodeSourceFileQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (csfgb *CodeSourceFileGroupBy) Aggregate(fns ...AggregateFunc) *CodeSourceFileGroupBy {
	csfgb.fns = append(csfgb.fns, fns...)
	return csfgb
}

// Scan applies the selector query and scans the result into the given value.
func (csfgb *CodeSourceFileGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, csfgb.build.ctx, ent.OpQueryGroupBy)
	if err := csfgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*CodeSourceFileQuery, *CodeSourceFileGroupBy](ctx, csfgb.build, csfgb, csfgb.build.inters, v)
}

func (csfgb *CodeSourceFileGroupBy) sqlScan(ctx context.Context, root *CodeSourceFileQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(csfgb.fns))
	for _, fn := range csfgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*csfgb.flds)+len(csfgb.fns))
		for _, f := range *csfgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*csfgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := csfgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// CodeSourceFileSelect is the builder for selecting fields of CodeSourceFile entities.
type CodeSourceFileSelect struct {
	*CodeSourceFileQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (csfs *CodeSourceFileSelect) Aggregate(fns ...AggregateFunc) *CodeSourceFileSelect {
	csfs.fns = append(csfs.fns, fns...)
	return csfs
}

// Scan applies the selector query and scans the result into the given value.
func (csfs *CodeSourceFileSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, csfs.ctx, ent.OpQuerySelect)
	if err := csfs.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*CodeSourceFileQuery, *CodeSourceFileSelect](ctx, csfs.CodeSourceFileQuery, csfs, csfs.inters, v)
}

func (csfs *CodeSourceFileSelect) sqlScan(ctx context.Context, root *CodeSourceFileQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(csfs.fns))
	for _, fn := range csfs.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*csfs.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := csfs.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
