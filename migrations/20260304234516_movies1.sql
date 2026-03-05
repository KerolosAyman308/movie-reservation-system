-- Create "genres" table
CREATE TABLE `genres` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `idx_genres_name` (`name`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "movies" table
CREATE TABLE `movies` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(100) NOT NULL,
  `description` varchar(500) NOT NULL,
  `image` varchar(350) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE INDEX `idx_movies_title` (`title`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Create "movie_genres" table
CREATE TABLE `movie_genres` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `movie_id` bigint unsigned NULL,
  `genre_id` bigint unsigned NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_genres_movies` (`genre_id`),
  INDEX `fk_movies_genres` (`movie_id`),
  CONSTRAINT `fk_genres_movies` FOREIGN KEY (`genre_id`) REFERENCES `genres` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION,
  CONSTRAINT `fk_movies_genres` FOREIGN KEY (`movie_id`) REFERENCES `movies` (`id`) ON UPDATE NO ACTION ON DELETE NO ACTION
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
