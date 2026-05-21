package productrepo

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/your-org/your-service/internal/domain/product"
	productdb "github.com/your-org/your-service/internal/infrastructure/postgres/product/sqlc"
)

// Repository is the PostgreSQL implementation of product.Repository.
type Repository struct {
	q *productdb.Queries
}

// NewRepository constructs a Repository backed by the given connection pool.
func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{q: productdb.New(pool)}
}

// Save persists a new product row.
func (r *Repository) Save(ctx context.Context, p *product.Product) error {
	_, err := r.q.CreateProduct(ctx, productdb.CreateProductParams{
		ID:          p.ID(),
		Name:        p.Name().String(),
		Description: p.Description().String(),
		Price:       p.Price().Value(),
		SKU:         p.SKU().String(),
	})
	if err != nil {
		return fmt.Errorf("create product row: %w", err)
	}
	return nil
}

// FindByID returns the product with the given ID or product.ErrNotFound.
func (r *Repository) FindByID(ctx context.Context, id uuid.UUID) (*product.Product, error) {
	row, err := r.q.GetProductByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, product.ErrNotFound
		}
		return nil, fmt.Errorf("get product row: %w", err)
	}
	return toDomain(row)
}

// List returns a paginated page of products matching the given filter.
func (r *Repository) List(ctx context.Context, f product.ListFilter) (product.Page, error) {
	total, err := r.q.CountProducts(ctx)
	if err != nil {
		return product.Page{}, fmt.Errorf("count products: %w", err)
	}

	rows, err := r.q.ListProducts(ctx, productdb.ListProductsParams{
		Limit:  f.Limit,
		Offset: f.Offset,
	})
	if err != nil {
		return product.Page{}, fmt.Errorf("list products: %w", err)
	}

	items := make([]*product.Product, 0, len(rows))
	for _, row := range rows {
		p, err := toDomain(row)
		if err != nil {
			return product.Page{}, err
		}
		items = append(items, p)
	}
	return product.Page{Items: items, Total: total}, nil
}

// Update persists changes to an existing product row.
func (r *Repository) Update(ctx context.Context, p *product.Product) error {
	_, err := r.q.UpdateProduct(ctx, productdb.UpdateProductParams{
		ID:          p.ID(),
		Name:        p.Name().String(),
		Description: p.Description().String(),
		Price:       p.Price().Value(),
		SKU:         p.SKU().String(),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return product.ErrNotFound
		}
		return fmt.Errorf("update product row: %w", err)
	}
	return nil
}

// Delete removes the product row with the given ID.
func (r *Repository) Delete(ctx context.Context, id uuid.UUID) error {
	if err := r.q.DeleteProduct(ctx, id); err != nil {
		return fmt.Errorf("delete product row: %w", err)
	}
	return nil
}

// ExistsBySKU reports whether a product with the given SKU already exists.
func (r *Repository) ExistsBySKU(ctx context.Context, sku product.SKU) (bool, error) {
	exists, err := r.q.ExistsProductBySKU(ctx, sku.String())
	if err != nil {
		return false, fmt.Errorf("exists by sku: %w", err)
	}
	return exists, nil
}

// toDomain maps a sqlc row to the domain aggregate.
// This is the ONLY place that knows about both layers.
func toDomain(row productdb.Product) (*product.Product, error) {
	name, err := product.NewName(row.Name)
	if err != nil {
		return nil, fmt.Errorf("reconstitute name: %w", err)
	}
	desc, err := product.NewDescription(row.Description)
	if err != nil {
		return nil, fmt.Errorf("reconstitute description: %w", err)
	}
	price, err := product.NewPrice(row.Price)
	if err != nil {
		return nil, fmt.Errorf("reconstitute price: %w", err)
	}
	sku, err := product.NewSKU(row.SKU)
	if err != nil {
		return nil, fmt.Errorf("reconstitute sku: %w", err)
	}
	return product.Reconstitute(row.ID, name, desc, price, sku, row.CreatedAt, row.UpdatedAt), nil
}
