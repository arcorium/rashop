package types

import "github.com/arcorium/rashop/shared/errors"

type Aggregate interface {
  Entity
  // AddEvents register event into aggregate
  AddEvents(events ...Event)
  // ApplyEvents will apply func parameter into each event
  ApplyEvents(events ...Event) errors.IndicesError
  // ApplyEvent will apply the event parameter into the object (object will be mutated)
  ApplyEvent(ev Event) error
  // Events get all registered events
  Events() []Event
}

type AggregateBase struct {
  Aggregate
  events []Event
}

func (a AggregateBase) AddEvents(events ...Event) {
  a.events = append(a.events, events...)
}

func (a AggregateBase) ApplyEvents(events ...Event) errors.IndicesError {
  var errs errors.IndicesError
  for i, event := range events {
    err := a.ApplyEvent(event)
    if err != nil {
      errs.Append(errors.NewIndex(i, err))
    }
  }
  return errs
}

func (a AggregateBase) Events() []Event {
  return a.events
}

// ChildEntityHelper to create helper for aggregate to determine what should be done for the child entities
type ChildEntityHelper[ID comparable] struct {
  updatedIndices Set[int]
  deletedIndices Set[ID]
  totalAdded     uint64
}

func (a *ChildEntityHelper[ID]) Update(idx int) {
  a.updatedIndices.Add(idx)
}

func (a *ChildEntityHelper[ID]) Updated() []int {
  return a.updatedIndices.Values()
}

func (a *ChildEntityHelper[ID]) Delete(idx ID) {
  a.deletedIndices.Add(idx)
}

func (a *ChildEntityHelper[ID]) Deleted() []ID {
  return a.deletedIndices.Values()
}

// Add delta into total added, when the delta is empty it will use 1
func (a *ChildEntityHelper[ID]) Add(delta ...int32) {
  var modifier int32 = 1
  if len(delta) > 0 {
    modifier = 0
    for _, dt := range delta {
      modifier += dt
    }
  }
  a.totalAdded += uint64(modifier)
}

func (a *ChildEntityHelper[ID]) Added() uint64 {
  return a.totalAdded
}

func NewChildEntityHelper[ID comparable, T any](val []T) ChildEntityHelperWithObject[ID, T] {
  return ChildEntityHelperWithObject[ID, T]{
    ChildEntityHelper: ChildEntityHelper[ID]{},
    Elm:               val,
  }
}

type ChildEntityHelperWithObject[ID comparable, T any] struct {
  ChildEntityHelper[ID]
  Elm []T // Embedded slice objects
}

func (e *ChildEntityHelperWithObject[ID, T]) HasElement() bool {
  return len(e.Elm) > 0
}

// Elements helper method for person that more like to use function instead of field directly
func (e *ChildEntityHelperWithObject[ID, T]) Elements() []T {
  return e.Elm
}

type AggregateHelper struct {
  // All values defaulted to false
  deleted bool
  updated bool
  created bool
}

func (a *AggregateHelper) Created() bool {
  return a.created
}

func (a *AggregateHelper) MarkCreated() {
  a.created = true
}

func (a *AggregateHelper) Deleted() bool {
  return a.deleted
}

func (a *AggregateHelper) MarkDeleted() {
  a.deleted = true
}

func (a *AggregateHelper) Updated() bool {
  return a.updated
}

func (a *AggregateHelper) MarkUpdated() {
  a.updated = true
}
