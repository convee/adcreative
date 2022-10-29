SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for advertiser_audit
-- ----------------------------
DROP TABLE IF EXISTS `advertiser_audit`;
CREATE TABLE `advertiser_audit` (
`id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'id',
`advertiser_name` varchar(255) NOT NULL DEFAULT '0' COMMENT '广告主名称',
`publisher_id` int(11) NOT NULL DEFAULT '0' COMMENT '媒体ID',
`customer_id` int(11) NOT NULL DEFAULT '0' COMMENT '客户ID',
`advertiser_id` int(11) NOT NULL DEFAULT '0' COMMENT '广告主id',
`status` tinyint(3) NOT NULL DEFAULT '0' COMMENT '审核状态：0待审核，1审核通过，2审核不通过',
`info` longtext,
`is_rsync` tinyint(3) DEFAULT '0' COMMENT '是否同步',
`created_at` timestamp NULL DEFAULT NULL COMMENT '创建时间',
`updated_at` timestamp NULL DEFAULT NULL COMMENT '更新时间',
`deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
`publisher_account_id` int(11) NOT NULL DEFAULT '0',
`err_code` tinyint(3) DEFAULT NULL COMMENT '内部错误码1上传失败、2上传不通过、3待送审、4送审失败、5审核中、6查询失败、7审核通过、8审核不通过',
`err_msg` text COMMENT '内部错误提示',
`media_cid` varchar(255) DEFAULT NULL COMMENT '媒体方广告主ID',
`extra` text COMMENT 'sohu token等信息',
PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of advertiser_audit
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for advertiser_rules
-- ----------------------------
DROP TABLE IF EXISTS `advertiser_rules`;
CREATE TABLE `advertiser_rules` (
`id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'id',
`publisher_id` int(11) NOT NULL DEFAULT '0' COMMENT '媒体ID',
`info` json NOT NULL COMMENT '广告主送审规则信息',
`created_at` timestamp NULL DEFAULT NULL COMMENT '创建时间',
`updated_at` timestamp NULL DEFAULT NULL COMMENT '更新时间',
`deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of advertiser_rules
-- ----------------------------
BEGIN;
INSERT INTO `advertiser_rules` (`id`, `publisher_id`, `info`, `created_at`, `updated_at`, `deleted_at`) VALUES (1, 34, '[{\"key\": \"advertiser_name\", \"desc\": \"广告主名称\", \"size\": 0, \"type\": \"txt\", \"limit\": 0, \"format\": [\"text\"], \"measure\": [], \"disabled\": 1, \"required\": 1}, {\"key\": \"authorize_state\", \"desc\": \"客户URL\", \"size\": 0, \"type\": \"txt\", \"limit\": 0, \"format\": [\"text\"], \"measure\": [], \"required\": 1}, {\"key\": \"qualifications\", \"desc\": \"广告主资质文件\", \"size\": 102400, \"type\": \"file\", \"limit\": 1, \"format\": [\"zip\", \"jpg\", \"jpeg\", \"png\", \"bmp\"], \"measure\": [], \"required\": 1}, {\"key\": \"industry\", \"desc\": \"行业\", \"size\": 0, \"type\": \"select\", \"limit\": 0, \"format\": [\"text\"], \"measure\": [], \"required\": 1}]', '2022-03-01 11:25:21', '2022-03-01 11:25:25', NULL);
COMMIT;

-- ----------------------------
-- Table structure for creative
-- ----------------------------
DROP TABLE IF EXISTS `creative`;
CREATE TABLE `creative` (
`id` int(11) unsigned NOT NULL AUTO_INCREMENT,
`position_id` int(11) NOT NULL DEFAULT '0' COMMENT '广告位ID',
`publisher_id` int(11) NOT NULL DEFAULT '0' COMMENT '媒体ID',
`creative_id` varchar(100) NOT NULL DEFAULT '' COMMENT '送审方创意ID',
`media_cid` varchar(2000) DEFAULT NULL COMMENT '媒体方创意ID',
`industry` varchar(255) DEFAULT NULL COMMENT '媒体行业',
`name` varchar(200) NOT NULL DEFAULT '' COMMENT '送审方创意名称',
`template_id` varchar(100) NOT NULL DEFAULT '' COMMENT '广告位模板ID',
`info` longtext COMMENT '素材信息',
`land_url` longtext COMMENT '落地页',
`deeplink_url` longtext COMMENT 'app唤起地址',
`start_date` varchar(20) NOT NULL DEFAULT '' COMMENT '开始时间',
`end_date` varchar(20) NOT NULL DEFAULT '' COMMENT '结束时间',
`monitor` json DEFAULT NULL COMMENT '曝光监测',
`cm` json DEFAULT NULL COMMENT '点击监测',
`vm` json DEFAULT NULL COMMENT '可见性监测',
`action` int(11) NOT NULL DEFAULT '0' COMMENT '广告交互类型：1-打开网页 2-下载 3-deeplink',
`material_id` varchar(255) NOT NULL DEFAULT '' COMMENT '素材ID',
`mini_program_id` varchar(100) NOT NULL DEFAULT '' COMMENT '小程序ID',
`mini_program_path` text NOT NULL COMMENT '小程序路径',
`customer_id` int(11) NOT NULL DEFAULT '0' COMMENT '客户ID',
`advertiser_id` int(11) NOT NULL DEFAULT '0' COMMENT '广告主ID',
`status` int(11) NOT NULL DEFAULT '0' COMMENT '素材状态：0待审核，1审核通过，2审核不通过',
`media_info` text NOT NULL COMMENT '媒体送审所需信息',
`reason` text COMMENT '审核失败原因',
`created_at` timestamp NULL DEFAULT NULL COMMENT '创建时间',
`updated_at` timestamp NULL DEFAULT NULL COMMENT '更新时间',
`deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
`is_rsync` tinyint(3) NOT NULL DEFAULT '0' COMMENT '是否同步',
`publisher_account_id` int(11) NOT NULL DEFAULT '0' COMMENT '媒体账号ID',
`err_code` tinyint(3) DEFAULT NULL COMMENT '内部错误码1上传失败、2上传不通过、3待送审、4送审失败、5审核中、6查询失败、7审核通过、8审核不通过',
`err_msg` text COMMENT '内部错误提示',
`request_id` varchar(128) NOT NULL DEFAULT '' COMMENT '最近一次送审requestId',
`extra` text COMMENT '额外信息',
`video_cdn_url` text COMMENT '媒体视频url',
`pic_cdn_url` text COMMENT '媒体图片url',
`pub_return_url` text COMMENT '媒体返回素材cdn地址',
`priority` tinyint(3) NOT NULL DEFAULT '1' COMMENT '送审优先级权重（1-10，10最高）',
PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of creative
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for customer
-- ----------------------------
DROP TABLE IF EXISTS `customer`;
CREATE TABLE `customer` (
`id` int(11) unsigned NOT NULL AUTO_INCREMENT,
`appid` varchar(32) NOT NULL DEFAULT '' COMMENT '应用标识',
`secret` varchar(128) NOT NULL DEFAULT '' COMMENT '密钥',
`name` varchar(100) NOT NULL DEFAULT '' COMMENT '客户名称',
`is_private` tinyint(3) NOT NULL DEFAULT '0' COMMENT '是否私有账号',
`updated_at` timestamp NULL DEFAULT NULL COMMENT '更新时间',
`created_at` timestamp NULL DEFAULT NULL COMMENT '创建时间',
`deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除',
`creative_callback_url` varchar(255) DEFAULT NULL COMMENT '创意回调地址',
`advertiser_callback_url` varchar(255) DEFAULT NULL COMMENT '广告主回调地址',
PRIMARY KEY (`id`) USING BTREE,
UNIQUE KEY `idx_appid` (`appid`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of customer
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for position
-- ----------------------------
DROP TABLE IF EXISTS `position`;
CREATE TABLE `position` (
`id` int(11) NOT NULL AUTO_INCREMENT COMMENT '广告位Id',
`publisher_id` int(11) NOT NULL DEFAULT '0' COMMENT '媒体Id',
`name` varchar(50) NOT NULL DEFAULT '' COMMENT '广告位名称',
`position` varchar(64) NOT NULL DEFAULT '' COMMENT '广告位标识',
`material_info` longtext NOT NULL COMMENT '素材详细信息',
`is_support_deeplink` tinyint(3) NOT NULL DEFAULT '1' COMMENT '是否支持deeplink',
`pv_limit` tinyint(1) NOT NULL DEFAULT '0' COMMENT '曝光监测条数上限，默认为0，则为不限',
`cl_limit` tinyint(1) NOT NULL DEFAULT '0' COMMENT '点击监测条数上限，默认为0，则为不限',
`created_at` timestamp NULL DEFAULT NULL COMMENT '创建时间',
`updated_at` timestamp NULL DEFAULT NULL COMMENT '更新时间',
`deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
PRIMARY KEY (`id`) USING BTREE,
UNIQUE KEY `idx_position` (`position`) USING BTREE,
KEY `idx_publisher_id` (`publisher_id`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of position
-- ----------------------------
BEGIN;
INSERT INTO `position` (`id`, `publisher_id`, `name`, `position`, `material_info`, `is_support_deeplink`, `pv_limit`, `cl_limit`, `created_at`, `updated_at`, `deleted_at`) VALUES (1, 1, '手机QQ空间_Mobile_信息流_第七条横板视频', 'Tencent_Mobile_QQKongjianAPP_If_SShipin', '{\"list\": [{\"id\": \"1\", \"info\": [{\"key\": \"abstract\", \"name\": \"广告文案\", \"format\": [\"txt\"], \"measure\": [\"1\", \"30\"], \"required\": 1}, {\"key\": \"logo_name\", \"name\": \"商标名称\", \"format\": [\"txt\"], \"measure\": [\"1\", \"20\"], \"required\": 1}, {\"key\": \"cover\", \"name\": \"封面图\", \"size\": 100, \"format\": [\"jpg\", \"png\", \"jpeg\"], \"measure\": [\"1280*720\"], \"required\": 1}, {\"key\": \"icon\", \"name\": \"商标图片\", \"size\": 400, \"format\": [\"jpg\", \"jpeg\", \"png\"], \"measure\": [\"512*512\"], \"required\": 1}, {\"key\": \"video\", \"name\": \"视频\", \"size\": 51200, \"format\": [\"mp4\"], \"measure\": [\"1280*720\"], \"required\": 1}], \"name\": \"手机QQ空间_Mobile_信息流_第七条横板视频\", \"display_id\": 113003}]}', 0, 4, 3, NULL, NULL, NULL);
COMMIT;

-- ----------------------------
-- Table structure for publisher
-- ----------------------------
DROP TABLE IF EXISTS `publisher`;
CREATE TABLE `publisher` (
`id` int(11) NOT NULL AUTO_INCREMENT COMMENT '媒体Id',
`name` varchar(50) NOT NULL COMMENT '媒体名称',
`is_rsync_advertiser` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否广告主送审',
`is_rsync_creative` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否创意送审',
`created_at` timestamp NULL DEFAULT NULL COMMENT '创建时间',
`updated_at` timestamp NULL DEFAULT NULL COMMENT '更新时间',
`deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
`advertiser_urls` json DEFAULT NULL COMMENT '广告主送审url',
`creative_urls` json DEFAULT NULL COMMENT '创意送审url',
`pv_limit` tinyint(1) NOT NULL DEFAULT '0' COMMENT '曝光监测条数上限，默认为0，则为不限',
`cl_limit` tinyint(1) NOT NULL DEFAULT '0' COMMENT '点击监测条数上限，默认为0，则为不限',
`nickname` varchar(255) NOT NULL DEFAULT '' COMMENT '昵称',
PRIMARY KEY (`id`) USING BTREE,
UNIQUE KEY `idx_name` (`name`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of publisher
-- ----------------------------
BEGIN;
INSERT INTO `publisher` (`id`, `name`, `is_rsync_advertiser`, `is_rsync_creative`, `created_at`, `updated_at`, `deleted_at`, `advertiser_urls`, `creative_urls`, `pv_limit`, `cl_limit`, `nickname`) VALUES (1, 'Tencent', 0, 1, '2022-10-29 21:25:32', '2022-10-29 21:25:32', NULL, '{}', '{\"query_url\": \"https://open.adx.qq.com/adcreative/list\", \"create_url\": \"https://open.adx.qq.com/adcreative/add\", \"extend_url\": \"https://open.adx.qq.com/adcreative/extend\", \"update_url\": \"https://open.adx.qq.com/adcreative/update\"}', 4, 3, '腾讯');
COMMIT;

-- ----------------------------
-- Table structure for publisher_account
-- ----------------------------
DROP TABLE IF EXISTS `publisher_account`;
CREATE TABLE `publisher_account` (
`id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'account_id',
`publisher_id` int(11) NOT NULL DEFAULT '0' COMMENT '媒体账户的ID',
`dsp_id` varchar(256) NOT NULL DEFAULT '' COMMENT 'DSP的ID',
`token` varchar(256) NOT NULL DEFAULT '' COMMENT 'DSP的令牌',
`customer_id` int(11) NOT NULL DEFAULT '0' COMMENT '客户ID',
`remark` text NOT NULL COMMENT '备注',
`extra` json DEFAULT NULL COMMENT '额外信息',
`callback_url` text COMMENT '送审回调',
`created_at` timestamp NULL DEFAULT NULL COMMENT '创建时间',
`updated_at` timestamp NULL DEFAULT NULL COMMENT '更新时间',
`deleted_at` timestamp NULL DEFAULT NULL COMMENT '删除时间',
PRIMARY KEY (`id`) USING BTREE,
KEY `publisher_id` (`publisher_id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of publisher_account
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for publisher_industry
-- ----------------------------
DROP TABLE IF EXISTS `publisher_industry`;
CREATE TABLE `publisher_industry` (
`id` bigint(20) unsigned NOT NULL AUTO_INCREMENT,
`name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL DEFAULT '' COMMENT '行业名称',
`pid` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '父级ID',
`level` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '类目级别',
`sort` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '序号',
`created_at` timestamp NULL DEFAULT NULL,
`updated_at` timestamp NULL DEFAULT NULL,
`deleted_at` timestamp NULL DEFAULT NULL,
`type_id` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '行业id',
`publisher` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '媒体',
PRIMARY KEY (`id`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Records of publisher_industry
-- ----------------------------
BEGIN;
COMMIT;

SET FOREIGN_KEY_CHECKS = 1;
