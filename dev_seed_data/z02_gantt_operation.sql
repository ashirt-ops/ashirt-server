
SELECT id INTO @op_owner FROM users WHERE slug='adefaultuser';

INSERT INTO operations(`slug`, `name`, `description`, active, `status`) VALUES
    ('gantt-test',   'Gantt Test',   'An operation viewing overview', 1, 0);

SELECT id INTO @op_id FROM operations WHERE slug='gantt-test';

INSERT INTO user_operation_permissions (`user_id`, `operation_id`, `role`) VALUES (@op_owner, @op_id, 'admin');

SET @tag_roy = 'roy';
SET @tag_iggy = 'iggy';
SET @tag_larry = 'larry';
SET @tag_lemmy = 'lemmy';
SET @tag_ludwig = 'ludwig';
SET @tag_morgon = 'morgon';
SET @tag_wendy = 'wendy';
SET @tag_bowser = 'bowser';
SET @tag_kamek = 'kamek';
SET @tag_bowser_jr = 'bowser jr';

INSERT INTO tags (name, color_name, operation_id) VALUES
      (@tag_roy,      'red',    @op_id)
    , (@tag_iggy,     'orange', @op_id)
    , (@tag_larry,    'yellow', @op_id)
    , (@tag_lemmy,    'green',  @op_id)
    , (@tag_ludwig,   'blue',   @op_id)
    , (@tag_morgon,   'lightRed',    @op_id)
    , (@tag_wendy,    'lightOrange', @op_id)
    , (@tag_bowser,   'lightYellow', @op_id)
    , (@tag_kamek,    'lightGreen',  @op_id)
    , (@tag_bowser_jr,'lightBlue',   @op_id)
    ;

SELECT id into @tag_id_01 from tags where name = @tag_roy;
SELECT id into @tag_id_02 from tags where name = @tag_iggy;
SELECT id into @tag_id_03 from tags where name = @tag_larry;
SELECT id into @tag_id_04 from tags where name = @tag_lemmy;
SELECT id into @tag_id_05 from tags where name = @tag_ludwig;
SELECT id into @tag_id_06 from tags where name = @tag_morgon;
SELECT id into @tag_id_07 from tags where name = @tag_wendy;
SELECT id into @tag_id_08 from tags where name = @tag_bowser;
SELECT id into @tag_id_09 from tags where name = @tag_kamek;
SELECT id into @tag_id_10 from tags where name = @tag_bowser_jr;


SET @evi_uuid_01 = 'f01E0000-0000-4000-8000-000000000000';
SET @evi_uuid_02 = 'f02E0000-0000-4000-8000-000000000000';
SET @evi_uuid_03 = 'f03E0000-0000-4000-8000-000000000000';
SET @evi_uuid_04 = 'f04E0000-0000-4000-8000-000000000000';
SET @evi_uuid_05 = 'f05E0000-0000-4000-8000-000000000000';
SET @evi_uuid_06 = 'f06E0000-0000-4000-8000-000000000000';
SET @evi_uuid_07 = 'f07E0000-0000-4000-8000-000000000000';
SET @evi_uuid_08 = 'f08E0000-0000-4000-8000-000000000000';
SET @evi_uuid_09 = 'f09E0000-0000-4000-8000-000000000000';
SET @evi_uuid_10 = 'f0AE0000-0000-4000-8000-000000000000';
SET @evi_uuid_11 = 'f0BE0000-0000-4000-8000-000000000000';
SET @evi_uuid_12 = 'f0CE0000-0000-4000-8000-000000000000';
SET @evi_uuid_13 = 'f0DE0000-0000-4000-8000-000000000000';
SET @evi_uuid_14 = 'f0EE0000-0000-4000-8000-000000000000';
SET @evi_uuid_15 = 'f0FE0000-0000-4000-8000-000000000000';
SET @evi_uuid_16 = 'f10E0000-0000-4000-8000-000000000000';
SET @evi_uuid_17 = 'f11E0000-0000-4000-8000-000000000000';
SET @evi_uuid_18 = 'f12E0000-0000-4000-8000-000000000000';
SET @evi_uuid_19 = 'f13E0000-0000-4000-8000-000000000000';
SET @evi_uuid_20 = 'f14E0000-0000-4000-8000-000000000000';
SET @evi_uuid_extra = 'f01E0000-1000-4000-8000-000000000000';

INSERT INTO evidence
    (uuid, operation_id, operator_id, content_type, full_image_key, thumb_image_key, occurred_at, `description`)
VALUES
      (@evi_uuid_01, @op_id,    @op_owner, 'none', '', '', now() - INTERVAL 19 DAY, '')
    , (@evi_uuid_extra, @op_id, @op_owner, 'none', '', '', now() - INTERVAL 19 DAY + INTERVAL 12 HOUR, '')
    , (@evi_uuid_02, @op_id,    @op_owner, 'none', '', '', now() - INTERVAL 18 DAY, '')
    , (@evi_uuid_03, @op_id,    @op_owner, 'none', '', '', now() - INTERVAL 17 DAY, '')
    , (@evi_uuid_04, @op_id,    @op_owner, 'none', '', '', now() - INTERVAL 16 DAY, '')
    , (@evi_uuid_05, @op_id,    @op_owner, 'none', '', '', now() - INTERVAL 15 DAY, '')
    , (@evi_uuid_06, @op_id,    @op_owner, 'none', '', '', now() - INTERVAL 14 DAY, '')
    , (@evi_uuid_07, @op_id,    @op_owner, 'none', '', '', now() - INTERVAL 13 DAY, '')
    , (@evi_uuid_08, @op_id,    @op_owner, 'none', '', '', now() - INTERVAL 12 DAY, '')
    , (@evi_uuid_09, @op_id,    @op_owner, 'none', '', '', now() - INTERVAL 11 DAY, '')
    , (@evi_uuid_10, @op_id,    @op_owner, 'none', '', '', now() - INTERVAL 10 DAY, '')
    , (@evi_uuid_11, @op_id,    @op_owner, 'none', '', '', now() - INTERVAL  9 DAY, '')
    , (@evi_uuid_12, @op_id,    @op_owner, 'none', '', '', now() - INTERVAL  8 DAY, '')
    , (@evi_uuid_13, @op_id,    @op_owner, 'none', '', '', now() - INTERVAL  7 DAY, '')
    , (@evi_uuid_14, @op_id,    @op_owner, 'none', '', '', now() - INTERVAL  6 DAY, '')
    , (@evi_uuid_15, @op_id,    @op_owner, 'none', '', '', now() - INTERVAL  5 DAY, '')
    , (@evi_uuid_16, @op_id,    @op_owner, 'none', '', '', now() - INTERVAL  4 DAY, '')
    , (@evi_uuid_17, @op_id,    @op_owner, 'none', '', '', now() - INTERVAL  3 DAY, '')
    , (@evi_uuid_18, @op_id,    @op_owner, 'none', '', '', now() - INTERVAL  2 DAY, '')
    , (@evi_uuid_19, @op_id,    @op_owner, 'none', '', '', now() - INTERVAL  1 DAY, '')
    , (@evi_uuid_20, @op_id,    @op_owner, 'none', '', '', now(), '')
    ;

SELECT id INTO @evi_id_01    from evidence WHERE uuid = @evi_uuid_01;
SELECT id INTO @evi_id_02    from evidence WHERE uuid = @evi_uuid_02;
SELECT id INTO @evi_id_03    from evidence WHERE uuid = @evi_uuid_03;
SELECT id INTO @evi_id_04    from evidence WHERE uuid = @evi_uuid_04;
SELECT id INTO @evi_id_05    from evidence WHERE uuid = @evi_uuid_05;
SELECT id INTO @evi_id_06    from evidence WHERE uuid = @evi_uuid_06;
SELECT id INTO @evi_id_07    from evidence WHERE uuid = @evi_uuid_07;
SELECT id INTO @evi_id_08    from evidence WHERE uuid = @evi_uuid_08;
SELECT id INTO @evi_id_09    from evidence WHERE uuid = @evi_uuid_09;
SELECT id INTO @evi_id_10    from evidence WHERE uuid = @evi_uuid_10;
SELECT id INTO @evi_id_11    from evidence WHERE uuid = @evi_uuid_11;
SELECT id INTO @evi_id_12    from evidence WHERE uuid = @evi_uuid_12;
SELECT id INTO @evi_id_13    from evidence WHERE uuid = @evi_uuid_13;
SELECT id INTO @evi_id_14    from evidence WHERE uuid = @evi_uuid_14;
SELECT id INTO @evi_id_15    from evidence WHERE uuid = @evi_uuid_15;
SELECT id INTO @evi_id_16    from evidence WHERE uuid = @evi_uuid_16;
SELECT id INTO @evi_id_17    from evidence WHERE uuid = @evi_uuid_17;
SELECT id INTO @evi_id_18    from evidence WHERE uuid = @evi_uuid_18;
SELECT id INTO @evi_id_19    from evidence WHERE uuid = @evi_uuid_19;
SELECT id INTO @evi_id_20    from evidence WHERE uuid = @evi_uuid_20;
SELECT id INTO @evi_id_extra from evidence WHERE uuid = @evi_uuid_extra;


-- tags are in a pattern: the first 10 columns are dedicated to a pineapple, the second to an apple
-- -- pineapple      apple
--    1234567890     1234567890
--  1 .###.###..     ......##..
--  2 #########.     .....##...
--  3 #..###..#.     ..##.#.##.
--  4 ..#####...     .####.####
--  5 .#.#.#.#..     .#.#######
--  6 .#######..     .#.#######
--  7 .#.#.#.#..     .#.#######
--  8 .#######..     .##.######
--  9 .#.#.#.#..     ..#######.
-- 10 ..#####...     ...##.##..

INSERT INTO tag_evidence_map
    (tag_id, evidence_id)
VALUES
      (@tag_id_01, @evi_id_02)
    , (@tag_id_01, @evi_id_extra)
    , (@tag_id_01, @evi_id_03)
    , (@tag_id_01, @evi_id_04)
    , (@tag_id_01, @evi_id_06)
    , (@tag_id_01, @evi_id_07)
    , (@tag_id_01, @evi_id_08)
    , (@tag_id_01, @evi_id_17)
    , (@tag_id_01, @evi_id_18)

    , (@tag_id_02, @evi_id_01)
    , (@tag_id_02, @evi_id_02)
    , (@tag_id_02, @evi_id_03)
    , (@tag_id_02, @evi_id_04)
    , (@tag_id_02, @evi_id_05)
    , (@tag_id_02, @evi_id_06)
    , (@tag_id_02, @evi_id_07)
    , (@tag_id_02, @evi_id_08)
    , (@tag_id_02, @evi_id_09)
    , (@tag_id_02, @evi_id_16)
    , (@tag_id_02, @evi_id_17)

    , (@tag_id_03, @evi_id_01)
    , (@tag_id_03, @evi_id_04)
    , (@tag_id_03, @evi_id_05)
    , (@tag_id_03, @evi_id_06)
    , (@tag_id_03, @evi_id_09)
    , (@tag_id_03, @evi_id_13)
    , (@tag_id_03, @evi_id_14)
    , (@tag_id_03, @evi_id_16)
    , (@tag_id_03, @evi_id_18)
    , (@tag_id_03, @evi_id_19)

    , (@tag_id_04, @evi_id_03)
    , (@tag_id_04, @evi_id_04)
    , (@tag_id_04, @evi_id_05)
    , (@tag_id_04, @evi_id_06)
    , (@tag_id_04, @evi_id_07)
    , (@tag_id_04, @evi_id_13)
    , (@tag_id_04, @evi_id_14)
    , (@tag_id_04, @evi_id_15)
    , (@tag_id_04, @evi_id_17)
    , (@tag_id_04, @evi_id_18)
    , (@tag_id_04, @evi_id_19)
    , (@tag_id_04, @evi_id_20)

    , (@tag_id_05, @evi_id_02)
    , (@tag_id_05, @evi_id_04)
    , (@tag_id_05, @evi_id_06)
    , (@tag_id_05, @evi_id_08)
    , (@tag_id_05, @evi_id_12)
    , (@tag_id_05, @evi_id_14)
    , (@tag_id_05, @evi_id_15)
    , (@tag_id_05, @evi_id_16)
    , (@tag_id_05, @evi_id_17)
    , (@tag_id_05, @evi_id_18)
    , (@tag_id_05, @evi_id_19)
    , (@tag_id_05, @evi_id_20)

    , (@tag_id_06, @evi_id_02)
    , (@tag_id_06, @evi_id_03)
    , (@tag_id_06, @evi_id_04)
    , (@tag_id_06, @evi_id_05)
    , (@tag_id_06, @evi_id_06)
    , (@tag_id_06, @evi_id_07)
    , (@tag_id_06, @evi_id_08)
    , (@tag_id_06, @evi_id_12)
    , (@tag_id_06, @evi_id_14)
    , (@tag_id_06, @evi_id_15)
    , (@tag_id_06, @evi_id_16)
    , (@tag_id_06, @evi_id_17)
    , (@tag_id_06, @evi_id_18)
    , (@tag_id_06, @evi_id_19)
    , (@tag_id_06, @evi_id_20)

    , (@tag_id_07, @evi_id_02)
    , (@tag_id_07, @evi_id_04)
    , (@tag_id_07, @evi_id_06)
    , (@tag_id_07, @evi_id_08)
    , (@tag_id_07, @evi_id_12)
    , (@tag_id_07, @evi_id_14)
    , (@tag_id_07, @evi_id_15)
    , (@tag_id_07, @evi_id_16)
    , (@tag_id_07, @evi_id_17)
    , (@tag_id_07, @evi_id_18)
    , (@tag_id_07, @evi_id_19)
    , (@tag_id_07, @evi_id_20)

    , (@tag_id_08, @evi_id_02)
    , (@tag_id_08, @evi_id_03)
    , (@tag_id_08, @evi_id_04)
    , (@tag_id_08, @evi_id_05)
    , (@tag_id_08, @evi_id_06)
    , (@tag_id_08, @evi_id_07)
    , (@tag_id_08, @evi_id_08)
    , (@tag_id_08, @evi_id_12)
    , (@tag_id_08, @evi_id_13)
    , (@tag_id_08, @evi_id_15)
    , (@tag_id_08, @evi_id_16)
    , (@tag_id_08, @evi_id_17)
    , (@tag_id_08, @evi_id_18)
    , (@tag_id_08, @evi_id_19)
    , (@tag_id_08, @evi_id_20)

    , (@tag_id_09, @evi_id_02)
    , (@tag_id_09, @evi_id_04)
    , (@tag_id_09, @evi_id_06)
    , (@tag_id_09, @evi_id_08)
    , (@tag_id_09, @evi_id_13)
    , (@tag_id_09, @evi_id_14)
    , (@tag_id_09, @evi_id_15)
    , (@tag_id_09, @evi_id_16)
    , (@tag_id_09, @evi_id_17)
    , (@tag_id_09, @evi_id_18)
    , (@tag_id_09, @evi_id_19)

    , (@tag_id_10, @evi_id_03)
    , (@tag_id_10, @evi_id_04)
    , (@tag_id_10, @evi_id_05)
    , (@tag_id_10, @evi_id_06)
    , (@tag_id_10, @evi_id_07)
    , (@tag_id_10, @evi_id_14)
    , (@tag_id_10, @evi_id_15)
    , (@tag_id_10, @evi_id_17)
    , (@tag_id_10, @evi_id_18)

    ;
