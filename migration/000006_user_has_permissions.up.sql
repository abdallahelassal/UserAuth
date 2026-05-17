CREATE TABLE user_has_permissions(
    user_id UUID NOT NULL,
    permission_id UUID NOT NULL,

    CONSTRAINT pk_user_has_permissions PRIMARY KEY (user_id, permission_id),
    
    CONSTRAINT fk_user_has_permissions_users FOREIGN KEY (user_id)
           REFERENCES users(id)
           ON DELETE CASCADE 
           ON UPDATE CASCADE,
    CONSTRAINT fk_user_has_permissions_permission FOREIGN KEY (permission_id)
        REFERENCES permissions(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE       
)