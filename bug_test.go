package bug

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"testing"

	"entgo.io/bug/ent/user"
	"entgo.io/ent/dialect"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"

	"entgo.io/bug/ent"
	"entgo.io/bug/ent/enttest"
)

func TestBugSQLite(t *testing.T) {
	client := enttest.Open(t, dialect.SQLite, "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()
	test(t, client)
}

func TestBugMySQL(t *testing.T) {
	for version, port := range map[string]int{"56": 3306, "57": 3307, "8": 3308} {
		addr := net.JoinHostPort("localhost", strconv.Itoa(port))
		t.Run(version, func(t *testing.T) {
			client := enttest.Open(t, dialect.MySQL, fmt.Sprintf("root:pass@tcp(%s)/test?parseTime=True", addr))
			defer client.Close()
			test(t, client)
		})
	}
}

func TestBugPostgres(t *testing.T) {
	for version, port := range map[string]int{"10": 5430, "11": 5431, "12": 5432, "13": 5433, "14": 5434} {
		t.Run(version, func(t *testing.T) {
			client := enttest.Open(t, dialect.Postgres, fmt.Sprintf("host=localhost port=%d user=postgres dbname=test password=pass sslmode=disable", port))
			defer client.Close()
			test(t, client)
		})
	}
}

func TestBugMaria(t *testing.T) {
	for version, port := range map[string]int{"10.5": 4306, "10.2": 4307, "10.3": 4308} {
		t.Run(version, func(t *testing.T) {
			addr := net.JoinHostPort("localhost", strconv.Itoa(port))
			client := enttest.Open(t, dialect.MySQL, fmt.Sprintf("root:pass@tcp(%s)/test?parseTime=True", addr))
			defer client.Close()
			test(t, client)
		})
	}
}

func test(t *testing.T, client *ent.Client) {
	ctx := context.Background()
	client.User.Delete().ExecX(ctx)

	t.Run("custom type solution", func(t *testing.T) {
		r := client.User.Create().
			SetName("Ariel").
			SetContent(`<a onblur="alert(secret)" href="http://www.google.com">Google</a>`).
			SetAge(30).
			SaveX(ctx)

		// content returned from Save/SaveX
		fmt.Println(r.Content) // <a onblur="alert(secret)" href="http://www.google.com">Google</a>

		// content save
		r = client.User.Query().Where(user.Name("Ariel")).OnlyX(ctx)
		fmt.Println(r.Content) // <a href="http://www.google.com" rel="nofollow">Google</a>
	})

	t.Run("hook solution", func(t *testing.T) {
		p := client.Post.Create().SetTitle(`<a onblur="alert(secret)" href="http://www.google.com">Google</a>`).SaveX(ctx)

		fmt.Println(p.Title) // <a href="http://www.google.com" rel="nofollow">Google</a>

		p = client.Post.Query().FirstX(ctx)
		fmt.Println(p.Title) // <a href="http://www.google.com" rel="nofollow">Google</a>

		p = client.Post.UpdateOne(p).SetTitle(`<a onblur="alert(secret)" href="http://www.google.com">Google</a>`).SaveX(ctx)
		fmt.Println(p.Title) // <a href="http://www.google.com" rel="nofollow">Google</a>

		p = client.Post.Query().FirstX(ctx)
		fmt.Println(p.Title) // <a href="http://www.google.com" rel="nofollow">Google</a>

		client.Post.Update().SetTitle(`<a onblur="alert(secret)" href="http://www.google.com">Google</a>`).ExecX(ctx)

		p = client.Post.Query().FirstX(ctx)
		fmt.Println(p.Title) // <a href="http://www.google.com" rel="nofollow">Google</a>
	})
}
