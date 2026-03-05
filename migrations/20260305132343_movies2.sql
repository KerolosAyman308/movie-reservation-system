-- Modify "movie_genres" table
ALTER TABLE `movie_genres` DROP INDEX `fk_movies_genres`, ADD UNIQUE INDEX `idx_movie_genre` (`movie_id`, `genre_id`);
-- Create "files" table
CREATE TABLE `files` (
  `object_key` varchar(40) NOT NULL,
  `bucket_name` varchar(250) NOT NULL,
  `original_name` varchar(250) NOT NULL,
  `file_name` varchar(250) NOT NULL,
  `url` varchar(250) NOT NULL,
  `size` bigint NOT NULL,
  `hash` varchar(65) NOT NULL,
  PRIMARY KEY (`object_key`)
) CHARSET utf8mb4 COLLATE utf8mb4_0900_ai_ci;
-- Modify "movies" table
ALTER TABLE `movies` DROP COLUMN `image`, ADD COLUMN `image_object_key` varchar(40) NULL, ADD INDEX `fk_movies_image` (`image_object_key`), ADD CONSTRAINT `fk_movies_image` FOREIGN KEY (`image_object_key`) REFERENCES `files` (`object_key`) ON UPDATE NO ACTION ON DELETE SET NULL;
