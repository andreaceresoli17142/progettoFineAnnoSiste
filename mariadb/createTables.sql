DROP DATABASE IF EXISTS instanTex_db;

CREATE DATABASE instanTex_db;

USE instanTex_db;

CREATE TABLE Users (
    id INT NOT NULL AUTO_INCREMENT,
    username VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL,
    date_of_join DATE NOT NULL,
    salt INT NOT NULL,
    pHash VARCHAR(64) NOT NULL,
	 last_login TIMESTAMP NOT NULL,
    PRIMARY KEY ( id )
);

CREATE TABLE FriendRequests (
    id INT NOT NULL AUTO_INCREMENT,
    senderId INT NOT NULL,
    reciverId INT NOT NULL,
    PRIMARY KEY ( id )
);

CREATE TABLE Conversations (
    id INT NOT NULL,
    participantId INT NOT NULL,
    date_of_join DATE NOT NULL,
    PRIMARY KEY ( id, participantId )
);

CREATE TABLE Token (
    userid INT NOT NULL,
    accessToken VARCHAR(64) NOT NULL,
    expireTime INT NOT NULL,
    refreshToken VARCHAR(64) NOT NULL,
    PRIMARY KEY ( userid )
);


CREATE TABLE LoginState (
    idstring VARCHAR(15) NOT NULL,
    PRIMARY KEY ( idstring )
);


INSERT INTO Users ( username, email, date_of_join, salt, pHash, last_login ) VALUES ( "pima", "pippo.mario@gimelli.com", CURRENT_DATE(), 123456, "62d18522b74d75b2a84776c91ba5498377441d4c4af0cea22ca7de9e09475d3a", CURRENT_TIMESTAMP() );

INSERT INTO Users ( username, email, date_of_join ,salt, pHash, last_login) VALUES ( "taurone", "taurone.mario@gimelli.com", CURRENT_DATE(), 123456, "62d18522b74d75b2a84776c91ba5498377441d4c4af0cea22ca7de9e09475d3a", CURRENT_TIMESTAMP());


INSERT INTO Conversations VALUES ( 0, 0, CURRENT_DATE() );

INSERT INTO Conversations VALUES ( 0, 1, CURRENT_DATE() );

CREATE TABLE MessageTable0 (
    id INT NOT NULL AUTO_INCREMENT,
    senderId INT NOT NULL,
    sendDate INT NOT NULL,
    messageText VARCHAR(1000) NOT NULL,
    attachment VARCHAR(100),
    PRIMARY KEY ( id )
);

INSERT INTO MessageTable0 ( senderId, sendDate,  messageText ) VALUES ( 0, UNIX_TIMESTAMP(), "taurone aiudo mi si è rotto lu picci" );
