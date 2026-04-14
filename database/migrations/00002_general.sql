-- +goose Up
CREATE TABLE projects (
    id TEXT PRIMARY KEY,

    name TEXT NOT NULL CHECK(name<>''),

    clone_url TEXT NOT NULL CHECK(clone_url<>''),

    installable TEXT NOT NULL,

    created INTEGER NOT NULL,
    updated INTEGER NOT NULL
);

CREATE TABLE project_tags (
    project_id TEXT NOT NULL REFERENCES proejcts(id),
    tag TEXT NOT NULL,

    PRIMARY KEY(project_id, tag)
);
