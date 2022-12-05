-- MySQL dump 10.13  Distrib 8.0.22, for Linux (x86_64)
--
-- Host: localhost    Database: migrate_db
-- ------------------------------------------------------
-- Server version	8.0.22

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `api_keys`
--

DROP TABLE IF EXISTS `api_keys`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `api_keys` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` int NOT NULL,
  `access_key` varbinary(255) NOT NULL,
  `secret_key` varbinary(255) NOT NULL,
  `last_auth` timestamp NULL DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `access_key` (`access_key`),
  KEY `user_id` (`user_id`),
  CONSTRAINT `api_keys_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `auth_scheme_data`
--

DROP TABLE IF EXISTS `auth_scheme_data`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `auth_scheme_data` (
  `id` int NOT NULL AUTO_INCREMENT,
  `auth_scheme` varchar(255) NOT NULL,
  `auth_type` varchar(255) NOT NULL,
  `username` varchar(255) NOT NULL,
  `user_id` int DEFAULT NULL,
  `encrypted_password` varbinary(255) DEFAULT NULL,
  `must_reset_password` tinyint(1) DEFAULT '0',
  `totp_secret` varchar(255) DEFAULT NULL,
  `json_data` json DEFAULT NULL,
  `last_login` timestamp NULL DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `auth_scheme_user_key` (`auth_scheme`,`username`),
  KEY `fk_user_id__users_id` (`user_id`),
  CONSTRAINT `fk_user_id__users_id` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `default_tags`
--

DROP TABLE IF EXISTS `default_tags`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `default_tags` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(63) NOT NULL,
  `color_name` varchar(63) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `email_queue`
--

DROP TABLE IF EXISTS `email_queue`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `email_queue` (
  `id` int NOT NULL AUTO_INCREMENT,
  `to_email` varchar(255) NOT NULL,
  `user_id` int NOT NULL DEFAULT '0',
  `template` varchar(255) NOT NULL,
  `email_status` varchar(32) NOT NULL DEFAULT 'created',
  `error_count` int NOT NULL DEFAULT '0',
  `error_text` text,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `email_queue__email_status` (`email_status`),
  KEY `email_queue__email_to` (`to_email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `evidence`
--

DROP TABLE IF EXISTS `evidence`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `evidence` (
  `id` int NOT NULL AUTO_INCREMENT,
  `uuid` varchar(36) NOT NULL,
  `operation_id` int NOT NULL,
  `operator_id` int NOT NULL,
  `description` text,
  `content_type` varchar(31) NOT NULL,
  `full_image_key` varchar(255) DEFAULT NULL,
  `thumb_image_key` varchar(255) DEFAULT NULL,
  `occurred_at` timestamp NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uuid` (`uuid`),
  KEY `operation_id` (`operation_id`),
  KEY `operator_id` (`operator_id`),
  CONSTRAINT `evidence_ibfk_1` FOREIGN KEY (`operation_id`) REFERENCES `operations` (`id`),
  CONSTRAINT `evidence_ibfk_2` FOREIGN KEY (`operator_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `evidence_finding_map`
--

DROP TABLE IF EXISTS `evidence_finding_map`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `evidence_finding_map` (
  `evidence_id` int NOT NULL,
  `finding_id` int NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`evidence_id`,`finding_id`),
  KEY `event_id` (`finding_id`),
  CONSTRAINT `evidence_finding_map_ibfk_1` FOREIGN KEY (`evidence_id`) REFERENCES `evidence` (`id`) ON DELETE CASCADE,
  CONSTRAINT `evidence_finding_map_ibfk_2` FOREIGN KEY (`finding_id`) REFERENCES `findings` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `evidence_metadata`
--

DROP TABLE IF EXISTS `evidence_metadata`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `evidence_metadata` (
  `id` int NOT NULL AUTO_INCREMENT,
  `evidence_id` int NOT NULL,
  `source` varchar(255) NOT NULL,
  `body` text NOT NULL,
  `status` varchar(255) DEFAULT NULL,
  `last_run_message` text,
  `can_process` tinyint(1) DEFAULT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `work_started_at` timestamp NULL DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `evidence_id` (`evidence_id`,`source`),
  KEY `evidence_id_2` (`evidence_id`),
  CONSTRAINT `evidence_metadata_ibfk_1` FOREIGN KEY (`evidence_id`) REFERENCES `evidence` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `finding_categories`
--

DROP TABLE IF EXISTS `finding_categories`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `finding_categories` (
  `id` int NOT NULL AUTO_INCREMENT,
  `category` varchar(255) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `category` (`category`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `findings`
--

DROP TABLE IF EXISTS `findings`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `findings` (
  `id` int NOT NULL AUTO_INCREMENT,
  `uuid` varchar(36) NOT NULL,
  `operation_id` int NOT NULL,
  `ready_to_report` tinyint(1) NOT NULL DEFAULT '0',
  `ticket_link` varchar(255) DEFAULT NULL,
  `category_id` int DEFAULT NULL,
  `title` varchar(255) NOT NULL,
  `description` text NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uuid` (`uuid`),
  KEY `operation_id` (`operation_id`),
  KEY `fk_category_id__finding_categories_id` (`category_id`),
  CONSTRAINT `findings_ibfk_1` FOREIGN KEY (`operation_id`) REFERENCES `operations` (`id`),
  CONSTRAINT `fk_category_id__finding_categories_id` FOREIGN KEY (`category_id`) REFERENCES `finding_categories` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `gorp_migrations`
--

DROP TABLE IF EXISTS `gorp_migrations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `gorp_migrations` (
  `id` varchar(255) NOT NULL,
  `applied_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `operations`
--

DROP TABLE IF EXISTS `operations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `operations` (
  `id` int NOT NULL AUTO_INCREMENT,
  `slug` varchar(255) NOT NULL,
  `name` varchar(255) NOT NULL,
  `description` varchar(255) DEFAULT NULL,
  `active` tinyint(1) DEFAULT '1',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `slug` (`slug`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `queries`
--

DROP TABLE IF EXISTS `queries`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `queries` (
  `id` int NOT NULL AUTO_INCREMENT,
  `operation_id` int NOT NULL,
  `name` varchar(255) NOT NULL,
  `query` varchar(255) NOT NULL,
  `type` varchar(15) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`,`operation_id`,`type`),
  UNIQUE KEY `query` (`query`,`operation_id`,`type`),
  KEY `operation_id` (`operation_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `service_workers`
--

DROP TABLE IF EXISTS `service_workers`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `service_workers` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `config` json NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `sessions`
--

DROP TABLE IF EXISTS `sessions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sessions` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` int DEFAULT NULL,
  `session_data` longblob,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `modified_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `expires_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `user_id` (`user_id`),
  CONSTRAINT `sessions_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `tag_evidence_map`
--

DROP TABLE IF EXISTS `tag_evidence_map`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `tag_evidence_map` (
  `tag_id` int NOT NULL,
  `evidence_id` int NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`tag_id`,`evidence_id`),
  KEY `evidence_id` (`evidence_id`),
  CONSTRAINT `tag_evidence_map_ibfk_1` FOREIGN KEY (`tag_id`) REFERENCES `tags` (`id`),
  CONSTRAINT `tag_evidence_map_ibfk_2` FOREIGN KEY (`evidence_id`) REFERENCES `evidence` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `tags`
--

DROP TABLE IF EXISTS `tags`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `tags` (
  `id` int NOT NULL AUTO_INCREMENT,
  `operation_id` int NOT NULL,
  `name` varchar(63) NOT NULL,
  `color_name` varchar(63) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`,`operation_id`),
  KEY `operation_id` (`operation_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `user_operation_permissions`
--

DROP TABLE IF EXISTS `user_operation_permissions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user_operation_permissions` (
  `user_id` int NOT NULL,
  `operation_id` int NOT NULL,
  `role` varchar(255) NOT NULL,
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`user_id`,`operation_id`),
  KEY `operation_id` (`operation_id`),
  CONSTRAINT `user_operation_permissions_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`),
  CONSTRAINT `user_operation_permissions_ibfk_2` FOREIGN KEY (`operation_id`) REFERENCES `operations` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `user_operation_preferences`
--

DROP TABLE IF EXISTS `user_operation_preferences`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user_operation_preferences` (
  `user_id` int NOT NULL,
  `operation_id` int NOT NULL,
  `is_favorite` tinyint(1) NOT NULL DEFAULT '0',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`user_id`,`operation_id`),
  KEY `operation_id` (`operation_id`),
  CONSTRAINT `user_operation_preferences_ibfk_1` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`),
  CONSTRAINT `user_operation_preferences_ibfk_2` FOREIGN KEY (`operation_id`) REFERENCES `operations` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `users` (
  `id` int NOT NULL AUTO_INCREMENT,
  `slug` varchar(255) NOT NULL,
  `first_name` varchar(255) NOT NULL,
  `last_name` varchar(255) NOT NULL,
  `email` varchar(255) NOT NULL,
  `admin` tinyint(1) DEFAULT '0',
  `headless` tinyint(1) DEFAULT '0',
  `disabled` tinyint(1) DEFAULT '0',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT NULL,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `slug` (`slug`),
  UNIQUE KEY `unique_email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2022-11-11 22:48:03
-- MySQL dump 10.13  Distrib 8.0.22, for Linux (x86_64)
--
-- Host: localhost    Database: migrate_db
-- ------------------------------------------------------
-- Server version	8.0.22

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Dumping data for table `gorp_migrations`
--

LOCK TABLES `gorp_migrations` WRITE;
/*!40000 ALTER TABLE `gorp_migrations` DISABLE KEYS */;
INSERT INTO `gorp_migrations` VALUES ('20190705190058-create-users-table.sql','2022-11-11 22:47:58'),('20190708185420-create-operations-table.sql','2022-11-11 22:47:58'),('20190708185427-create-events-table.sql','2022-11-11 22:47:58'),('20190708185432-create-evidence-table.sql','2022-11-11 22:47:58'),('20190708185441-create-evidence-event-map-table.sql','2022-11-11 22:47:59'),('20190716190100-create-user-operation-map-table.sql','2022-11-11 22:47:59'),('20190722193434-create-tags-table.sql','2022-11-11 22:47:59'),('20190722193937-create-tag-event-map.sql','2022-11-11 22:47:59'),('20190909183500-add-short-name-to-users-table.sql','2022-11-11 22:47:59'),('20190909190416-add-short-name-index.sql','2022-11-11 22:47:59'),('20190926205116-evidence-name.sql','2022-11-11 22:47:59'),('20190930173342-add-saved-searches.sql','2022-11-11 22:47:59'),('20191001182541-evidence-tags.sql','2022-11-11 22:47:59'),('20191008005212-add-uuid-to-events-evidence.sql','2022-11-11 22:48:00'),('20191015235306-add-slug-to-operations.sql','2022-11-11 22:48:00'),('20191018172105-modular-auth.sql','2022-11-11 22:48:00'),('20191023170906-codeblock.sql','2022-11-11 22:48:00'),('20191101185207-replace-events-with-findings.sql','2022-11-11 22:48:00'),('20191114211948-add-operation-to-tags.sql','2022-11-11 22:48:00'),('20191205182830-create-api-keys-table.sql','2022-11-11 22:48:00'),('20191213222629-users-with-email.sql','2022-11-11 22:48:01'),('20200103194053-rename-short-name-to-slug.sql','2022-11-11 22:48:01'),('20200104013804-rework-ashirt-auth.sql','2022-11-11 22:48:01'),('20200116070736-add-admin-flag.sql','2022-11-11 22:48:01'),('20200130175541-fix-color-truncation.sql','2022-11-11 22:48:01'),('20200205200208-disable-user-support.sql','2022-11-11 22:48:01'),('20200215015330-optional-user-id.sql','2022-11-11 22:48:01'),('20200221195107-deletable-user.sql','2022-11-11 22:48:01'),('20200303215004-move-last-login.sql','2022-11-11 22:48:01'),('20200306221628-add-explicit-headless.sql','2022-11-11 22:48:02'),('20200331155258-finding-status.sql','2022-11-11 22:48:02'),('20200617193248-case-senitive-apikey.sql','2022-11-11 22:48:02'),('20200928160958-add-totp-secret-to-auth-table.sql','2022-11-11 22:48:02'),('20210120205510-create-email-queue-table.sql','2022-11-11 22:48:02'),('20210401220807-dynamic-categories.sql','2022-11-11 22:48:02'),('20210408212206-remove-findings-category.sql','2022-11-11 22:48:02'),('20210730170543-add-auth-type.sql','2022-11-11 22:48:02'),('20220211181557-add-default-tags.sql','2022-11-11 22:48:02'),('20220512174013-evidence-metadata.sql','2022-11-11 22:48:02'),('20220516163424-add-worker-services.sql','2022-11-11 22:48:03'),('20220811153414-webauthn-credentials.sql','2022-11-11 22:48:03'),('20220908193523-switch-to-username.sql','2022-11-11 22:48:03'),('20220912185024-add-is_favorite.sql','2022-11-11 22:48:03'),('20220916190855-remove-null-as-value-for-is_favorite.sql','2022-11-11 22:48:03'),('20221027152757-remove-operation-status.sql','2022-11-11 22:48:03'),('20221111221242-create-user-operation-preferences.sql','2022-11-11 22:48:03');
/*!40000 ALTER TABLE `gorp_migrations` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2022-11-11 22:48:03
