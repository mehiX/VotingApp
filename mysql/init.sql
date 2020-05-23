
CREATE TABLE IF NOT EXISTS `votes` (
  `voterID` varchar(64) NOT NULL,
  `vote` VARCHAR(20),
  PRIMARY KEY `pk_id`(`voterID`)
) ENGINE = InnoDB;
