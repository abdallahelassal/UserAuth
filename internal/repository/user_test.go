package repository

import (
	"context"
	"fmt"
	"testing"

	"github.com/abdallahelassal/UserAuth/domain"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)


func setupTestpostgresDB(t *testing.T)(*gorm.DB, func() ){
	t.Helper()
	dsn := "host=localhost user=postgres password=postgres port=5432 sslmode=disable dbname=test_db"
	db , err := gorm.Open(postgres.Open(dsn),&gorm.Config{})
	if err != nil {
		t.Fatalf("failed connect db: %v",err)
	}
	tx := db.Begin()
	cleanup := func ()  {
		tx.Rollback()
	}
	return tx ,cleanup
}

func setupRepo(t *testing.T)(*UserRepository,func()){
	db , cleanup := setupTestpostgresDB(t)
	repo  := NewUserRepository(db)
	return repo,cleanup
}

func TestCreate(t *testing.T){
	t.Helper()
	repo,cleanup := setupRepo(t)
	defer cleanup()
	
	user := &domain.User{
		Base: domain.Base{
			ID: uuid.New(),
		},
		UserName: "testuser",
		Email: "testuser@example.com",
	}
	err := repo.Create(context.Background(),user)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	stored, err := repo.FindByID(context.Background(),user.ID)
	if err != nil {
		t.Fatalf("failed to featch user : %v",err)
	}
	if  stored.Email != user.Email{
		t.Errorf("expected email %s got %s", stored.Email,user.Email)
	} 
}

func TestGetByEmail(t *testing.T){
	t.Helper()
	repo , cleanup := setupRepo(t)
	defer cleanup()
	

	user := &domain.User{
		Base: domain.Base{ID: uuid.New()},
		UserName: "abdallah",
		Email: "testuser@test.com",
	}
	 err := repo.Create(context.Background(), user)
	if err != nil {
		t.Fatalf("faild create user : %s", err)
	}
	get , err := repo.GetByEmail(context.Background(),user.Email)
	if err != nil {
		t.Fatalf("failed fetch data : %s", err)
	}
	if get.Email != user.Email{
		t.Errorf("expected email %s got %s", user.Email,get.Email)
	}
}
 func TestGetByName(t *testing.T){
	repo , cleanup := setupRepo(t)
	defer cleanup()

	user := &domain.User{
		Base: domain.Base{ID: uuid.New()},
		UserName: "ahmed",
		Email: "ahmed@test.com",
	}
	if err := repo.Create(context.Background(),user);err != nil {
		t.Fatalf("failded create user : %s",err)
	}
	
	get , err := repo.GetByName(context.Background(),user.UserName)
	if err != nil {
		t.Fatalf("failed get name : %s",err)
	}
	if get.UserName != user.UserName{

		t.Errorf("expected user name %s got %s", user.UserName , get.UserName)
	}
 }

 func TestFindByID(t *testing.T) {
	t.Helper()
	repo , cleanup := setupRepo(t)
	defer cleanup()
	user := &domain.User{
		Base: domain.Base{ID: uuid.New()},
		UserName: "yousef",
		Email: "yousef@test.com",
	}
	if err := repo.Create(context.Background(),user); err != nil {
		t.Fatalf("failed craete user : %s",err)
	}
	get , err := repo.FindByID(context.Background(),user.ID)
	if err != nil {
		t.Fatalf("failed find userID please check : %s",err)
	}
	if get.ID != user.ID {
		t.Errorf("expected user ID %s got %s",user.ID , get.ID)
	}
 }

func TestFetch(t *testing.T){
	repo , cleanup := setupRepo(t)
	defer cleanup()
	ctx := context.Background()
	users := []domain.User{
		{
			Base: domain.Base{ID: uuid.New()},
			UserName: "user1",
			Email:    "u1@test.com",
		},
		{
			Base: domain.Base{ID: uuid.New()},
			UserName: "user2",
			Email:    "u2@test.com",
		},
		{
			Base: domain.Base{ID: uuid.New()},
			UserName: "user3",
			Email:    "u3@test.com",
		},
	}
	for _ , v := range users{
		if err := repo.Create(ctx, &v);err != nil {
			t.Fatalf("faild create users :%s",err)
		}
	}

	result ,nextCursor , err := repo.Fetch(ctx,"",2)
	if err != nil {
		t.Fatalf("fetch error %s",err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 users got %d", len(result))
	}
	if nextCursor == ""{
		t.Errorf("expected next cursor but got empty")
	}
}

func TestFetch_SecoundPage(t *testing.T){
	repo , cleanup := setupRepo(t)
	defer cleanup()
	ctx := context.Background()
	for i := 0; i <3; i++ {
		u := domain.User{
			Base: domain.Base{ID: uuid.New()},
			UserName: fmt.Sprintf("user%d",i),
			Email: fmt.Sprintf("u%d@test.com",i),
		}
		if err := repo.Create(ctx,&u); err != nil {
			t.Fatalf("faild created user %d", err)
		}
	}

	firstPage , cursor , err := repo.Fetch(ctx,"",2)
	if err != nil {
		t.Fatalf("failed fetch first page %d",err)
	}
	if len(firstPage) != 2{
		t.Errorf("expected 2 users got %d",err)
	}
	secoundPage , _,err := repo.Fetch(ctx,cursor,2)
	if err != nil {
		t.Fatalf("failed fetch secound page %d", err)
	}
	if len(secoundPage) == 0 {
		t.Errorf("expected secound page data %d",err)
	} 
}

func TestAssignRole(t *testing.T) {
	db , cleanup := setupTestpostgresDB(t)
	defer cleanup()
	userRepo := NewUserRepository(db)	
	roleRepo := NewRoleRepository(db)

	ctx := context.Background()
	user := &domain.User{
		Base: domain.Base{ID: uuid.New()},
		UserName: "ahmed",
		Email: "atest@test.com",
	}
	role := &domain.Role{
		Base: domain.Base{ID: uuid.New()},
		Name: "admin",
	}
	if err := roleRepo.Create(ctx,role); err != nil {
		t.Fatalf("failed create role %d",err)
	}
	if err := userRepo.Create(ctx,user);err != nil {
		t.Fatalf("failed create user %d", err )
	}
	assignRole := userRepo.AssignRole(ctx,user.ID,role.ID)
	if assignRole != nil {
		t.Errorf("failed fetch assignRole %d",assignRole)
	}

}