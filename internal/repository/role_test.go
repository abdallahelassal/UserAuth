package repository

import (
	"context"
	"errors"
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

func TestFindByName(t *testing.T){
	repo , cleanup := setupRoleRepo(t)
	defer cleanup()
	role := &domain.Role{
		Base: domain.Base{ID: uuid.New()},
		Name: "admin",
	}
	if err := repo.Create(context.Background(),role); err != nil {
		t.Fatalf("failed create role %v",err)
	}

 	a, err := repo.FindByName(context.Background(), role.Name)
  	if err != nil {
		t.Fatalf("failed fetch role %v", err)
	}
	if a.Name != role.Name{
		t.Errorf("expected role %v got %v", role.Name , a.Name)
	}
}

func TestFindAll(t *testing.T){
	repo , cleanup := setupRoleRepo(t)
	defer cleanup()
	ctx := context.Background()
	roles := []domain.Role{
		{
			Base: domain.Base{ID: uuid.New()},
			Name: "admin",
		},
				{
			Base: domain.Base{ID: uuid.New()},
			Name: "user",
		},
	}
	for i:= range roles{
		if err := repo.Create(ctx,&roles[i]); err != nil {
			t.Fatalf("failed create role %v", err)
		}
	}
	got , err := repo.FindAll(ctx)
	if err != nil {
		t.Fatalf("failed fetch roles %v", err)
	}
	if len(got) != len(roles){
		t.Fatalf("expected role %v got %v", len(roles),len(got))
	}
	
}

func TestGetRoleByUserID(t *testing.T){
	db , cleanup := setupTestPostgresDB(t)
	defer cleanup()
	userRepo := NewUserRepository(db)
	roleRepo := NewRoleRepository(db)
	ctx := context.Background()
	user := &domain.User{
		Base: domain.Base{ID: uuid.New()},
		UserName: "abdallah",
		Email: "abdallah@test.com",
	}
	role := &domain.Role{
		Base: domain.Base{ID: uuid.New()},
		Name: "admin",
	}
	if err := userRepo.Create(ctx,user) ;err != nil {
		t.Fatalf("failed create user %v", err)
	}
	if err := roleRepo.Create(ctx , role); err != nil {
		t.Fatalf("failed create role %v", err)
	}
	if err := userRepo.AssignRole(ctx, user.ID, role.ID); err != nil {
		t.Fatalf("failed assign role to user %v", err)
	}
	got , err := roleRepo.GetRolesByUserID(ctx, user.ID)
	if err != nil {
		t.Fatalf("failed fetch role by user %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected role %v got %v", 1, len(got))
	}
	if got[0].Name != role.Name{
		t.Errorf("expected role %v got %v", role.Name, got[0].Name)
	}
}
	
func TestUpdateRole(t *testing.T){
	repo , cleanup := setupRoleRepo(t)
	defer cleanup()
	ctx := context.Background()
	role := &domain.Role{
		Base: domain.Base{ID: uuid.New()},
		Name: "admin",
	}
	if err := repo.Create(ctx,role); err != nil {
		t.Fatalf("failed create role %v", err)
	}
	role.Name = "user"
	if err := repo.Update(ctx,role); err != nil {
		t.Fatalf("failed update role %v", err)
	}
	stored , err := repo.FindByID(ctx,role.ID)
	if err != nil {
		t.Fatalf("failed fetch role %v", err)
	}
	if stored.Name != role.Name{
		t.Errorf("expected role %v got %v", stored.Name,role.Name)
	}
}
func TestDeleteRole(t *testing.T){
	repo , cleanup := setupRoleRepo(t)
	defer cleanup()
	ctx := context.Background()
	role := &domain.Role{
		Base: domain.Base{ID: uuid.New()},
		Name: "admin",
	}
	if err := repo.Create(ctx,role); err != nil {
		t.Fatalf("failed create role %v", err)
	}
	if err := repo.Delete(ctx,role.ID); err != nil {
		t.Fatalf("failed delete role %v", err)
	}
role, err := repo.FindByID(ctx, role.ID)

if !errors.Is(err, domain.ErrNotFound) {
	t.Fatalf("expected ErrRoleNotFound, got %v", err)
}

if role != nil {
	t.Fatalf("expected nil role after delete")
}
}

func TestRemovePermissions(t *testing.T){
	db , cleanup := setupTestPostgresDB(t)
	defer cleanup()
	roleRepo := NewRoleRepository(db)
	permissionRepo := NewPermissionRepository(db)
	ctx := context.Background()

	role := &domain.Role{
		Base: domain.Base{ID: uuid.New()},
		Name: "admin",
	}
	permission := []domain.Permission{
		{
			Base: domain.Base{ID: uuid.New()},
			Name: "read",
		},
		{
			Base: domain.Base{ID: uuid.New()},
			Name: "write",
		},
	}
	if err := roleRepo.Create(ctx,role); err != nil {
		t.Fatalf("failed create role %v", err)
	}
	for i := range permission{
		if err := permissionRepo.Create(ctx,&permission[i]); err != nil {
			t.Fatalf("failed create permission %v", err)
		}
	var permissionIDs []uuid.UUID
	for _, p := range permission{
		permissionIDs = append(permissionIDs, p.ID)
	}
	if err := roleRepo.AssignPermission(ctx, role.ID, permissionIDs); err != nil {
		t.Fatalf("failed assign permissions %v", err)
	}
	if err := roleRepo.RemoveAllPermission(ctx, role.ID); err != nil {
		t.Fatalf("failed remove permissions %v", err)
	}
		
	}
}
