-- +goose Up
-- +goose StatementBegin

INSERT INTO companies (name, country, creation_date) VALUES ('SM Entertainment', 'South Korea', '1995-02-14');
INSERT INTO companies (name, country, creation_date) VALUES ('JYP Entertainment', 'South Korea', '1997-04-25');
INSERT INTO companies (name, country, creation_date) VALUES ('YG Entertainment', 'South Korea', '1996-02-24');

INSERT INTO groups (name, type, debut_date, company_id) VALUES ('BTS (Bangtan Sonyeondan)', 'BG', '2013-06-13', 1);
INSERT INTO groups (name, type, debut_date, company_id) VALUES ('Twice', 'GG', '2015-10-20', 2);
INSERT INTO groups (name, type, debut_date, company_id) VALUES ('Blackpink', 'GG', '2016-08-08', 3);

INSERT INTO idols (name, stage_name, gender, idol_info_id) VALUES ('Kim Namjoon', 'RM', 'M', 1);
INSERT INTO idols (name, stage_name, gender, idol_info_id) VALUES ('Im Na-yeon', 'Nayeon', 'F', 2);
INSERT INTO idols (name, stage_name, gender, idol_info_id) VALUES ('Lalisa Manobal', 'Lisa', 'F', 3);

INSERT INTO idol_info (birthdate, height_cm, weight_kg) VALUES ('1994-09-12', 181.0, 67.0);
INSERT INTO idol_info (birthdate, height_cm, weight_kg) VALUES ('1995-09-22', 163.0, 47.0);
INSERT INTO idol_info (birthdate, height_cm, weight_kg) VALUES ('1997-03-27', 166.0, 46.0);

INSERT INTO group_members (group_id, idol_id, since_date) VALUES (1, 1, '2013-06-13');
INSERT INTO group_members (group_id, idol_id, since_date) VALUES (2, 2, '2015-10-20');
INSERT INTO group_members (group_id, idol_id, since_date) VALUES (3, 3, '2016-08-08');

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DELETE FROM group_members;
DELETE FROM idols;
DELETE FROM idol_info;
DELETE FROM groups;
DELETE FROM companies;

-- +goose StatementEnd
