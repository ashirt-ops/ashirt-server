-- +migrate Up
-- color_name got truncated to 6 characters due to migration mishap
-- this migration sets the length and restores all color names longer
-- than 6 characters to their full name.
--
-- The only color we can't determine is 'lightV' which could have been
-- lightVermilion or lightViolet so we just assume lightVermilion
ALTER TABLE `tags`
  MODIFY `color_name` VARCHAR(63) NOT NULL
;

UPDATE `tags`
  SET `color_name` = CASE `color_name`
    WHEN "vermil" THEN "vermilion"
    WHEN "lightB" THEN "lightBlue"
    WHEN "lightY" THEN "lightYellow"
    WHEN "lightG" THEN "lightGreen"
    WHEN "lightI" THEN "lightIndigo"
    WHEN "lightO" THEN "lightOrange"
    WHEN "lightP" THEN "lightPink"
    WHEN "lightR" THEN "lightRed"
    WHEN "lightT" THEN "lightTeal"
    WHEN "lightV" THEN "lightVermilion"
    ELSE `color_name`
  END
;

-- +migrate Down
UPDATE `tags`
  SET `color_name` = LEFT(`color_name`, 6)
;

ALTER TABLE `tags`
  MODIFY `color_name` VARCHAR(6) NOT NULL
;
