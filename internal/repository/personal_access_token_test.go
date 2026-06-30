package repository

import (
	"context"
	"testing"
	

	"github.com/abdallahelassal/UserAuth/domain"
	"github.com/google/uuid"
)



func setupPersonRepo(t *testing.T)(*PersonalAccessTokenRepository,func()){
	db , cleanup := setupTestPostgresDB(t)
	repo := NewPersonalAccessTokenRepository(db)
	return repo, cleanup
}

func TestCreatePersonalAccessToken(t *testing.T){
	db  , cleanup := setupTestPostgresDB(t)
	defer cleanup()
	persRepo := NewPersonalAccessTokenRepository(db)
	userRepo := NewUserRepository(db)

	ctx := context.Background()
	user := &domain.User{
		Base: domain.Base{ID: uuid.New()},
		UserName: "abdallah",
		Email: "abdallah@test.com",
	}
	if err  := userRepo.Create(ctx,user); err != nil{
		t.Fatalf("failed create user %v", err)
	}
	token := &domain.PersonalAccessToken{
		ID: uuid.New(),
		TokenHash: "testhash",
		UserID: user.ID,
		
	}


	if err := persRepo.Create(ctx,token); err != nil {
		t.Fatalf("failed create personal access token %v",err)
	}
	stored, err := persRepo.FindByToken(ctx, token.TokenHash)
	if err != nil {
		t.Fatalf("failed fetch personal access token %v", err)
	}
	if stored.TokenHash != token.TokenHash{
		t.Errorf("expected personal access token %v got %v", stored.TokenHash,token.TokenHash)
	}
}

func TestFindByToken(t *testing.T) {
	db  , cleanup := setupTestPostgresDB(t)
	defer cleanup()
	persRepo := NewPersonalAccessTokenRepository(db)
	userRepo := NewUserRepository(db)

	ctx := context.Background()
	user := &domain.User{
		Base: domain.Base{ID: uuid.New()},
		UserName: "abdallah",
		Email: "abdallah@test.com",
	}
	if err  := userRepo.Create(ctx,user); err != nil{
		t.Fatalf("failed create user %v", err)
	}
	token := &domain.PersonalAccessToken{
		ID: uuid.New(),
		TokenHash: "testhash",
		UserID: user.ID,
		
	}


	if err := persRepo.Create(ctx,token); err != nil {
		t.Fatalf("failed create personal access token %v",err)
	}
	stored , err := persRepo.FindByToken(ctx,token.TokenHash)
	if err != nil {
		t.Fatalf("faied find By token %v", err)
	}
	if token.TokenHash != stored.TokenHash{
		t.Errorf("expected token %v got %v", stored.TokenHash, stored.TokenHash)
	}	
}

func TestUpdateLastUseAt(t *testing.T) {
	db  , cleanup := setupTestPostgresDB(t)
	defer cleanup()
	persRepo := NewPersonalAccessTokenRepository(db)
	userRepo := NewUserRepository(db)

	ctx := context.Background()
	user := &domain.User{
		Base: domain.Base{ID: uuid.New()},
		UserName: "abdallah",
		Email: "abdallah@test.com",
	}
	if err  := userRepo.Create(ctx,user); err != nil{
		t.Fatalf("failed create user %v", err)
	}
	token := &domain.PersonalAccessToken{
		ID: uuid.New(),
		TokenHash: "testhash",
		UserID: user.ID,
		
	}
	if err := persRepo.Create(ctx,token); err != nil {
		t.Fatalf("failed create personal access token %v",err)
	}

	 err := persRepo.UpdateLastUseAt(ctx, token.TokenHash)
	if err != nil {
		t.Fatalf("failed update %v", err)
	}


	stored , err := persRepo.FindByToken(ctx,token.TokenHash)
	if err != nil {
		t.Fatalf("faied find By token %v", err)
	}
	if token.TokenHash != stored.TokenHash{
		t.Errorf("expected token %v got %v", stored.TokenHash, stored.TokenHash)
	}



	if stored.LastUsedAt.IsZero(){
		t.Errorf("expected last use at to be updated but got zero value")
	}

}

func TestDelete_personalToken(t *testing.T) {
		db  , cleanup := setupTestPostgresDB(t)
	defer cleanup()
	persRepo := NewPersonalAccessTokenRepository(db)
	userRepo := NewUserRepository(db)

	ctx := context.Background()
	user := &domain.User{
		Base: domain.Base{ID: uuid.New()},
		UserName: "abdallah",
		Email: "abdallah@test.com",
	}
	if err  := userRepo.Create(ctx,user); err != nil{
		t.Fatalf("failed create user %v", err)
	}
	token := &domain.PersonalAccessToken{
		ID: uuid.New(),
		TokenHash: "testhash",
		UserID: user.ID,
		
	}
	if err := persRepo.Create(ctx,token); err != nil {
		t.Fatalf("failed create personal access token %v",err)
	}
	if err := persRepo.Delete(ctx, token.ID); err != nil {
		t.Fatalf("failed delete personal token %v ", err)
	}
}

func TestDeleteByUserID(t *testing.T) {
	db  , cleanup := setupTestPostgresDB(t)
	defer cleanup()
	persRepo := NewPersonalAccessTokenRepository(db)
	userRepo := NewUserRepository(db)

	ctx := context.Background()
	user := &domain.User{
		Base: domain.Base{ID: uuid.New()},
		UserName: "abdallah",
		Email: "abdallah@test.com",
	}
	if err  := userRepo.Create(ctx,user); err != nil{
		t.Fatalf("failed create user %v", err)
	}
	token := &domain.PersonalAccessToken{
		ID: uuid.New(),
		TokenHash: "testhash",
		UserID: user.ID,
		
	}
	if err := persRepo.Create(ctx,token); err != nil {
		t.Fatalf("failed create personal access token %v",err)
	}
	if err := persRepo.DeleteByUserID(ctx, user.ID); err != nil {
		t.Fatalf("failed delete personal token by user id %v ", err)
	}
}	