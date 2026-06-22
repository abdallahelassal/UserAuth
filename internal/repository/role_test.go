package repository

import (
	"context"
	"testing"

	"github.com/abdallahelassal/UserAuth/domain"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


func setupTestPostgresDB(t *testing.T)(*gorm.DB,func()){
	dsn := "host=localhost dbname=test_db port=5432 user=postgres sslmode=disable password=postgres"
	db , err := gorm.Open(postgres.Open(dsn),&gorm.Config{})
	if err != nil {
		t.Fatalf("failed connect database %v",err)
	}

	tx := db.Begin()
	
	cleanup := func ()  {
		 tx.Rollback()
	}
	return tx,cleanup
}

func setupRoleRepo(t *testing.T)(*RoleRepository, func()){
	db , cleanup := setupTestPostgresDB(t)
	repo := NewRoleRepository(db)
	return repo , cleanup
}

func TestCreate_Role(t *testing.T) {
	repo , cleanup := setupRoleRepo(t)
	defer cleanup()
	role := &domain.Role{
		Base: domain.Base{ID: uuid.New()},
		Name: "admin",
	}
	if err := repo.Create(context.Background(),role); err != nil {
		t.Fatalf("failed create role %v",err)
	}
	stored , err := repo.FindByID(context.Background(),role.ID)
	if err != nil {
		t.Fatalf("failed fetch role %v",err)
	}
	if role.Name != stored.Name{
		t.Errorf("expected %v got %v", role.Name,stored.Name)
	}
}

func TestAssignPermissions(t *testing.T){
	db , cleanup := setupTestPostgresDB(t)
	defer cleanup()
	roleRepo := NewRoleRepository(db)
	permissionRepo := NewPermissionRepository(db)
	ctx := context.Background()

	role := &domain.Role{
		Base: domain.Base{ID: uuid.New()},
		Name: "admin",
	}
	permissions := []domain.Permission{
		{
			Base: domain.Base{ID: uuid.New()},
			Name: "admin_create",
		},
		{
			Base: domain.Base{ID: uuid.New()},
			Name: "admin_delete",
		},
	}
	for i:= range permissions{
		 err := permissionRepo.Create(ctx,&permissions[i])
		if err != nil {
			t.Fatalf("failed create permissions %v",err)
		}
	}
	if err := roleRepo.Create(ctx,role); err != nil {
		t.Fatalf("failed create role %v",err)
	}
	var permissionIDs []uuid.UUID
	for _ , p := range permissions{
		permissionIDs = append(permissionIDs, p.ID)
	}
	assignPermissions := roleRepo.AssignPermission(ctx, role.ID, permissionIDs)
	if assignPermissions != nil {
		t.Errorf("failed assign permission %v",assignPermissions)
	}
}

func TestByID(t *testing.T){
	repo, cleanup := setupRoleRepo(t)
	defer cleanup()
	ctx := context.Background()
	role := &domain.Role{
		Base: domain.Base{ID: uuid.New()},
		Name: "admin",
	}
	if err := repo.Create(ctx,role); err != nil {
		t.Fatalf("failed create role %v", err)
	}
	stored , err := repo.FindByID(ctx,role.ID)
	if err != nil {
		t.Fatalf("failed fetch role %v", err)
	}
	if stored.Name != role.Name{
		t.Errorf("expected role %v got %v", stored.Name,role.Name)
	}
}


	


