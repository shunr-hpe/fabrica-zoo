# fix postgres


```
CREATE UNIQUE INDEX resources_user_username_unique
  ON resources (resource_type, (spec ->> 'username'))
  WHERE resource_type = 'User';
```
```
CREATE UNIQUE INDEX resources_user_username_unique
  ON resources (resource_type, (spec ->> 'username'))
  WHERE resource_type = 'User';
```

```
psql "$PG" -c "CREATE UNIQUE INDEX resources_username_unique ON resources ((spec ->> 'username'));"
```
```
psql "$PG" -c "CREATE UNIQUE INDEX resources_user_username_unique ON resources (resource_type, (spec ->> 'username')) WHERE resource_type = 'User';"
```
