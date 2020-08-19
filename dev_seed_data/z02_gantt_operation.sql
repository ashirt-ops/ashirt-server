
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
SET @tag_koopa = 'koopa';

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
    , (@tag_koopa, 'teal', @op_id)
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
SELECT id into @tag_id_11 from tags where name = @tag_koopa;


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

SET @evi_uuid_long_00 = 'f14E0000-0000-4008-8000-000000000000';
SET @evi_uuid_long_01 = 'f14E0000-0001-4008-8000-000000000000';
SET @evi_uuid_long_02 = 'f14E0000-0002-4008-8000-000000000000';
SET @evi_uuid_long_03 = 'f14E0000-0003-4008-8000-000000000000';
SET @evi_uuid_long_04 = 'f14E0000-0004-4008-8000-000000000000';
SET @evi_uuid_long_05 = 'f14E0000-0005-4008-8000-000000000000';
SET @evi_uuid_long_06 = 'f14E0000-0006-4008-8000-000000000000';
SET @evi_uuid_long_07 = 'f14E0000-0007-4008-8000-000000000000';
SET @evi_uuid_long_08 = 'f14E0000-0008-4008-8000-000000000000';
SET @evi_uuid_long_09 = 'f14E0000-0009-4008-8000-000000000000';
SET @evi_uuid_long_10 = 'f14E0000-0010-4008-8000-000000000000';
SET @evi_uuid_long_11 = 'f14E0000-0011-4008-8000-000000000000';
SET @evi_uuid_long_12 = 'f14E0000-0012-4008-8000-000000000000';
SET @evi_uuid_long_13 = 'f14E0000-0013-4008-8000-000000000000';
SET @evi_uuid_long_14 = 'f14E0000-0014-4008-8000-000000000000';
SET @evi_uuid_long_15 = 'f14E0000-0015-4008-8000-000000000000';
SET @evi_uuid_long_16 = 'f14E0000-0016-4008-8000-000000000000';
SET @evi_uuid_long_17 = 'f14E0000-0017-4008-8000-000000000000';
SET @evi_uuid_long_18 = 'f14E0000-0018-4008-8000-000000000000';
SET @evi_uuid_long_19 = 'f14E0000-0019-4008-8000-000000000000';
SET @evi_uuid_long_20 = 'f14E0000-0020-4008-8000-000000000000';
SET @evi_uuid_long_21 = 'f14E0000-0021-4008-8000-000000000000';
SET @evi_uuid_long_22 = 'f14E0000-0022-4008-8000-000000000000';
SET @evi_uuid_long_23 = 'f14E0000-0023-4008-8000-000000000000';
SET @evi_uuid_long_24 = 'f14E0000-0024-4008-8000-000000000000';
SET @evi_uuid_long_25 = 'f14E0000-0025-4008-8000-000000000000';
SET @evi_uuid_long_26 = 'f14E0000-0026-4008-8000-000000000000';
SET @evi_uuid_long_27 = 'f14E0000-0027-4008-8000-000000000000';
SET @evi_uuid_long_28 = 'f14E0000-0028-4008-8000-000000000000';
SET @evi_uuid_long_29 = 'f14E0000-0029-4008-8000-000000000000';
SET @evi_uuid_long_30 = 'f14E0000-0030-4008-8000-000000000000';
SET @evi_uuid_long_31 = 'f14E0000-0031-4008-8000-000000000000';
SET @evi_uuid_long_32 = 'f14E0000-0032-4008-8000-000000000000';
SET @evi_uuid_long_33 = 'f14E0000-0033-4008-8000-000000000000';
SET @evi_uuid_long_34 = 'f14E0000-0034-4008-8000-000000000000';
SET @evi_uuid_long_35 = 'f14E0000-0035-4008-8000-000000000000';
SET @evi_uuid_long_36 = 'f14E0000-0036-4008-8000-000000000000';
SET @evi_uuid_long_37 = 'f14E0000-0037-4008-8000-000000000000';
SET @evi_uuid_long_38 = 'f14E0000-0038-4008-8000-000000000000';
SET @evi_uuid_long_39 = 'f14E0000-0039-4008-8000-000000000000';
SET @evi_uuid_long_40 = 'f14E0000-0040-4008-8000-000000000000';
SET @evi_uuid_long_41 = 'f14E0000-0041-4008-8000-000000000000';
SET @evi_uuid_long_42 = 'f14E0000-0042-4008-8000-000000000000';
SET @evi_uuid_long_43 = 'f14E0000-0043-4008-8000-000000000000';
SET @evi_uuid_long_44 = 'f14E0000-0044-4008-8000-000000000000';
SET @evi_uuid_long_45 = 'f14E0000-0045-4008-8000-000000000000';
SET @evi_uuid_long_46 = 'f14E0000-0046-4008-8000-000000000000';
SET @evi_uuid_long_47 = 'f14E0000-0047-4008-8000-000000000000';
SET @evi_uuid_long_48 = 'f14E0000-0048-4008-8000-000000000000';
SET @evi_uuid_long_49 = 'f14E0000-0049-4008-8000-000000000000';
SET @evi_uuid_long_50 = 'f14E0000-0050-4008-8000-000000000000';
SET @evi_uuid_long_51 = 'f14E0000-0051-4008-8000-000000000000';
SET @evi_uuid_long_52 = 'f14E0000-0052-4008-8000-000000000000';
SET @evi_uuid_long_53 = 'f14E0000-0053-4008-8000-000000000000';
SET @evi_uuid_long_54 = 'f14E0000-0054-4008-8000-000000000000';
SET @evi_uuid_long_55 = 'f14E0000-0055-4008-8000-000000000000';
SET @evi_uuid_long_56 = 'f14E0000-0056-4008-8000-000000000000';
SET @evi_uuid_long_57 = 'f14E0000-0057-4008-8000-000000000000';
SET @evi_uuid_long_58 = 'f14E0000-0058-4008-8000-000000000000';
SET @evi_uuid_long_59 = 'f14E0000-0059-4008-8000-000000000000';


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

, (@evi_uuid_long_00, @op_id, @op_owner, 'none', '', '', now(), '')
, (@evi_uuid_long_01, @op_id, @op_owner, 'none', '', '', now() - INTERVAL 1 DAY, '')
, (@evi_uuid_long_02, @op_id, @op_owner, 'none', '', '', now() - INTERVAL 2 DAY, '')
, (@evi_uuid_long_03, @op_id, @op_owner, 'none', '', '', now() - INTERVAL 3 DAY, '')
, (@evi_uuid_long_04, @op_id, @op_owner, 'none', '', '', now() - INTERVAL 4 DAY, '')
, (@evi_uuid_long_05, @op_id, @op_owner, 'none', '', '', now() - INTERVAL 5 DAY, '')
, (@evi_uuid_long_06, @op_id, @op_owner, 'none', '', '', now() - INTERVAL 6 DAY, '')
, (@evi_uuid_long_07, @op_id, @op_owner, 'none', '', '', now() - INTERVAL 7 DAY, '')
, (@evi_uuid_long_08, @op_id, @op_owner, 'none', '', '', now() - INTERVAL 8 DAY, '')
, (@evi_uuid_long_09, @op_id, @op_owner, 'none', '', '', now() - INTERVAL 9 DAY, '')
, (@evi_uuid_long_10, @op_id, @op_owner, 'none', '', '', now() - INTERVAL 10 DAY, '')
, (@evi_uuid_long_11, @op_id, @op_owner, 'none', '', '', now() - INTERVAL 11 DAY, '')
, (@evi_uuid_long_12, @op_id, @op_owner, 'none', '', '', now() - INTERVAL 12 DAY, '')
, (@evi_uuid_long_13, @op_id, @op_owner, 'none', '', '', now() - INTERVAL 13 DAY, '')
, (@evi_uuid_long_14, @op_id, @op_owner, 'none', '', '', now() - INTERVAL 14 DAY, '')
, (@evi_uuid_long_15, @op_id, @op_owner, 'none', '', '', now() - INTERVAL 15 DAY, '')
, (@evi_uuid_long_16, @op_id, @op_owner, 'none', '', '', now() - INTERVAL 16 DAY, '')
, (@evi_uuid_long_17, @op_id, @op_owner, 'none', '', '', now() - INTERVAL 17 DAY, '')
, (@evi_uuid_long_18, @op_id, @op_owner, 'none', '', '', now() - INTERVAL 18 DAY, '')
, (@evi_uuid_long_19, @op_id, @op_owner, 'none', '', '', now() - INTERVAL 19 DAY, '')
, (@evi_uuid_long_20, @op_id, @op_owner, 'none', '', '', now() - INTERVAL 20 DAY, '')
, (@evi_uuid_long_21, @op_id, @op_owner, 'none', '', '', now() - INTERVAL 21 DAY, '')
, (@evi_uuid_long_22, @op_id, @op_owner, 'none', '', '', now() - INTERVAL 22 DAY, '')
, (@evi_uuid_long_23, @op_id, @op_owner, 'none', '', '', now() - INTERVAL 23 DAY, '')
, (@evi_uuid_long_24, @op_id, @op_owner, 'none', '', '', now() - INTERVAL 24 DAY, '')
, (@evi_uuid_long_25, @op_id, @op_owner, 'none', '', '', now() - INTERVAL 25 DAY, '')
, (@evi_uuid_long_26, @op_id, @op_owner, 'none', '', '', now() - INTERVAL 26 DAY, '')
, (@evi_uuid_long_27, @op_id, @op_owner, 'none', '', '', now() - INTERVAL 27 DAY, '')
, (@evi_uuid_long_28, @op_id, @op_owner, 'none', '', '', now() - INTERVAL 28 DAY, '')
, (@evi_uuid_long_29, @op_id, @op_owner, 'none', '', '', now() - INTERVAL 29 DAY, '')
, (@evi_uuid_long_30, @op_id, @op_owner, 'none', '', '', now() - INTERVAL 30 DAY, '')
, (@evi_uuid_long_31, @op_id, @op_owner, 'none', '', '', now() - INTERVAL 31 DAY, '')
, (@evi_uuid_long_32, @op_id, @op_owner, 'none', '', '', now() - INTERVAL 32 DAY, '')
, (@evi_uuid_long_33, @op_id, @op_owner, 'none', '', '', now() - INTERVAL 33 DAY, '')
, (@evi_uuid_long_34, @op_id, @op_owner, 'none', '', '', now() - INTERVAL 34 DAY, '')
, (@evi_uuid_long_35, @op_id, @op_owner, 'none', '', '', now() - INTERVAL 35 DAY, '')
, (@evi_uuid_long_36, @op_id, @op_owner, 'none', '', '', now() - INTERVAL 36 DAY, '')
, (@evi_uuid_long_37, @op_id, @op_owner, 'none', '', '', now() - INTERVAL 37 DAY, '')
, (@evi_uuid_long_38, @op_id, @op_owner, 'none', '', '', now() - INTERVAL 38 DAY, '')
, (@evi_uuid_long_39, @op_id, @op_owner, 'none', '', '', now() - INTERVAL 39 DAY, '')
, (@evi_uuid_long_40, @op_id, @op_owner, 'none', '', '', now() - INTERVAL 40 DAY, '')
, (@evi_uuid_long_41, @op_id, @op_owner, 'none', '', '', now() - INTERVAL 41 DAY, '')
, (@evi_uuid_long_42, @op_id, @op_owner, 'none', '', '', now() - INTERVAL 42 DAY, '')
, (@evi_uuid_long_43, @op_id, @op_owner, 'none', '', '', now() - INTERVAL 43 DAY, '')
, (@evi_uuid_long_44, @op_id, @op_owner, 'none', '', '', now() - INTERVAL 44 DAY, '')
, (@evi_uuid_long_45, @op_id, @op_owner, 'none', '', '', now() - INTERVAL 45 DAY, '')
, (@evi_uuid_long_46, @op_id, @op_owner, 'none', '', '', now() - INTERVAL 46 DAY, '')
, (@evi_uuid_long_47, @op_id, @op_owner, 'none', '', '', now() - INTERVAL 47 DAY, '')
, (@evi_uuid_long_48, @op_id, @op_owner, 'none', '', '', now() - INTERVAL 48 DAY, '')
, (@evi_uuid_long_49, @op_id, @op_owner, 'none', '', '', now() - INTERVAL 49 DAY, '')
, (@evi_uuid_long_50, @op_id, @op_owner, 'none', '', '', now() - INTERVAL 50 DAY, '')
, (@evi_uuid_long_51, @op_id, @op_owner, 'none', '', '', now() - INTERVAL 51 DAY, '')
, (@evi_uuid_long_52, @op_id, @op_owner, 'none', '', '', now() - INTERVAL 52 DAY, '')
, (@evi_uuid_long_53, @op_id, @op_owner, 'none', '', '', now() - INTERVAL 53 DAY, '')
, (@evi_uuid_long_54, @op_id, @op_owner, 'none', '', '', now() - INTERVAL 54 DAY, '')
, (@evi_uuid_long_55, @op_id, @op_owner, 'none', '', '', now() - INTERVAL 55 DAY, '')
, (@evi_uuid_long_56, @op_id, @op_owner, 'none', '', '', now() - INTERVAL 56 DAY, '')
, (@evi_uuid_long_57, @op_id, @op_owner, 'none', '', '', now() - INTERVAL 57 DAY, '')
, (@evi_uuid_long_58, @op_id, @op_owner, 'none', '', '', now() - INTERVAL 58 DAY, '')
, (@evi_uuid_long_59, @op_id, @op_owner, 'none', '', '', now() - INTERVAL 59 DAY, '')

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


SELECT id INTO @evi_uuid_long_00_id from evidence where uuid = @evi_uuid_long_00;
SELECT id INTO @evi_uuid_long_01_id from evidence where uuid = @evi_uuid_long_01;
SELECT id INTO @evi_uuid_long_02_id from evidence where uuid = @evi_uuid_long_02;
SELECT id INTO @evi_uuid_long_03_id from evidence where uuid = @evi_uuid_long_03;
SELECT id INTO @evi_uuid_long_04_id from evidence where uuid = @evi_uuid_long_04;
SELECT id INTO @evi_uuid_long_05_id from evidence where uuid = @evi_uuid_long_05;
SELECT id INTO @evi_uuid_long_06_id from evidence where uuid = @evi_uuid_long_06;
SELECT id INTO @evi_uuid_long_07_id from evidence where uuid = @evi_uuid_long_07;
SELECT id INTO @evi_uuid_long_08_id from evidence where uuid = @evi_uuid_long_08;
SELECT id INTO @evi_uuid_long_09_id from evidence where uuid = @evi_uuid_long_09;
SELECT id INTO @evi_uuid_long_10_id from evidence where uuid = @evi_uuid_long_10;
SELECT id INTO @evi_uuid_long_11_id from evidence where uuid = @evi_uuid_long_11;
SELECT id INTO @evi_uuid_long_12_id from evidence where uuid = @evi_uuid_long_12;
SELECT id INTO @evi_uuid_long_13_id from evidence where uuid = @evi_uuid_long_13;
SELECT id INTO @evi_uuid_long_14_id from evidence where uuid = @evi_uuid_long_14;
SELECT id INTO @evi_uuid_long_15_id from evidence where uuid = @evi_uuid_long_15;
SELECT id INTO @evi_uuid_long_16_id from evidence where uuid = @evi_uuid_long_16;
SELECT id INTO @evi_uuid_long_17_id from evidence where uuid = @evi_uuid_long_17;
SELECT id INTO @evi_uuid_long_18_id from evidence where uuid = @evi_uuid_long_18;
SELECT id INTO @evi_uuid_long_19_id from evidence where uuid = @evi_uuid_long_19;
SELECT id INTO @evi_uuid_long_20_id from evidence where uuid = @evi_uuid_long_20;
SELECT id INTO @evi_uuid_long_21_id from evidence where uuid = @evi_uuid_long_21;
SELECT id INTO @evi_uuid_long_22_id from evidence where uuid = @evi_uuid_long_22;
SELECT id INTO @evi_uuid_long_23_id from evidence where uuid = @evi_uuid_long_23;
SELECT id INTO @evi_uuid_long_24_id from evidence where uuid = @evi_uuid_long_24;
SELECT id INTO @evi_uuid_long_25_id from evidence where uuid = @evi_uuid_long_25;
SELECT id INTO @evi_uuid_long_26_id from evidence where uuid = @evi_uuid_long_26;
SELECT id INTO @evi_uuid_long_27_id from evidence where uuid = @evi_uuid_long_27;
SELECT id INTO @evi_uuid_long_28_id from evidence where uuid = @evi_uuid_long_28;
SELECT id INTO @evi_uuid_long_29_id from evidence where uuid = @evi_uuid_long_29;
SELECT id INTO @evi_uuid_long_30_id from evidence where uuid = @evi_uuid_long_30;
SELECT id INTO @evi_uuid_long_31_id from evidence where uuid = @evi_uuid_long_31;
SELECT id INTO @evi_uuid_long_32_id from evidence where uuid = @evi_uuid_long_32;
SELECT id INTO @evi_uuid_long_33_id from evidence where uuid = @evi_uuid_long_33;
SELECT id INTO @evi_uuid_long_34_id from evidence where uuid = @evi_uuid_long_34;
SELECT id INTO @evi_uuid_long_35_id from evidence where uuid = @evi_uuid_long_35;
SELECT id INTO @evi_uuid_long_36_id from evidence where uuid = @evi_uuid_long_36;
SELECT id INTO @evi_uuid_long_37_id from evidence where uuid = @evi_uuid_long_37;
SELECT id INTO @evi_uuid_long_38_id from evidence where uuid = @evi_uuid_long_38;
SELECT id INTO @evi_uuid_long_39_id from evidence where uuid = @evi_uuid_long_39;
SELECT id INTO @evi_uuid_long_40_id from evidence where uuid = @evi_uuid_long_40;
SELECT id INTO @evi_uuid_long_41_id from evidence where uuid = @evi_uuid_long_41;
SELECT id INTO @evi_uuid_long_42_id from evidence where uuid = @evi_uuid_long_42;
SELECT id INTO @evi_uuid_long_43_id from evidence where uuid = @evi_uuid_long_43;
SELECT id INTO @evi_uuid_long_44_id from evidence where uuid = @evi_uuid_long_44;
SELECT id INTO @evi_uuid_long_45_id from evidence where uuid = @evi_uuid_long_45;
SELECT id INTO @evi_uuid_long_46_id from evidence where uuid = @evi_uuid_long_46;
SELECT id INTO @evi_uuid_long_47_id from evidence where uuid = @evi_uuid_long_47;
SELECT id INTO @evi_uuid_long_48_id from evidence where uuid = @evi_uuid_long_48;
SELECT id INTO @evi_uuid_long_49_id from evidence where uuid = @evi_uuid_long_49;
SELECT id INTO @evi_uuid_long_50_id from evidence where uuid = @evi_uuid_long_50;
SELECT id INTO @evi_uuid_long_51_id from evidence where uuid = @evi_uuid_long_51;
SELECT id INTO @evi_uuid_long_52_id from evidence where uuid = @evi_uuid_long_52;
SELECT id INTO @evi_uuid_long_53_id from evidence where uuid = @evi_uuid_long_53;
SELECT id INTO @evi_uuid_long_54_id from evidence where uuid = @evi_uuid_long_54;
SELECT id INTO @evi_uuid_long_55_id from evidence where uuid = @evi_uuid_long_55;
SELECT id INTO @evi_uuid_long_56_id from evidence where uuid = @evi_uuid_long_56;
SELECT id INTO @evi_uuid_long_57_id from evidence where uuid = @evi_uuid_long_57;
SELECT id INTO @evi_uuid_long_58_id from evidence where uuid = @evi_uuid_long_58;
SELECT id INTO @evi_uuid_long_59_id from evidence where uuid = @evi_uuid_long_59;

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

    , (@tag_id_11, @evi_uuid_long_00_id)
    , (@tag_id_11, @evi_uuid_long_01_id)
    , (@tag_id_11, @evi_uuid_long_02_id)
    , (@tag_id_11, @evi_uuid_long_03_id)
    , (@tag_id_11, @evi_uuid_long_04_id)
    , (@tag_id_11, @evi_uuid_long_05_id)
    , (@tag_id_11, @evi_uuid_long_06_id)
    , (@tag_id_11, @evi_uuid_long_07_id)
    , (@tag_id_11, @evi_uuid_long_08_id)
    , (@tag_id_11, @evi_uuid_long_09_id)
    , (@tag_id_11, @evi_uuid_long_10_id)
    , (@tag_id_11, @evi_uuid_long_11_id)
    , (@tag_id_11, @evi_uuid_long_12_id)
    , (@tag_id_11, @evi_uuid_long_13_id)
    , (@tag_id_11, @evi_uuid_long_14_id)
    , (@tag_id_11, @evi_uuid_long_15_id)
    , (@tag_id_11, @evi_uuid_long_16_id)
    , (@tag_id_11, @evi_uuid_long_17_id)
    , (@tag_id_11, @evi_uuid_long_18_id)
    , (@tag_id_11, @evi_uuid_long_19_id)
    , (@tag_id_11, @evi_uuid_long_20_id)
    , (@tag_id_11, @evi_uuid_long_21_id)
    , (@tag_id_11, @evi_uuid_long_22_id)
    , (@tag_id_11, @evi_uuid_long_23_id)
    , (@tag_id_11, @evi_uuid_long_24_id)
    , (@tag_id_11, @evi_uuid_long_25_id)
    , (@tag_id_11, @evi_uuid_long_26_id)
    , (@tag_id_11, @evi_uuid_long_27_id)
    , (@tag_id_11, @evi_uuid_long_28_id)
    , (@tag_id_11, @evi_uuid_long_29_id)
    , (@tag_id_11, @evi_uuid_long_30_id)
    , (@tag_id_11, @evi_uuid_long_31_id)
    , (@tag_id_11, @evi_uuid_long_32_id)
    , (@tag_id_11, @evi_uuid_long_33_id)
    , (@tag_id_11, @evi_uuid_long_34_id)
    , (@tag_id_11, @evi_uuid_long_35_id)
    , (@tag_id_11, @evi_uuid_long_36_id)
    , (@tag_id_11, @evi_uuid_long_37_id)
    , (@tag_id_11, @evi_uuid_long_38_id)
    , (@tag_id_11, @evi_uuid_long_39_id)
    , (@tag_id_11, @evi_uuid_long_40_id)
    , (@tag_id_11, @evi_uuid_long_41_id)
    , (@tag_id_11, @evi_uuid_long_42_id)
    , (@tag_id_11, @evi_uuid_long_43_id)
    , (@tag_id_11, @evi_uuid_long_44_id)
    , (@tag_id_11, @evi_uuid_long_45_id)
    , (@tag_id_11, @evi_uuid_long_46_id)
    , (@tag_id_11, @evi_uuid_long_47_id)
    , (@tag_id_11, @evi_uuid_long_48_id)
    , (@tag_id_11, @evi_uuid_long_49_id)
    , (@tag_id_11, @evi_uuid_long_50_id)
    , (@tag_id_11, @evi_uuid_long_51_id)
    , (@tag_id_11, @evi_uuid_long_52_id)
    , (@tag_id_11, @evi_uuid_long_53_id)
    , (@tag_id_11, @evi_uuid_long_54_id)
    , (@tag_id_11, @evi_uuid_long_55_id)
    , (@tag_id_11, @evi_uuid_long_56_id)
    , (@tag_id_11, @evi_uuid_long_57_id)
    , (@tag_id_11, @evi_uuid_long_58_id)
    , (@tag_id_11, @evi_uuid_long_59_id)

    ;
