CREATE TABLE IF NOT EXISTS request
(
    id       SERIAL NOT NULL PRIMARY KEY,
    method   text   NOT NULL,
    scheme   text   NOT NULL,
    host     text   NOT NULL,
    path     text   NOT NULL,
    cookies  text   NOT NULL,
    header   text   default '',
    body     text   default ''
);

CREATE TABLE IF NOT EXISTS response
(
    id           SERIAL NOT NULL PRIMARY KEY,
    request_id   SERIAL REFERENCES request (id) NOT NULL,
    code         INT                            NOT NULL,
    message      text                           NOT NULL,
    cookies      text                           NOT NULL,
    header       text                           default '',
    body         text                           default ''
);