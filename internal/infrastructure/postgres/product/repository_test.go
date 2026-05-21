package productrepo_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"github.com/your-org/your-service/internal/domain/product"
	productrepo "github.com/your-org/your-service/internal/infrastructure/postgres/product"
)

func TestRepository_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()
	pool := setupDB(t, ctx)
	repo := productrepo.NewRepository(pool)

	t.Run("Save and FindByID", func(t *testing.T) {
		p := newProduct(t)
		if err := repo.Save(ctx, p); err != nil {
			t.Fatalf("Save: %v", err)
		}

		found, err := repo.FindByID(ctx, p.ID())
		if err != nil {
			t.Fatalf("FindByID: %v", err)
		}
		if found.ID() != p.ID() {
			t.Errorf("expected id %s, got %s", p.ID(), found.ID())
		}
		if found.Name().String() != p.Name().String() {
			t.Errorf("expected name %s, got %s", p.Name(), found.Name())
		}
	})

	t.Run("FindByID returns ErrNotFound for unknown id", func(t *testing.T) {
		_, err := repo.FindByID(ctx, uuid.New())
		if !errors.Is(err, product.ErrNotFound) {
			t.Errorf("expected ErrNotFound, got %v", err)
		}
	})

	t.Run("List with pagination", func(t *testing.T) {
		pool2 := setupDB(t, ctx)
		repo2 := productrepo.NewRepository(pool2)

		for range 5 {
			p := newProduct(t)
			if err := repo2.Save(ctx, p); err != nil {
				t.Fatalf("Save: %v", err)
			}
		}

		page, err := repo2.List(ctx, product.ListFilter{Limit: 3, Offset: 0})
		if err != nil {
			t.Fatalf("List: %v", err)
		}
		if page.Total != 5 {
			t.Errorf("expected total 5, got %d", page.Total)
		}
		if len(page.Items) != 3 {
			t.Errorf("expected 3 items, got %d", len(page.Items))
		}
	})

	t.Run("Update", func(t *testing.T) {
		p := newProduct(t)
		if err := repo.Save(ctx, p); err != nil {
			t.Fatalf("Save: %v", err)
		}
		newName, _ := product.NewName("Updated")
		p.Rename(newName)

		if err := repo.Update(ctx, p); err != nil {
			t.Fatalf("Update: %v", err)
		}

		found, _ := repo.FindByID(ctx, p.ID())
		if found.Name().String() != "Updated" {
			t.Errorf("expected Updated, got %s", found.Name())
		}
	})

	t.Run("Delete", func(t *testing.T) {
		p := newProduct(t)
		if err := repo.Save(ctx, p); err != nil {
			t.Fatalf("Save: %v", err)
		}
		if err := repo.Delete(ctx, p.ID()); err != nil {
			t.Fatalf("Delete: %v", err)
		}

		_, err := repo.FindByID(ctx, p.ID())
		if !errors.Is(err, product.ErrNotFound) {
			t.Error("expected product to be deleted")
		}
	})

	t.Run("ExistsBySKU", func(t *testing.T) {
		p := newProduct(t)
		if err := repo.Save(ctx, p); err != nil {
			t.Fatalf("Save: %v", err)
		}

		exists, err := repo.ExistsBySKU(ctx, p.SKU())
		if err != nil {
			t.Fatalf("ExistsBySKU: %v", err)
		}
		if !exists {
			t.Error("expected SKU to exist")
		}

		other, _ := product.NewSKU("NO-EXIST")
		exists, _ = repo.ExistsBySKU(ctx, other)
		if exists {
			t.Error("expected SKU to not exist")
		}
	})
}

// setupDB starts a postgres container and returns a connected pool with migrations applied.
func setupDB(t *testing.T, ctx context.Context) *pgxpool.Pool {
	t.Helper()

	pgContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("test"),
		postgres.WithPassword("test"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second),
		),
	)
	if err != nil {
		t.Fatalf("start postgres container: %v", err)
	}
	t.Cleanup(func() { _ = pgContainer.Terminate(ctx) })

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		t.Fatalf("get connection string: %v", err)
	}

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		t.Fatalf("create pool: %v", err)
	}
	t.Cleanup(pool.Close)

	if err := applyMigrations(t, ctx, pool); err != nil {
		t.Fatalf("apply migrations: %v", err)
	}

	return pool
}

func applyMigrations(t *testing.T, ctx context.Context, pool *pgxpool.Pool) error {
	t.Helper()
	migrationSQL, err := os.ReadFile(filepath.Join(projectRoot(t), "migrations", "000001_create_products.up.sql"))
	if err != nil {
		return err
	}
	_, err = pool.Exec(ctx, string(migrationSQL))
	return err
}

// projectRoot walks up from the test file's directory to find the go.mod root.
func projectRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("could not determine caller path")
	}
	dir := filepath.Dir(file)
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("could not find project root (go.mod not found)")
		}
		dir = parent
	}
}

var skuCounter int

func newProduct(t *testing.T) *product.Product {
	t.Helper()
	skuCounter++
	name, _ := product.NewName("Test Product")
	desc, _ := product.NewDescription("A test product")
	price, _ := product.NewPrice(9.99)
	sku, _ := product.NewSKU(fmt.Sprintf("TST-%04d", skuCounter))
	return product.New(name, desc, price, sku)
}
