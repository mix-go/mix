/*
 Navicat Premium Data Transfer

 Source Server         : localhost_3306
 Source Server Type    : MySQL
 Source Server Version : 80027
 Source Host           : localhost:3306
 Source Schema         : test

 Target Server Type    : MySQL
 Target Server Version : 80027
 File Encoding         : 65001

 Date: 15/04/2022 22:20:33
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for xsql
-- ----------------------------
DROP TABLE IF EXISTS `xsql`;
CREATE TABLE `xsql` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `foo` varchar(255) DEFAULT NULL,
  `bar` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of xsql
-- ----------------------------
BEGIN;
INSERT INTO `xsql` (`id`, `foo`, `bar`) VALUES (1, 'v', '2022-04-14 23:49:48');
INSERT INTO `xsql` (`id`, `foo`, `bar`) VALUES (2, 'v1', '2022-04-14 23:50:00');
COMMIT;

SET FOREIGN_KEY_CHECKS = 1;
