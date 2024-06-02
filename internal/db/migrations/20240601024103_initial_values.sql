-- +goose Up
-- +goose StatementBegin

INSERT INTO companies (name, country, creation_date) VALUES ('SM Entertainment', 'South Korea', '1995-02-14');
INSERT INTO companies (name, country, creation_date) VALUES ('JYP Entertainment', 'South Korea', '1997-04-25');
INSERT INTO companies (name, country, creation_date) VALUES ('YG Entertainment', 'South Korea', '1996-02-24');

INSERT INTO groups (name, type, debut_date, company_id) VALUES ('BTS (Bangtan Sonyeondan)', 'BG', '2013-06-13',
    (SELECT id FROM companies WHERE name = 'Big Hit Entertainment')
);
INSERT INTO groups (name, type, debut_date, company_id) VALUES ('Twice', 'GG', '2015-10-20',
    (SELECT id FROM companies WHERE name = 'JYP Entertainment')
);
INSERT INTO groups (name, type, debut_date, company_id) VALUES ('Blackpink', 'GG', '2016-08-08',
    (SELECT id FROM companies WHERE name = 'YG Entertainment')
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DELETE FROM groups;
DELETE FROM companies;

-- +goose StatementEnd
