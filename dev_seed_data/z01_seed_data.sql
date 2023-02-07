
-- encrypted bcrypt password (Plain Text: "password")
-- equivalent Go code: bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

SET @defaultPassword = '$2a$10$MLooRpCcdyyxoXwe3ZiCFuQZfsGeVC7TPCSyYhTs8Bl/sFPd4K67W';

SET @alice_uid = 1;
SET @bob_uid = 2;
SET @carl_uid = 3;
SET @debbie_uid = 4;
SET @edith_uid = 5;
SET @fred_uid = 6;
SET @ginna_uid = 7;
SET @hazel_uid = 8;
SET @ivy_uid = 9;
SET @justin_uid = 10;
SET @kevin_uid = 11;
SET @liz_uid = 12;
SET @martin_uid = 13;
SET @nick_uid = 14;
SET @ophelia_uid = 15;
SET @paul_uid = 16;
SET @quinton_uid = 17;
SET @rachel_uid = 18;
SET @tess_uid = 19;
SET @ursula_uid = 20;
SET @vera_uid = 21;
SET @winston_uid = 22;
SET @xavier_uid = 23;
SET @yavonne_uid = 24;
SET @zoey_uid = 25;
SET @smoker_uid = 26;

INSERT INTO users
    (id, slug, first_name, last_name, email, `admin`)
VALUES
      (@alice_uid,   'adefaultuser', 'Alice',   'DefaultUser', 'adefaultuser@example.com', true)
    , (@bob_uid,     'bsnooper',     'Bob',     'Snooper',     'bsnooper@example.com', false)
    , (@carl_uid,    'cother',       'Carl',    'Other',       'cother@example.com', false)
    , (@debbie_uid,  'dother',       'Debbie',  'Other',       'dother@example.com', false)
    , (@edith_uid,   'eother',       'Edith',   'Other',       'eother@example.com', false)
    , (@fred_uid,    'fother',       'Fred',    'Other',       'fother@example.com', false)
    , (@ginna_uid,   'gother',       'Ginna',   'Other',       'gother@example.com', false)
    , (@hazel_uid,   'hother',       'Hazel',   'Other',       'hother@example.com', false)
    , (@ivy_uid,     'iother',       'Ivy',     'Other',       'iother@example.com', false)
    , (@justin_uid,  'jother',       'Justin',  'Other',       'jother@example.com', false)
    , (@kevin_uid,   'kother',       'Kevin',   'Other',       'kother@example.com', false)
    , (@liz_uid,     'lother',       'Liz',     'Other',       'lother@example.com', false)
    , (@martin_uid,  'mother',       'Martin',  'Other',       'mother@example.com', false)
    , (@nick_uid,    'nother',       'Nick',    'Other',       'nother@example.com', false)
    , (@ophelia_uid, 'oother',       'Ophelia', 'Other',       'oother@example.com', false)
    , (@paul_uid,    'pother',       'Paul',    'Other',       'pother@example.com', false)
    , (@quinton_uid, 'qother',       'Quinton', 'Other',       'qother@example.com', false)
    , (@rachel_uid,  'rother',       'Rachel',  'Other',       'rother@example.com', false)
    , (@tess_uid,    'tother',       'Tess',    'Other',       'tother@example.com', false)
    , (@ursula_uid,  'uother',       'Ursula',  'Other',       'uother@example.com', false)
    , (@vera_uid,    'vother',       'Vera',    'Other',       'vother@example.com', false)
    , (@winston_uid, 'wother',       'Winston', 'Other',       'wother@example.com', false)
    , (@xavier_uid,  'xother',       'Xavier',  'Other',       'xother@example.com', false)
    , (@yavonne_uid, 'yother',       'Yavonne', 'Other',       'yother@example.com', false)
    , (@zoey_uid,    'zother',       'Zoey',    'Other',       'zother@example.com', false)
    , (@smoker_uid,  'smoker',       'Smoke',   'Tester',      'smoker', false)
    ;

INSERT INTO auth_scheme_data
    (auth_scheme, auth_type, username, user_id, encrypted_password)
VALUES
      ('local', 'local', 'adefaultuser@example.com', @alice_uid,     @defaultPassword)
    , ('local', 'local', 'bsnooper@example.com',     @bob_uid,       @defaultPassword)
    , ('local', 'local', 'Carl',                     @carl_uid,      @defaultPassword)
    , ('local', 'local', 'Debbie',                   @debbie_uid,    @defaultPassword)
    , ('local', 'local', 'Edith',                    @edith_uid,     @defaultPassword)
    , ('local', 'local', 'Fred',                     @fred_uid,      @defaultPassword)
    , ('local', 'local', 'Ginna',                    @ginna_uid,     @defaultPassword)
    , ('local', 'local', 'Hazel',                    @hazel_uid,     @defaultPassword)
    , ('local', 'local', 'Ivy',                      @ivy_uid,       @defaultPassword)
    , ('local', 'local', 'Justin',                   @justin_uid,    @defaultPassword)
    , ('local', 'local', 'Kevin',                    @kevin_uid,     @defaultPassword)
    , ('local', 'local', 'Liz',                      @liz_uid,       @defaultPassword)
    , ('local', 'local', 'Martin',                   @martin_uid,    @defaultPassword)
    , ('local', 'local', 'Nick',                     @nick_uid,      @defaultPassword)
    , ('local', 'local', 'Ophelia',                  @ophelia_uid,   @defaultPassword)
    , ('local', 'local', 'Paul',                     @paul_uid,      @defaultPassword)
    , ('local', 'local', 'Quinton',                  @quinton_uid,   @defaultPassword)
    , ('local', 'local', 'Rachel',                   @rachel_uid,    @defaultPassword)
    , ('local', 'local', 'Tess',                     @tess_uid,      @defaultPassword)
    , ('local', 'local', 'Ursula',                   @ursula_uid,    @defaultPassword)
    , ('local', 'local', 'Vera',                     @vera_uid,      @defaultPassword)
    , ('local', 'local', 'Winston',                  @winston_uid,   @defaultPassword)
    , ('local', 'local', 'Xavier',                   @xavier_uid,    @defaultPassword)
    , ('local', 'local', 'Yavonne',                  @yavonne_uid,   @defaultPassword)
    , ('local', 'local', 'Zoey',                     @zoey_uid,      @defaultPassword)
    ;

SET @alice_op_id = 1;
SET @bob_op_id = 2;
SET @co_op_id = 3;
SET @no_op_id = 4;
SET @big_op_id = 5;

INSERT INTO operations
    (`id`, `slug`, `name`, `description`, active)
VALUES
      (@alice_op_id, 'alice-op',     'AliceOp',      'An operation for Alice',  1)
    , (@bob_op_id,   'bob-op',       'BobOp',        'An operation for Bob',    1)
    , (@co_op_id,    'co-op',        'Co-Op',        'A Cooperative Operation', 1)
    , (@no_op_id,    'no-op',        'No-Op',        'An Orphaned Operation',   1)
    , (@big_op_id,   'big-op',       'BigOp',        'An operation with lots of users', 1)
    ;

INSERT INTO user_operation_permissions
    (`user_id`, `operation_id`, `role`)
VALUES
      (@alice_uid,   @alice_op_id, 'admin')
    , (@bob_uid,     @bob_op_id,   'admin')
    , (@alice_uid,   @co_op_id,    'admin')
    , (@bob_uid,     @co_op_id,    'admin')
    , (@alice_uid,   @big_op_id,   'admin')
    , (@carl_uid,    @big_op_id,   'read')
    , (@debbie_uid,  @big_op_id,   'write')
    , (@edith_uid,   @big_op_id,   'read')
    , (@fred_uid,    @big_op_id,   'write')
    , (@ginna_uid,   @big_op_id,   'read')
    , (@hazel_uid,   @big_op_id,   'read')
    , (@ivy_uid,     @big_op_id,   'write')
    , (@justin_uid,  @big_op_id,   'read')
    , (@kevin_uid,   @big_op_id,   'write')
    , (@liz_uid,     @big_op_id,   'read')
    , (@martin_uid,  @big_op_id,   'read')
    , (@nick_uid,    @big_op_id,   'read')
    , (@ophelia_uid, @big_op_id,   'write')
    , (@paul_uid,    @big_op_id,   'read')
    , (@quinton_uid, @big_op_id,   'write')
    , (@rachel_uid,  @big_op_id,   'read')
    , (@tess_uid,    @big_op_id,   'write')
    , (@ursula_uid,  @big_op_id,   'write')
    , (@vera_uid,    @big_op_id,   'read')
    , (@winston_uid, @big_op_id,   'write')
    , (@xavier_uid,  @big_op_id,   'read')
    , (@yavonne_uid, @big_op_id,   'write')
    , (@zoey_uid,    @big_op_id,   'read')
    ;

-- Alice op tags
SET @tag_moon_1_id = 1;
SET @tag_moon_2_id = 2;
SET @tag_moon_3_id = 3;
SET @tag_moon_4_id = 4;
SET @tag_moon_5_id = 5;

-- Bob op tags
SET @tag_dwarf_1_id = 6;
SET @tag_dwarf_2_id = 7;
SET @tag_dwarf_3_id = 8;
SET @tag_dwarf_4_id = 9;
SET @tag_dwarf_5_id = 10;
SET @tag_dwarf_6_id = 11;
SET @tag_dwarf_7_id = 12;

-- Co op tags
SET @tag_osi_1_id = 13;
SET @tag_osi_2_id = 14;
SET @tag_osi_3_id = 15;
SET @tag_osi_4_id = 16;
SET @tag_osi_5_id = 17;
SET @tag_osi_6_id = 18;
SET @tag_osi_7_id = 19;

INSERT INTO tags
    (`id`, `name`, color_name, operation_id)
VALUES
      (@tag_moon_1_id, 'Europa', 'red',    @alice_op_id)
    , (@tag_moon_2_id, 'Titan',  'orange', @alice_op_id)
    , (@tag_moon_3_id, 'Io',     'yellow', @alice_op_id)
    , (@tag_moon_4_id, 'Ceres',  'green',  @alice_op_id)
    , (@tag_moon_5_id, 'Triton', 'blue',   @alice_op_id)

    , (@tag_dwarf_1_id, 'Doc',     'red',    @bob_op_id)
    , (@tag_dwarf_2_id, 'Grumpy',  'orange', @bob_op_id)
    , (@tag_dwarf_3_id, 'Happy',   'yellow', @bob_op_id)
    , (@tag_dwarf_4_id, 'Sleepy',  'green',  @bob_op_id)
    , (@tag_dwarf_5_id, 'Bashful', 'blue',   @bob_op_id)
    , (@tag_dwarf_6_id, 'Sneezy',  'indigo', @bob_op_id)
    , (@tag_dwarf_7_id, 'Dopey',   'violet', @bob_op_id)

    , (@tag_osi_1_id, 'application',  'lightRed',    @co_op_id)
    , (@tag_osi_2_id, 'presentation', 'lightOrange', @co_op_id)
    , (@tag_osi_3_id, 'session',      'lightYellow', @co_op_id)
    , (@tag_osi_4_id, 'transport',    'lightGreen',  @co_op_id)
    , (@tag_osi_5_id, 'network',      'lightBlue',   @co_op_id)
    , (@tag_osi_6_id, 'data link',    'lightIndigo', @co_op_id)
    , (@tag_osi_7_id, 'physical',     'lightViolet', @co_op_id)
    ;

--                      0 1 2 3  4 5  6 7  8 9  A B C D E F
SET @a_op_evi_uuid_1 = 'a10E0000-0000-4000-a000-000000000000';
SET @a_op_evi_uuid_2 = 'a20E0000-0000-4000-a000-000000000000';
SET @a_op_evi_uuid_3 = 'a30E0000-c0de-4000-a000-000000000000';

SET @b_op_evi_uuid_1 = 'b10E0000-0000-4000-b000-000000000000';
SET @b_op_evi_uuid_2 = 'b20E0000-0000-4000-b000-000000000000';
SET @b_op_evi_uuid_3 = 'b30E0000-c0de-4000-b000-000000000000';

SET @c_op_evi_uuid_1 = 'c10E0000-0000-4000-8000-000000000000';
SET @c_op_evi_uuid_2 = 'c20E0000-0000-4000-8000-000000000000';
SET @c_op_evi_uuid_3 = 'c30E0000-c0de-4000-8000-000000000000';

SET @a_op_evi_1 = 1;
SET @a_op_evi_2 = 2;
SET @a_op_evi_3 = 3;
SET @b_op_evi_1 = 4;
SET @b_op_evi_2 = 5;
SET @b_op_evi_3 = 6;
SET @c_op_evi_1 = 7;
SET @c_op_evi_2 = 8;
SET @c_op_evi_3 = 9;

INSERT INTO evidence
    (`id`, uuid, operation_id, operator_id, content_type, full_image_key, thumb_image_key, occurred_at, `description`)
VALUES
      (@a_op_evi_1, @a_op_evi_uuid_1, @alice_op_id, @alice_uid, 'image',     'seed_movie_full',    'seed_movie_thumb',    now(), CONCAT_WS(CHAR(10 using utf8), '# Movie Reel', '', 'Cinema''s favorite feature. Action. Excitement. Comedy.', 'Reviews:', '', '* A loud, long and pointless spectacle.', '* Schumacher''s storytelling is limp, and the characters lack energy.', '', '[Click here](https://www.rottentomatoes.com/m/1077027_batman_and_robin) for more info' ))
    , (@a_op_evi_2, @a_op_evi_uuid_2, @alice_op_id, @alice_uid, 'image',     'seed_popcorn_full',  'seed_popcorn_thumb',  now(), 'Popcorn Box')
    , (@a_op_evi_3, @a_op_evi_uuid_3, @alice_op_id, @alice_uid, 'codeblock', 'seed_go_aoc201614',  'seed_go_aoc201614',   now(), 'Go AOC 2016 Day 14 (https://adventofcode.com/2016/day/14)')

    , (@b_op_evi_1, @b_op_evi_uuid_1, @bob_op_id,   @bob_uid,   'image',     'seed_magazine_full', 'seed_magazine_thumb', now(), 'Magazine')
    , (@b_op_evi_2, @b_op_evi_uuid_2, @bob_op_id,   @bob_uid,   'image',     'seed_chicken_full',  'seed_chicken_thumb',  now(), 'Chickens with a Frisbee')
    , (@b_op_evi_3, @b_op_evi_uuid_3, @bob_op_id,   @bob_uid,   'codeblock', 'seed_py_aoc201717',  'seed_py_aoc201717',   now(), 'Python AOC 2017 Day 17 (https://adventofcode.com/2017/day/17)')

    , (@c_op_evi_1, @c_op_evi_uuid_1, @co_op_id,    @alice_uid, 'image',     'seed_pocky_full',    'seed_pocky_thumb',    now(), 'A miko holding a bunch of ofudas')
    , (@c_op_evi_2, @c_op_evi_uuid_2, @co_op_id,    @bob_uid,   'image',     'seed_rocky_full',    'seed_rocky_thumb',    now(), 'A raccoon juggling some leaves')
    , (@c_op_evi_3, @c_op_evi_uuid_3, @co_op_id,    @bob_uid,   'codeblock', 'seed_rs_aoc201501',  'seed_rs_aoc201501',   now(), 'Rust AOC 2015 Day 1 (https://adventofcode.com/2015/day/1)')
    ;

SET @finding_c_opsec_id = 1;
SET @finding_c_dectection_gap_id = 2;

INSERT INTO finding_categories
    (`id`, `category`)
VALUES
      (@finding_c_opsec_id, 'OPSEC')
    , (@finding_c_dectection_gap_id, 'Detection Gap')
    ;

SET @a_op_finding_1 = 1;
SET @a_op_finding_2 = 2;
SET @b_op_finding_1 = 3;
SET @b_op_finding_2 = 4;
SET @c_op_finding_1 = 5;
SET @c_op_finding_2 = 6;

SET @a_op_finding_uuid_1 = 'a10F0000-0000-4000-a000-000000000000';
SET @a_op_finding_uuid_2 = 'a20F0000-0000-4000-a000-000000000000';

SET @b_op_finding_uuid_1 = 'b10F0000-0000-4000-b000-000000000000';
SET @b_op_finding_uuid_2 = 'b20F0000-0000-4000-b000-000000000000';

SET @c_op_finding_uuid_1 = 'c10F0000-0000-4000-8000-000000000000';
SET @c_op_finding_uuid_2 = 'c20F0000-0000-4000-8000-000000000000';

INSERT INTO findings
    (`id`, `uuid`, `operation_id`, `category_id`, `title`, `description`, `ready_to_report`, `ticket_link`)
VALUES
      (@a_op_finding_1, @a_op_finding_uuid_1, @alice_op_id, @finding_category_opsec_id, 'Main Event',                'body', true,  'http://google.com') -- 1
    , (@a_op_finding_2, @a_op_finding_uuid_2, @alice_op_id, @finding_category_opsec_id, 'Side Show left',            'body', true,  null) -- 2
    , (@b_op_finding_1, @b_op_finding_uuid_1, @bob_op_id,   @finding_category_opsec_id, 'Bob Sees an Issue',         'body', false, null) -- 3
    , (@b_op_finding_2, @b_op_finding_uuid_2, @bob_op_id,   @finding_category_opsec_id, 'Bob Suspects Fowl Play',    'body', false, null) -- 4
    , (@c_op_finding_1, @c_op_finding_uuid_1, @co_op_id,    @finding_category_opsec_id, 'I get Pocky',               'body', false, null) -- 5
    , (@c_op_finding_2, @c_op_finding_uuid_2, @co_op_id,    @finding_category_opsec_id, 'Bob gets stuck with Rocky', 'body', false, null) -- 6
    ;


INSERT INTO evidence_finding_map 
    (`evidence_id`, `finding_id`)
VALUES
      (@a_op_evi_1, @a_op_finding_1)
    , (@a_op_evi_2, @a_op_finding_2)
    , (@a_op_evi_3, @a_op_finding_2)
    , (@b_op_evi_1, @b_op_finding_1)
    , (@b_op_evi_2, @b_op_finding_2)
    , (@b_op_evi_3, @b_op_finding_2)
    , (@c_op_evi_1, @c_op_finding_1)
    , (@c_op_evi_2, @c_op_finding_2)
    , (@c_op_evi_3, @c_op_finding_2)
    ;

INSERT INTO tag_evidence_map
    (`tag_id`, `evidence_id`)
VALUES
      (@tag_moon_1_id, @a_op_evi_1)
    , (@tag_moon_2_id, @a_op_evi_1)
    , (@tag_moon_3_id, @a_op_evi_2)
    , (@tag_moon_2_id, @a_op_evi_2)
    , (@tag_moon_4_id, @a_op_evi_3)
    , (@tag_moon_5_id, @a_op_evi_3)

    , (@tag_dwarf_4_id, @b_op_evi_1)
    , (@tag_dwarf_5_id, @b_op_evi_1)
    , (@tag_dwarf_1_id, @b_op_evi_2)
    , (@tag_dwarf_4_id, @b_op_evi_2)
    , (@tag_dwarf_7_id, @b_op_evi_3)
    , (@tag_dwarf_6_id, @b_op_evi_3)

    , (@tag_osi_6_id, @c_op_evi_1)
    , (@tag_osi_5_id, @c_op_evi_1)
    , (@tag_osi_7_id, @c_op_evi_2)
    , (@tag_osi_1_id, @c_op_evi_2)
    , (@tag_osi_3_id, @c_op_evi_3)
    , (@tag_osi_2_id, @c_op_evi_3)
    ;

SET @api_key = 'DAYPFGHnm1Pqes-l0Fm76_y1';
SET @secret_key = 0x1EA9AE5B294BCE747EB6A4A8B5900E73EC38EDBB9215A28A4C9A33A5711892E3418AE449830DCD7893AE54FEA46D00509A2613A9801A88829B70ED41C5C1F9D9;
-- secret_key: HqmuWylLznR+tqSotZAOc+w47buSFaKKTJozpXEYkuNBiuRJgw3NeJOuVP6kbQBQmiYTqYAaiIKbcO1BxcH52Q==

INSERT INTO api_keys
    (`user_id`, `access_key`, `secret_key`)
VALUES
      (@smoker_uid, @api_key, @secret_key)
    ;
