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

func UserBookCompare(book Book, user User) int {
	if book.Author == user.Id {
		return 0
	} else if book.Author > user.Id {
		return 1
	} else {
		return -1
	}
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

func TestJoinBNL(t *testing.T) {
	for _, tt := range []struct {
		name string
		size int
	}{
		{
			name: "1",
			size: 1,
		},
		{
			name: "2",
			size: 2,
		},
		{
			name: "3",
			size: 3,
		},
		{
			name: "4",
			size: 4,
		},
		{
			name: "5",
			size: 5,
		},
	} {
		setJoinBNLBufSize(tt.size)

		r := joinBNL(books, users, UserBookMatcher, func(book Book, user User) BookWithUser {
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
}

func TestJoinSortMerge(t *testing.T) {
	r := joinSortMerge(books, users, func(i, j int) bool {
		return books[i].Author < books[j].Author
	}, func(i, j int) bool {
		return users[i].Id < users[j].Id
	}, UserBookCompare, func(book Book, user User) BookWithUser {
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

func BenchmarkJoinBNL(b *testing.B) {
	setJoinBNLBufSize(2)

	for i := 0; i < b.N; i++ {
		joinBNL(books, users, UserBookMatcher, func(book Book, user User) BookWithUser {
			return BookWithUser{
				Book:     book,
				UserName: user.Name,
			}
		})
	}
}

func BenchmarkJoinBNL100X(b *testing.B) {
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

func BenchmarkJoinSortMerge(b *testing.B) {
	for i := 0; i < b.N; i++ {
		joinSortMerge(books, users, func(i, j int) bool {
			return books[i].Author < books[j].Author
		}, func(i, j int) bool {
			return users[i].Id < users[j].Id
		}, UserBookCompare, func(book Book, user User) BookWithUser {
			return BookWithUser{
				Book:     book,
				UserName: user.Name,
			}
		})
	}
}

var (
	usersX1MSM = lo.RepeatBy(1_000_000, func(index int) User {
		return users[1]
	})
	booksX1MSM = lo.RepeatBy(1_000_000, func(index int) Book {
		return books[3]
	})
)

func BenchmarkJoinSortMerge100X(b *testing.B) {
	for i := 0; i < b.N; i++ {
		joinSortMerge(booksX1MSM, usersX1MSM, func(i, j int) bool {
			return booksX1MSM[i].Author < booksX1MSM[j].Author
		}, func(i, j int) bool {
			return usersX1MSM[i].Id < usersX1MSM[j].Id
		}, UserBookCompare, func(book Book, user User) BookWithUser {
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
