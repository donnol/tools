package slicex

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/samber/lo"
)

type User struct {
	Id   uint64
	Name string
}

type Book struct {
	Id     uint64
	Title  string
	Author uint64 // User.Id
}

type BookWithUser struct {
	Book
	UserName string
}

func UserBookMatcher(book Book, user User) bool {
	return book.Author == user.Id
}

var (
	users = []User{
		{Id: 1, Name: "jd"},
		{Id: 2, Name: "jc"},
	}
	books = []Book{
		{Id: 1, Title: "hello", Author: 1},
		{Id: 2, Title: "world", Author: 1},
		{Id: 3, Title: "good", Author: 2},
		{Id: 4, Title: "job", Author: 2},
	}
)

func TestJoin(t *testing.T) {
	r := Join(books, users, UserBookMatcher, func(book Book, user User) BookWithUser {
		return BookWithUser{
			Book:     book,
			UserName: user.Name,
		}
	})
	want := []BookWithUser{
		{Book{1, "hello", 1}, "jd"},
		{Book{2, "world", 1}, "jd"},
		{Book{3, "good", 2}, "jc"},
		{Book{4, "job", 2}, "jc"},
	}
	if !reflect.DeepEqual(r, want) {
		t.Errorf("bad case, %+v != %+v", r, want)
	}
}

func TestJoinByKey(t *testing.T) {
	r := JoinByKey(books, users, func(item Book) uint64 {
		return item.Author
	}, func(item User) uint64 {
		return item.Id
	}, func(book Book, user User) BookWithUser {
		return BookWithUser{
			Book:     book,
			UserName: user.Name,
		}
	})
	want := []BookWithUser{
		{Book{1, "hello", 1}, "jd"},
		{Book{2, "world", 1}, "jd"},
		{Book{3, "good", 2}, "jc"},
		{Book{4, "job", 2}, "jc"},
	}
	if !reflect.DeepEqual(r, want) {
		t.Errorf("bad case, %+v != %+v", r, want)
	}
}

func TestJoinByKey_CompositeKey(t *testing.T) {
	type CompositeKey string

	r := JoinByKey(books, users, func(item Book) CompositeKey {
		return CompositeKey(fmt.Sprintf("%d|%d", item.Author, item.Author))
	}, func(item User) CompositeKey {
		return CompositeKey(fmt.Sprintf("%d|%d", item.Id, item.Id))
	}, func(book Book, user User) BookWithUser {
		return BookWithUser{
			Book:     book,
			UserName: user.Name,
		}
	})
	want := []BookWithUser{
		{Book{1, "hello", 1}, "jd"},
		{Book{2, "world", 1}, "jd"},
		{Book{3, "good", 2}, "jc"},
		{Book{4, "job", 2}, "jc"},
	}
	if !reflect.DeepEqual(r, want) {
		t.Errorf("bad case, %+v != %+v", r, want)
	}
}

func BenchmarkJoin(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Join(books, users, UserBookMatcher, func(book Book, user User) BookWithUser {
			return BookWithUser{
				Book:     book,
				UserName: user.Name,
			}
		})
	}
}

func BenchmarkJoin100X(b *testing.B) {
	users := lo.RepeatBy(100, func(index int) User {
		return users[1]
	})
	books := lo.RepeatBy(100, func(index int) Book {
		return books[3]
	})

	for i := 0; i < b.N; i++ {
		Join(books, users, UserBookMatcher, func(book Book, user User) BookWithUser {
			return BookWithUser{
				Book:     book,
				UserName: user.Name,
			}
		})
	}
}

func BenchmarkJoinByKey(b *testing.B) {
	JoinByKey(books, users, func(item Book) uint64 {
		return item.Author
	}, func(item User) uint64 {
		return item.Id
	}, func(book Book, user User) BookWithUser {
		return BookWithUser{
			Book:     book,
			UserName: user.Name,
		}
	})
}

var (
	usersX1M = lo.RepeatBy(1_000_000, func(index int) User {
		return users[1]
	})
	booksX1M = lo.RepeatBy(1_000_000, func(index int) Book {
		return books[3]
	})
)

func BenchmarkJoinByKey100X(b *testing.B) {
	JoinByKey(booksX1M, usersX1M, func(item Book) uint64 {
		return item.Author
	}, func(item User) uint64 {
		return item.Id
	}, func(book Book, user User) BookWithUser {
		return BookWithUser{
			Book:     book,
			UserName: user.Name,
		}
	})
}
