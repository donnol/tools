package worker

import (
	"log"
	"testing"
	"time"
)

func TestWorker(t *testing.T) {
	w := New(10)
	w.Start()

	if err := w.Push(Job{
		do: func() error {
			panic("terriable")
		},
	}); err != nil {
		t.Fatal(err)
	}
	if err := w.Push(Job{
		do: func() error {
			for i := 0; i < 10; i++ {
				log.Printf("i: %d\n", i)
			}
			return nil
		},
	}); err != nil {
		t.Fatal(err)
	}
	for i := 10; i < 1000; i++ {
		tmp := i
		if err := w.Push(Job{
			do: func() error {
				log.Printf("i: %d\n", tmp)
				return nil
			},
		}); err != nil {
			t.Fatal(err)
		}
	}
	if err := w.Push(Job{}); err != ErrNilJobDo {
		t.Fatal(err)
	}

	w.Stop()

	if err := w.Push(Job{
		do: func() error {
			log.Printf("Push after stop.")
			return nil
		},
	}); err != ErrWorkerIsStop {
		t.Fatal(err)
	}
}

func TestWorkerWithTimeout(t *testing.T) {
	w := New(0)
	w.Start()

	job := MakeJob(func() error {
		for i := 0; i < 10; i++ {
			log.Printf("i: %d\n", i)
			time.Sleep(1 * time.Second)
		}
		return nil
	}, 5*time.Second, func(err error) {
		log.Printf("err is %v\n", err)
	})
	if err := w.Push(job); err != nil {
		t.Fatal(err)
	}

	w.Stop()
}
