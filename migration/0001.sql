-- Dumping structure for table ethereum.histories
CREATE TABLE IF NOT EXISTS `histories` (
  `address` varchar(50) NOT NULL,
  `datetime` datetime DEFAULT NULL,
  PRIMARY KEY (`address`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;
