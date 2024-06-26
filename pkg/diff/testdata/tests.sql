---- Test: creates single table
----
CREATE TABLE a (id string PRIMARY KEY);
----
CREATE TABLE a (id STRING PRIMARY KEY);
----

---- Test: creates multiple tables
----
CREATE TABLE a (id string PRIMARY KEY);
CREATE TABLE b (id string PRIMARY KEY);
----
CREATE TABLE a (id STRING PRIMARY KEY);
CREATE TABLE b (id STRING PRIMARY KEY);
----

---- Test: creates one of two tables
CREATE TABLE a (id string PRIMARY KEY);
----
CREATE TABLE a (id string PRIMARY KEY);
CREATE TABLE b (id string PRIMARY KEY);
----
CREATE TABLE b (id STRING PRIMARY KEY);
----

---- Test: create table noop
CREATE TABLE a (id string PRIMARY KEY);
----
CREATE TABLE a (id string PRIMARY KEY);
----
----

---- Test: drop table
CREATE TABLE a (id string PRIMARY KEY);
----
----
DROP TABLE a;
----



---- Test: adds a single column
CREATE TABLE a (id string PRIMARY KEY);
----
CREATE TABLE a (id string PRIMARY KEY, foo string NOT NULL);
----
ALTER TABLE a ADD COLUMN foo STRING NOT NULL;
----

---- Test: drops single column
CREATE TABLE a (id string PRIMARY KEY, foo string NOT NULL);
----
CREATE TABLE a (id string PRIMARY KEY);
----
ALTER TABLE a DROP COLUMN foo;
----

---- Test: adds multiple columns
CREATE TABLE a (id string PRIMARY KEY);
----
CREATE TABLE a (id string PRIMARY KEY, foo string NOT NULL, bar string NOT NULL);
----
ALTER TABLE a ADD COLUMN foo STRING NOT NULL;
ALTER TABLE a ADD COLUMN bar STRING NOT NULL;
----

---- Test: adds multiple columns across tables
CREATE TABLE a (id string PRIMARY KEY);
CREATE TABLE b (id string PRIMARY KEY);
----
CREATE TABLE a (id string PRIMARY KEY, foo string NOT NULL);
CREATE TABLE b (id string PRIMARY KEY, bar string NOT NULL);
----
ALTER TABLE a ADD COLUMN foo STRING NOT NULL;
ALTER TABLE b ADD COLUMN bar STRING NOT NULL;
----

---- Test: adds one column, drops another
CREATE TABLE a (id string PRIMARY KEY, foo string NOT NULL);
----
CREATE TABLE a (id string PRIMARY KEY, bar string NOT NULL);
----
ALTER TABLE a DROP COLUMN foo;
ALTER TABLE a ADD COLUMN bar STRING NOT NULL;
----



---- Test: creates a single index
CREATE TABLE a (id string PRIMARY KEY);
----
CREATE TABLE a (id string PRIMARY KEY);
CREATE INDEX idx ON a (id DESC);
----
CREATE INDEX idx ON a (id DESC);
----

---- Test: drops a single index
CREATE TABLE a ( id string PRIMARY KEY, foo string NOT NULL);
CREATE INDEX foo_idx ON a (foo DESC);
----
CREATE TABLE a ( id string PRIMARY KEY, foo string NOT NULL);
----
DROP INDEX a@foo_idx;
----

---- Test: adds multiple indexes
CREATE TABLE a (id string PRIMARY KEY);
----
CREATE TABLE a (id string PRIMARY KEY);
CREATE INDEX foo_idx ON a (id);
CREATE INDEX bar_idx ON a (id);
----
CREATE INDEX foo_idx ON a (id);
CREATE INDEX bar_idx ON a (id);
----

---- Test: Add unique and partial index
CREATE TABLE a (id string PRIMARY KEY);
----
CREATE TABLE a (id string PRIMARY KEY);
CREATE UNIQUE INDEX foo_idx ON a (id) WHERE id = 'baz';
----
CREATE UNIQUE INDEX foo_idx ON a (id) WHERE id = 'baz';
----

---- Test: adds multiple indexes across tables
CREATE TABLE a (id string PRIMARY KEY);
CREATE TABLE b (id string PRIMARY KEY);
----
CREATE TABLE a (id string PRIMARY KEY);
CREATE INDEX foo_idx ON a (id);
CREATE TABLE b (id string PRIMARY KEY);
CREATE INDEX bar_idx ON b (id);
----
CREATE INDEX foo_idx ON a (id);
CREATE INDEX bar_idx ON b (id);
----

---- Test: adds one index, drops another
CREATE TABLE a (id string PRIMARY KEY);
CREATE INDEX foo_idx ON a (id);
----
CREATE TABLE a (id string PRIMARY KEY, INDEX bar_idx (id));
----
DROP INDEX a@foo_idx;
CREATE INDEX bar_idx ON a (id);
----

---- Test: add table and index from blank slate
----
CREATE TABLE a (id string PRIMARY KEY, foo string NOT NULL);
CREATE INDEX foo_idx ON a (id);
----
CREATE TABLE a (id STRING PRIMARY KEY, foo STRING NOT NULL);
CREATE INDEX foo_idx ON a (id);
----

---- Test: drop table and index
CREATE TABLE a (id string PRIMARY KEY, foo string NOT NULL);
CREATE INDEX foo_idx ON a (id);
----
----
DROP TABLE a;
----

---- Test: add constraints
CREATE TABLE a (id string PRIMARY KEY);
CREATE TABLE b (id string PRIMARY KEY);
CREATE TABLE c (id string PRIMARY KEY);
----
CREATE TABLE a (id string PRIMARY KEY, CONSTRAINT bar_ctn UNIQUE (id));
CREATE TABLE b (id string PRIMARY KEY, CONSTRAINT bar_ctn FOREIGN KEY (id) REFERENCES bb);
CREATE TABLE c (id string PRIMARY KEY, CONSTRAINT bar_ctn CHECK (id IN ('yes', 'no', 'unknown')));
----
ALTER TABLE a ADD CONSTRAINT bar_ctn UNIQUE (id);
ALTER TABLE b ADD CONSTRAINT bar_ctn FOREIGN KEY (id) REFERENCES bb;
ALTER TABLE c ADD CONSTRAINT bar_ctn CHECK (id IN ('yes', 'no', 'unknown'));
----


---- Test: create table with inline constraint
----
CREATE TABLE a (id string PRIMARY KEY, CONSTRAINT bar_ctn UNIQUE (id));
----
CREATE TABLE a (id STRING PRIMARY KEY, CONSTRAINT bar_ctn UNIQUE (id));
----


---- Test: drop a constraint
CREATE TABLE a (id string PRIMARY KEY, CONSTRAINT bar_ctn UNIQUE (id));
CREATE TABLE b (id string PRIMARY KEY, CONSTRAINT bar_ctn FOREIGN KEY (id) REFERENCES bb);
CREATE TABLE c (id string PRIMARY KEY, CONSTRAINT bar_ctn CHECK (id IN ('yes', 'no', 'unknown')));
----
CREATE TABLE a (id string PRIMARY KEY);
CREATE TABLE b (id string PRIMARY KEY);
CREATE TABLE c (id string PRIMARY KEY);
----
ALTER TABLE a DROP CONSTRAINT bar_ctn;
ALTER TABLE b DROP CONSTRAINT bar_ctn;
ALTER TABLE c DROP CONSTRAINT bar_ctn;
----


---- Test: drop a foreign key reference
CREATE TABLE a (id string PRIMARY KEY);
CREATE TABLE b (id string PRIMARY KEY, a_id string NOT NULL REFERENCES a (id));
----
CREATE TABLE a (id string PRIMARY KEY);
CREATE TABLE b (id string PRIMARY KEY, a_id string NOT NULL);
----
ALTER TABLE b DROP CONSTRAINT b_a_id_fkey;
----


---- Test: modify a foreign key constraint
CREATE TABLE a (id string PRIMARY KEY);
CREATE TABLE b (id string PRIMARY KEY, a_id string NOT NULL REFERENCES a (id));
----
CREATE TABLE a (id string PRIMARY KEY);
CREATE TABLE b (id string PRIMARY KEY, a_id string NOT NULL REFERENCES a (id) ON DELETE CASCADE);
----
ALTER TABLE b DROP CONSTRAINT b_a_id_fkey;
ALTER TABLE b ADD CONSTRAINT b_a_id_fkey FOREIGN KEY (a_id) REFERENCES a (id) ON DELETE CASCADE;
----


---- Test: columns can be set NOT NULL
CREATE TABLE a (id string PRIMARY KEY, foo string NULL);
----
CREATE TABLE a (id string PRIMARY KEY, foo string NOT NULL);
----
ALTER TABLE a ALTER COLUMN foo SET NOT NULL;
----

---- Test: columns can drop NOT NULL
CREATE TABLE a (id string PRIMARY KEY, foo string NOT NULL);
----
CREATE TABLE a (id string PRIMARY KEY, foo string NULL);
----
ALTER TABLE a ALTER COLUMN foo DROP NOT NULL;
----

---- Test: columns can be set NOT NULL from silent null
CREATE TABLE a (id string PRIMARY KEY, foo string);
----
CREATE TABLE a (id string PRIMARY KEY, foo string NOT NULL);
----
ALTER TABLE a ALTER COLUMN foo SET NOT NULL;
----

---- Test: columns can drop NOT NULL to silent null
CREATE TABLE a (id string PRIMARY KEY, foo string NOT NULL);
----
CREATE TABLE a (id string PRIMARY KEY, foo string);
----
ALTER TABLE a ALTER COLUMN foo DROP NOT NULL;
----

---- Test: enums can be created
----
CREATE TYPE a AS ENUM ('foo', 'bar', 'baz');
----
CREATE TYPE a AS ENUM ('foo', 'bar', 'baz');
----

---- Test: enums can be dropped
CREATE TYPE a AS ENUM ('foo', 'bar', 'baz');
----
----
DROP TYPE a;
----

---- Test: enums can be modified to be larger
CREATE TYPE a AS ENUM ('foo', 'bar');
----
CREATE TYPE a AS ENUM ('foo', 'baz', 'bar');
----
ALTER TYPE a ADD VALUE 'baz' AFTER 'foo';
----

---- Test: enums can be enlarged at the end
CREATE TYPE a AS ENUM ('foo', 'bar');
----
CREATE TYPE a AS ENUM ('foo', 'bar', 'baz');
----
ALTER TYPE a ADD VALUE 'baz' AFTER 'bar';
----

---- Test: enums can be enlarged at the beginning
CREATE TYPE a AS ENUM ('foo', 'bar');
----
CREATE TYPE a AS ENUM ('baz', 'foo', 'bar');
----
ALTER TYPE a ADD VALUE 'baz' BEFORE 'foo';
----

---- Test: enums can be enlarged from empty
CREATE TYPE a AS ENUM ();
----
CREATE TYPE a AS ENUM ('foo');
----
ALTER TYPE a ADD VALUE 'foo';
----


---- Test: enums can be modified to be smaller
CREATE TYPE a AS ENUM ('foo', 'bar', 'baz');
----
CREATE TYPE a AS ENUM ('foo', 'bar');
----
ALTER TYPE a DROP VALUE 'baz';
----

---- Test: enums can be modified to be empty
CREATE TYPE a AS ENUM ('foo', 'bar');
----
CREATE TYPE a AS ENUM ();
----
ALTER TYPE a DROP VALUE 'foo';
ALTER TYPE a DROP VALUE 'bar';
----

---- Test: complicated enum modifications
CREATE TYPE a AS ENUM ('foo', 'bar', 'baz');
----
CREATE TYPE a AS ENUM ('corge', 'grault', 'foo','waldo', 'fred', 'plugh', 'bar', 'baz', 'qux', 'qix');
----
ALTER TYPE a ADD VALUE 'corge' BEFORE 'foo';
ALTER TYPE a ADD VALUE 'grault' BEFORE 'foo';
ALTER TYPE a ADD VALUE 'waldo' AFTER 'foo';
ALTER TYPE a ADD VALUE 'fred' AFTER 'waldo';
ALTER TYPE a ADD VALUE 'plugh' AFTER 'fred';
ALTER TYPE a ADD VALUE 'qux' AFTER 'baz';
ALTER TYPE a ADD VALUE 'qix' AFTER 'qux';
----

---- Test: enums are defined before the table they are used in.
CREATE TABLE b (id STRING PRIMARY KEY);
----
CREATE TYPE a AS ENUM ('foo', 'bar', 'baz');
CREATE TABLE b (id STRING PRIMARY KEY, a_id a);
CREATE TABLE c (id STRING PRIMARY KEY, a_id a);
----
CREATE TYPE a AS ENUM ('foo', 'bar', 'baz');
CREATE TABLE c (id STRING PRIMARY KEY, a_id a);
ALTER TABLE b ADD COLUMN a_id a;
----
