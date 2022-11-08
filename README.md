# roachdiff
An offline cockroachdb schema diffing tool

## Introduction
We built this tool to help developers craft schema migrations when they make changes to columns and indexes in a `schema.sql` file. 

By default, this tool compares that schema with a version in git, and creates the neccessary `ALTER TABLE` statements that would bring the various databases inline.

## Example
The default, with no arguments, compares the schema with the version in the users git `main` branch.

```bash
$ roachdiff ./schema.sql
ALTER TABLE a ADD COLUMN foo STRING NOT NULL;
```

But you can also use a file too.
```bash
$ roachdiff --prev ./schema.old.sql ./schema.sql
ALTER TABLE a ADD COLUMN foo STRING NOT NULL;
```

## Caveats
This only produces diffs for the common cases such as adding/dropping tables, columns, indexes and constraints. If this is missing something you desire, feel free to contribute.

## Library
Available as a go library to produce your own diffs programatically
```
github.com/cygnetdigital/roachdiff/pkg/diff
```