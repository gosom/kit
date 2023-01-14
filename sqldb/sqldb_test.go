package sqldb_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gosom/kit/sqldb"

	"github.com/stretchr/testify/require"
)

type testObject struct {
	ID   int
	Name string
}

func (t *testObject) Bind() []any {
	return []any{&t.ID, &t.Name}
}

func TestOpen(t *testing.T) {
	t.Run("TestThatNewDBReturnsDB", func(t *testing.T) {
		db := sqldb.NewDB("dummy", "dsn")
		require.NotNil(t, db)
		require.IsType(t, &sqldb.DB{}, db)
		require.Equal(t, "dummy", db.DriverName)
		require.Equal(t, "dsn", db.DSN)
		require.NotNil(t, db.Now)
	})
	t.Run("TestDbConnMethod", func(t *testing.T) {
		conn, mock, err := sqlmock.New()
		require.NoError(t, err)
		mock.ExpectClose()
		db := sqldb.NewDB("dummy", "dsn")
		db.SetPool(conn)
		require.NotNil(t, db.Conn())
		err = db.Close()
		require.NoError(t, err)
	})
	t.Run("TestDbBeginTxMethod", func(t *testing.T) {
		conn, mock, err := sqlmock.New()
		require.NoError(t, err)
		mock.ExpectBegin()
		db := sqldb.NewDB("dummy", "dsn")
		db.SetPool(conn)
		tx, err := db.BeginTx(context.Background(), nil)
		require.NoError(t, err)
		require.NotNil(t, tx)
		require.IsType(t, &sql.Tx{}, tx)
		require.Implements(t, (*sqldb.DBTX)(nil), tx)
	})
	t.Run("TestQueryFunc", func(t *testing.T) {
		conn, mock, err := sqlmock.New()
		require.NoError(t, err)
		mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "test"))
		db := sqldb.NewDB("dummy", "dsn")
		db.SetPool(conn)
		items, err := sqldb.Query[testObject](context.Background(), db.Conn(), "SELECT * from tests")
		require.NoError(t, err)
		require.Len(t, items, 1)
		require.Equal(t, 1, items[0].ID)
		require.Equal(t, "test", items[0].Name)
	})
	t.Run("TestQueryFuncErr", func(t *testing.T) {
		conn, mock, err := sqlmock.New()
		require.NoError(t, err)
		mock.ExpectQuery("SELECT").WillReturnError(sql.ErrConnDone)
		db := sqldb.NewDB("dummy", "dsn")
		db.SetPool(conn)
		items, err := sqldb.Query[testObject](context.Background(), db.Conn(), "SELECT * from tests")
		require.Error(t, err)
		require.Len(t, items, 0)
	})
	t.Run("TestQueryRowFunc", func(t *testing.T) {
		conn, mock, err := sqlmock.New()
		require.NoError(t, err)
		mock.ExpectQuery("SELECT").WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "test"))
		db := sqldb.NewDB("dummy", "dsn")
		db.SetPool(conn)
		item, err := sqldb.QueryRow[testObject](context.Background(), db.Conn(), "SELECT * from tests")
		require.NoError(t, err)
		require.Equal(t, 1, item.ID)
		require.Equal(t, "test", item.Name)
	})
	t.Run("TestQueryRowFuncNoRows", func(t *testing.T) {
		conn, mock, err := sqlmock.New()
		require.NoError(t, err)
		mock.ExpectQuery("SELECT").WillReturnError(sql.ErrNoRows)
		db := sqldb.NewDB("dummy", "dsn")
		db.SetPool(conn)
		_, err = sqldb.QueryRow[testObject](context.Background(), db.Conn(), "SELECT * from tests")
		require.Error(t, err)
		require.Equal(t, sql.ErrNoRows, err)
	})
}
