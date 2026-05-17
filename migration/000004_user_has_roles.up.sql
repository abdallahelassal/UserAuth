CREATE TABLE user_has_roles(
    user_id     UUID    NOT NULL,
    role_id     UUID    NOT NULL,

    CONSTRAINT pk_user_has_roles    PRIMARY KEY (user_id, role_id),
    CONSTRAINT fk_user_has_roles_user FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    CONSTRAINT fk_user_has_roles_role FOREIGN KEY (role_id)
        REFERENCES roles(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE    

);