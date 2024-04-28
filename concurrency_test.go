package main

import (
	"errors"
	"github.com/hallgren/eventsourcing"
	"github.com/hallgren/eventsourcing/eventstore/memory"
	"log"
	"sync"
	"testing"
)

func TestConcurrencySaves(t *testing.T) {
	repo := eventsourcing.NewEventRepository(memory.Create())

	repo.Register(&Person{})

	person, _ := CreatePerson("Alice")
	person.GrowOlder()

	for i := 0; i < 1000; i++ {
		repo.Save(person)
	}

	loops := 5
	wg := sync.WaitGroup{}
	wg.Add(loops)
	for i := 0; i < loops; i++ {
		go func(localCounter int) {
			wg.Add(1)
			defer wg.Done()
			if err := repo.Save(person); err != nil {
				t.Errorf("failed to save event: %v", err)
				return
			}
			log.Println("saved event")
		}(i)
	}
	wg.Wait()

	twin := Person{}
	repo.Get(person.ID(), &twin)
}

type Person struct {
	eventsourcing.AggregateRoot
	Name string
	Age  int
}

// GrowOlder command
func (person *Person) GrowOlder() {
	person.TrackChange(person, &AgedOneYear{})
}

// Transition the person state dependent on the events
func (person *Person) Transition(event eventsourcing.Event) {
	switch e := event.Data().(type) {
	case *Born:
		person.Age = 0
		person.Name = e.Name
	case *AgedOneYear:
		person.Age += 1
	}
}

// Register callback method that register Person events to the repository
func (person *Person) Register(r eventsourcing.RegisterFunc) {
	r(&Born{}, &AgedOneYear{})
}

// Initial event
type Born struct {
	Name string
}

// Event that happens once a year
type AgedOneYear struct{}

// CreatePerson constructor for Person
func CreatePerson(name string) (*Person, error) {
	if name == "" {
		return nil, errors.New("name can't be blank")
	}
	person := Person{}
	person.TrackChange(&person, &Born{Name: name})
	return &person, nil
}
