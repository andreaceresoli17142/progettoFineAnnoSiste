DROP DATABASE IF EXISTS instanTex_db;

CREATE DATABASE instanTex_db;

USE instanTex_db;

CREATE TABLE Users (
    id INT NOT NULL AUTO_INCREMENT,
    username VARCHAR(100) NOT NULL,
    email VARCHAR(100) NOT NULL,
    date_of_join DATE NOT NULL,
    salt INT NOT NULL,
    pHash CHAR(64) NOT NULL,
	 last_login TIMESTAMP NOT NULL,
    PRIMARY KEY ( id )
);

CREATE TABLE FriendRequests (
    id INT NOT NULL AUTO_INCREMENT,
    senderId INT NOT NULL,
    reciverId INT NOT NULL,
    PRIMARY KEY ( id )
);

CREATE TABLE GroupName (
	id INT NOT NULL AUTO_INCREMENT,
	name VARCHAR(30) NOT NULL,
	description TEXT,
	date_of_creation DATE NOT NULL DEFAULT CURRENT_DATE,
	PRIMARY KEY ( id )
);

CREATE TABLE Conversations (
    id INT NOT NULL,
    participantId INT NOT NULL,
    date_of_join DATE NOT NULL DEFAULT CURRENT_DATE,
    PRIMARY KEY ( id, participantId )
);

CREATE TABLE Token (
    userid INT NOT NULL,
    accessToken CHAR(64) NOT NULL,
    act_expt INT NOT NULL,
    refreshToken CHAR(64) NOT NULL,
	rft_expt INT NOT NULL,
    PRIMARY KEY ( userid )
);

CREATE TABLE AccessToken (
	userid INT NOT NULL,
	accessToken CHAR(64) NOT NULL,
	act_expt INT NOT NULL,
	PRIMARY KEY ( userid )
);

CREATE TABLE RefreshToken (
	Id INT AUTO_INCREMENT
	userid INT NOT NULL,
	refreshToken CHAR(64) NOT NULL,
	rft_expt INT NOT NULL,
	PRIMARY KEY ( userid )
);


CREATE TABLE PwOtp (
	userId INT NOT NULL,
	otp CHAR(32) NOT NULL,
	expt INT NOT NULL,
	PRIMARY KEY (userId)
);

CREATE TABLE LoginState (
    idstring CHAR(15) NOT NULL,
    userEmail VARCHAR(45),
    PRIMARY KEY ( idstring )
);


INSERT INTO Users ( username, email, date_of_join, salt, pHash, last_login ) VALUES ( "pima", "pippo.mario@gimelli.com", CURRENT_DATE(), 66858, "ac7c27a867f92dbecc637a14afae8657f2c2a65eb47faeb3a6cadcad21c17da0", CURRENT_TIMESTAMP() );

INSERT INTO Users ( username, email, date_of_join ,salt, pHash, last_login) VALUES ( "taurone", "taurone.mario@gimelli.com", CURRENT_DATE(), 66858, "ac7c27a867f92dbecc637a14afae8657f2c2a65eb47faeb3a6cadcad21c17da0", CURRENT_TIMESTAMP());

INSERT INTO GroupName ( name, description  ) VALUES ( "gruppo tennici (senza taurone)", "tutti i miei amici odiano taurone" );

INSERT INTO Conversations ( id, participantId ) VALUES ( 1, 1 );

INSERT INTO GroupName ( name, description  ) VALUES ( "gruppo tennici", "taurone trovati un lavoro" );

INSERT INTO Conversations ( id, participantId ) VALUES ( 2, 1 );

INSERT INTO Conversations ( id, participantId ) VALUES ( 2, 2 );

CREATE TABLE MessageTable0 (
    id INT NOT NULL AUTO_INCREMENT,
    senderId INT NOT NULL,
    sendDate INT NOT NULL,
    messageText VARCHAR(1000) NOT NULL,
    attachment VARCHAR(100),
    PRIMARY KEY ( id )
);

INSERT INTO MessageTable0 ( senderId, sendDate,  messageText ) VALUES ( 0, UNIX_TIMESTAMP(), "taurone aiudo mi si è rotto lu picci" );

-- token sig h1wEJ9UE1VCGqXYk6TNoWt9WfWI8PpCNha2MyTxOWhWQCyUXWNTr8JxFX49FLUdp.qNu+Qaie3EVDNS+/Xhlg7RUjEe+/sgO/ax3UBfK39I0IemXFOjPwSaOn0NipTfCJOuvQcYv6vCjvQK6vTQbwUPKIICnMfTyDOEhAnrS31/5wzQoyIT3a29q6F5I76I0jhvEUN1qTDhR7BGp/QEUhdgfUaD3s9Jzt68PJKl3nl512fRorPqPhKikwsQIWEUfQvZqDUw4p5OvltoWeQxuTZKjfHGQe81DXv95SkYbs8gZ3aHyFD0WfQJr6tTg9HnKnMC1TLIuDgJmAKNN/L8mTjJWy8Dh5z+YwVclbKDtG5081o3fUC1OzUaEXlynk94aExAry7wV/ApwchZsZyCfBYg==
INSERT INTO Token (userid, accessToken, act_expt, refreshToken, rft_expt ) VALUES ( 1, "h1wEJ9UE1VCGqXYk6TNoWt9WfWI8PpCNha2MyTxOWhWQCyUXWNTr8JxFX49FLUdp", 1747270630, "lLTasuWtulIVIH3C98DmfgezO8hTKqonGNQicZ7HBmm5oxcXLVVpJAgW7cb98WM7", 1747270630 );
-- token sig UAfGWwjpahrXdym4pRXuS8AsItCL6jwyuVXIfU2qaKhGXyHPVCvp7bpdeNvLU2TI.1nQog3U7AVywZ9KqA7jhKT4mij8Luw2T8RrGqY6a5gxhn7g5TdHdTzVZjxs9KL2awfr++CN10j/0oLMSrEZ8WK/yMZaeeR8ImyXxTiryTinyphzpdmGm2gwCBTdsPUmAzXBt4goyHOKrMPxR2ceqfhJ+/h7du3ZPJMt6iasBzDYDLHXhJys/z4TcWHCcWhxXQkgqjskAkxbpSjlgYYrTszdcsPeI1eCKnD4W7X+YMGFTJqXxamyLFuR5Em2rbMAk1+yRVPHeOLPdAnd1TcVjC7z9CFiFACeMKkKxR0P9iHo6Wcv4QfuVEo6+CDlNX+HSs4Z2WcSmjkGZmDC8wkhWog==
INSERT INTO Token (userid, accessToken, act_expt, refreshToken, rft_expt ) VALUES ( 2, "UAfGWwjpahrXdym4pRXuS8AsItCL6jwyuVXIfU2qaKhGXyHPVCvp7bpdeNvLU2TI", 1747270630, "meH7pBndzCrQJqI7W0Ii2VwFXLV8j8C6EyrApBrhIga4jqe8O8DHPtDmUZ75BKt7", 1747270630 );