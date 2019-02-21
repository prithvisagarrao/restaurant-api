#!/bin/sh

psql -p 5432 -d recipes_dev -U recipes_usr < rp.sql
