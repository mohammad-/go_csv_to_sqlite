1. Create a folder for data

2. Sync Data
aws s3 sync s3://bodybank-enterprise-shoplist-prd-request-data-bucket .

3. Create DB
sqlite3 user_data.db

4. Create Table

CREATE TABLE `userinfo` (
`createdAtCompound` INTEGER,
`userId` TEXT,
`height` INTEGER,
`status` TEXT,
`createdAt` INTEGER,
`id` TEXT,
`gender` TEXT,
`weight` INTEGER,
`age` INTEGER,
`updatedAt` INTEGER,
`errorCode` TEXT,
`created_idx` TEXT,
`errorDetail` TEXT,
`time_dif` TEXT,
 PRIMARY KEY (id, userId)
);

CREATE INDEX created_data_idx ON userinfo (`createdAt`);
CREATE INDEX userid_created_data_idx ON userinfo (`createdAt`, `userId`);

5. change mode to csv
.mode csv

6. List all the files to import (If possible in different terminal)
ls *.csv | xargs -I% echo '.import ./% userinfo' | pbcopy


7. Paste list copied in step 6 in sqlite3

8. delete header row
delete from userinfo where id=='id';

9. Run following sql to find request and user information

select id, userId,  datetime(`createdAt`, 'unixepoch', 'localtime')
from userinfo
where userId not in
    (
        select distinct userId from userinfo where
        datetime(`createdAt`, 'unixepoch', 'localtime') < '2019-07-31 23:59:59'
    )
and status=='completed' and
datetime(`createdAt`, 'unixepoch', 'localtime') < '2019-09-01 00:00:00';

select count(id)
from userinfo
where userId not in
    (
        select distinct userId from userinfo where
        datetime(`createdAt`, 'unixepoch', 'localtime') < '2019-07-31 23:59:59'
    )
and status=='completed' and
datetime(`createdAt`, 'unixepoch', 'localtime') < '2019-09-01 00:00:00';


select userId, id, datetime(`createdAt`, 'unixepoch', 'localtime')
from userinfo
where userId not in
    (
        select distinct userId from userinfo where
        datetime(`createdAt`, 'unixepoch', 'localtime') < '2019-07-31 23:59:59'
    )
and status=='completed' and
datetime(`createdAt`, 'unixepoch', 'localtime') < '2019-09-01 00:00:00';

select count(distinct userId)
from userinfo
where userId not in
    (
        select distinct userId from userinfo where
        datetime(createdAt, 'unixepoch', 'localtime') < '2019-08-01 00:00:00'
    )
and status=='completed'
and datetime(createdAt, 'unixepoch', 'localtime') < '2019-08-30 23:59:59';

select count(distinct userId)
from userinfo
where userId not in
    (
        select distinct userId from userinfo where
        datetime(`createdAt`, 'unixepoch', 'localtime') < '2019-08-14 23:59:59'
    )
and status=='completed'
and datetime(`createdAt`, 'unixepoch', 'localtime') < '2019-08-16 00:00:00';

select count(distinct userId) from userinfo where datetime(`createdAt`, 'unixepoch', 'localtime') > '2019-08-02 00:00:00' and datetime(`createdAt`, 'unixepoch', 'localtime') < '2019-08-02 23:59:59';