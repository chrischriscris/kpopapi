#!/usr/bin/env bash

GOOSE_DRIVER=sqlite3 GOOSE_DBSTRING=../../../db.sqlite3 goose up

