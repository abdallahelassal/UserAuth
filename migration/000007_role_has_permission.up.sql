CREATE TABLE role_has_permission(
    role_id    UUID     NOT NULL,
    permission_id UUID NOT NULL,

    CONSTRAINT pk_role_has_permission PRIMARY KEY (role_id,permission_id),

    CONSTRAINT fk_role_has_permission_roles FOREIGN KEY (role_id)
        REFERENCES roles(id)
        ON DELETE CASCADE 
        ON UPDATE CASCADE,
    CONSTRAINT fk_role_has_permission_permissions FOREIGN KEY (permission_id)
        REFERENCES permissions(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE    
);
