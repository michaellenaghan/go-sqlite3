package tests

import (
	"bytes"
	"crypto/rand"
	"errors"
	"io"
	"testing"

	"github.com/ncruces/go-sqlite3"
	_ "github.com/ncruces/go-sqlite3/embed"
)

func TestBlob(t *testing.T) {
	t.Parallel()

	db, err := sqlite3.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	err = db.Exec(`CREATE TABLE IF NOT EXISTS test (col)`)
	if err != nil {
		t.Fatal(err)
	}

	err = db.Exec(`INSERT INTO test VALUES (zeroblob(1024))`)
	if err != nil {
		t.Fatal(err)
	}

	blob, err := db.OpenBlob("main", "test", "col", db.LastInsertRowID(), true)
	if err != nil {
		t.Fatal(err)
	}
	defer blob.Close()

	size := blob.Size()
	if size != 1024 {
		t.Errorf("got %d, want 1024", size)
	}

	var data [1280]byte
	_, err = rand.Read(data[:])
	if err != nil {
		t.Fatal(err)
	}

	_, err = io.Copy(blob, bytes.NewReader(data[:size/2]))
	if err != nil {
		t.Fatal(err)
	}

	_, err = io.Copy(blob, bytes.NewReader(data[:]))
	if !errors.Is(err, sqlite3.ERROR) {
		t.Fatal("want error")
	}

	_, err = io.Copy(blob, bytes.NewReader(data[size/2:size]))
	if err != nil {
		t.Fatal(err)
	}

	_, err = blob.Seek(size/4, io.SeekStart)
	if err != nil {
		t.Fatal(err)
	}

	if got, err := io.ReadAll(blob); err != nil {
		t.Fatal(err)
	} else if !bytes.Equal(got, data[size/4:size]) {
		t.Errorf("got %q, want %q", got, data[size/4:size])
	}

	if n, err := blob.Read(make([]byte, 1)); n != 0 || err != io.EOF {
		t.Errorf("got (%d, %v), want (0, EOF)", n, err)
	}

	if err := blob.Close(); err != nil {
		t.Fatal(err)
	}

	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
}

func TestBlob_invalid(t *testing.T) {
	t.Parallel()

	db, err := sqlite3.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	err = db.Exec(`CREATE TABLE IF NOT EXISTS test (col)`)
	if err != nil {
		t.Fatal(err)
	}

	err = db.Exec(`INSERT INTO test VALUES (zeroblob(1024))`)
	if err != nil {
		t.Fatal(err)
	}

	_, err = db.OpenBlob("", "test", "col", db.LastInsertRowID(), false)
	if !errors.Is(err, sqlite3.ERROR) {
		t.Fatal("want error")
	}
}

func TestBlob_readonly(t *testing.T) {
	t.Parallel()

	db, err := sqlite3.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	err = db.Exec(`CREATE TABLE IF NOT EXISTS test (col)`)
	if err != nil {
		t.Fatal(err)
	}

	err = db.Exec(`INSERT INTO test VALUES (zeroblob(1024))`)
	if err != nil {
		t.Fatal(err)
	}

	blob, err := db.OpenBlob("main", "test", "col", db.LastInsertRowID(), false)
	if err != nil {
		t.Fatal(err)
	}
	defer blob.Close()

	_, err = blob.Write([]byte("data"))
	if !errors.Is(err, sqlite3.READONLY) {
		t.Fatal("want error")
	}
}

func TestBlob_expired(t *testing.T) {
	t.Parallel()

	db, err := sqlite3.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	err = db.Exec(`CREATE TABLE IF NOT EXISTS test (col)`)
	if err != nil {
		t.Fatal(err)
	}

	err = db.Exec(`INSERT INTO test VALUES (zeroblob(1024))`)
	if err != nil {
		t.Fatal(err)
	}

	blob, err := db.OpenBlob("main", "test", "col", db.LastInsertRowID(), false)
	if err != nil {
		t.Fatal(err)
	}
	defer blob.Close()

	err = db.Exec(`DELETE FROM test`)
	if err != nil {
		t.Fatal(err)
	}

	_, err = io.ReadAll(blob)
	if !errors.Is(err, sqlite3.ABORT) {
		t.Fatal("want error", err)
	}
}

func TestBlob_Seek(t *testing.T) {
	t.Parallel()

	db, err := sqlite3.Open(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	err = db.Exec(`CREATE TABLE IF NOT EXISTS test (col)`)
	if err != nil {
		t.Fatal(err)
	}

	err = db.Exec(`INSERT INTO test VALUES (zeroblob(1024))`)
	if err != nil {
		t.Fatal(err)
	}

	blob, err := db.OpenBlob("main", "test", "col", db.LastInsertRowID(), true)
	if err != nil {
		t.Fatal(err)
	}
	defer blob.Close()

	_, err = blob.Seek(0, 10)
	if err == nil {
		t.Fatal("want error")
	}

	_, err = blob.Seek(-1, io.SeekCurrent)
	if err == nil {
		t.Fatal("want error")
	}

	n, err := blob.Seek(1, io.SeekEnd)
	if err != nil {
		t.Fatal(err)
	}
	if n != blob.Size()+1 {
		t.Errorf("got %d, want %d", n, blob.Size())
	}

	_, err = blob.Write([]byte("data"))
	if !errors.Is(err, sqlite3.ERROR) {
		t.Fatal("want error")
	}
}
