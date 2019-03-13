
-- encrypted bcrypt password (Plain Text: "password")
-- equivalent Go code: bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

SET @defaultPassword = '$2a$10$MLooRpCcdyyxoXwe3ZiCFuQZfsGeVC7TPCSyYhTs8Bl/sFPd4K67W';

INSERT INTO users
    (slug, first_name, last_name, email, `admin`) 
VALUES
    ('adefaultuser', 'Alice',   'DefaultUser', 'adefaultuser@example.com', true),
    ('bsnooper',     'Bob',     'Snooper',     'bsnooper@example.com', false),
    ('cother',       'Carl',    'Other',       'cother@example.com', false),
    ('dother',       'Debbie',  'Other',       'dother@example.com', false),
    ('eother',       'Edith',   'Other',       'eother@example.com', false),
    ('fother',       'Fred',    'Other',       'fother@example.com', false),
    ('gother',       'Ginna',   'Other',       'gother@example.com', false),
    ('hother',       'Hazel',   'Other',       'hother@example.com', false),
    ('iother',       'Ivy',     'Other',       'iother@example.com', false),
    ('jother',       'Justin',  'Other',       'jother@example.com', false),
    ('kother',       'Kevin',   'Other',       'kother@example.com', false),
    ('lother',       'Liz',     'Other',       'lother@example.com', false),
    ('mother',       'Martin',  'Other',       'mother@example.com', false),
    ('nother',       'Nick',    'Other',       'nother@example.com', false),
    ('oother',       'Ophelia', 'Other',       'oother@example.com', false),
    ('pother',       'Paul',    'Other',       'pother@example.com', false),
    ('qother',       'Quinton', 'Other',       'qother@example.com', false),
    ('rother',       'Rachel',  'Other',       'rother@example.com', false),
    ('tother',       'Tess',    'Other',       'tother@example.com', false),
    ('uother',       'Ursula',  'Other',       'uother@example.com', false),
    ('vother',       'Vera',    'Other',       'vother@example.com', false),
    ('wother',       'Winston', 'Other',       'wother@example.com', false),
    ('xother',       'Xavier',  'Other',       'xother@example.com', false),
    ('yother',       'Yavonne', 'Other',       'yother@example.com', false),
    ('zother',       'Zoey',    'Other',       'zother@example.com', false),
    ('smoker',       'Smoke',   'Tester',      'smoker', false);

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

INSERT INTO auth_scheme_data
    (auth_scheme, user_key, user_id, encrypted_password)
VALUES
    ('local', 'adefaultuser@example.com', @alice_uid,     @defaultPassword),
    ('local', 'bsnooper@example.com',     @bob_uid,       @defaultPassword),
    ('local', 'Carl',                     @carl_uid,      @defaultPassword),
    ('local', 'Debbie',                   @debbie_uid,    @defaultPassword),
    ('local', 'Edith',                    @edith_uid,     @defaultPassword),
    ('local', 'Fred',                     @fred_uid,      @defaultPassword),
    ('local', 'Ginna',                    @ginna_uid,     @defaultPassword),
    ('local', 'Hazel',                    @hazel_uid,     @defaultPassword),
    ('local', 'Ivy',                      @ivy_uid,       @defaultPassword),
    ('local', 'Justin',                   @justin_uid,    @defaultPassword),
    ('local', 'Kevin',                    @kevin_uid,     @defaultPassword),
    ('local', 'Liz',                      @liz_uid,       @defaultPassword),
    ('local', 'Martin',                   @martin_uid,    @defaultPassword),
    ('local', 'Nick',                     @nick_uid,      @defaultPassword),
    ('local', 'Ophelia',                  @ophelia_uid,   @defaultPassword),
    ('local', 'Paul',                     @paul_uid,      @defaultPassword),
    ('local', 'Quinton',                  @quinton_uid,   @defaultPassword),
    ('local', 'Rachel',                   @rachel_uid,    @defaultPassword),
    ('local', 'Tess',                     @tess_uid,      @defaultPassword),
    ('local', 'Ursula',                   @ursula_uid,    @defaultPassword),
    ('local', 'Vera',                     @vera_uid,      @defaultPassword),
    ('local', 'Winston',                  @winston_uid,   @defaultPassword),
    ('local', 'Xavier',                   @xavier_uid,    @defaultPassword),
    ('local', 'Yavonne',                  @yavonne_uid,   @defaultPassword),
    ('local', 'Zoey',                     @zoey_uid,      @defaultPassword);

INSERT INTO
    operations (`slug`, `name`, `description`, active, `status`)
VALUES
    ('alice-op', 'AliceOp', 'An operation for Alice',  1, 0),
    ('bob-op',   'BobOp',   'An operation for Bob',    1, 0),
    ('co-op',    'Co-Op',   'A Cooperative Operation', 1, 0),
    ('no-op',    'No-Op',   'An Orphaned Operation',   1, 0),
    ('big-op',   'BigOp',   'An operation with lots of users', 1, 0);

SET @alice_op_id = 1;
SET @bob_op_id = 2;
SET @co_op_id = 3;
SET @no_op_id = 4;
SET @big_op_id = 5;

INSERT INTO user_operation_permissions
    (`user_id`, `operation_id`, `role`)
VALUES
    (@alice_uid, @alice_op_id, 'admin'),
    (@bob_uid,   @bob_op_id,   'admin'),
    (@alice_uid, @co_op_id,    'admin'),
    (@bob_uid,   @co_op_id,    'admin'),

    (@alice_uid,   @big_op_id, 'admin'),
    (@carl_uid,    @big_op_id, 'read'),
    (@debbie_uid,  @big_op_id, 'write'),
    (@edith_uid,   @big_op_id, 'read'),
    (@fred_uid,    @big_op_id, 'write'),
    (@ginna_uid,   @big_op_id, 'read'),
    (@hazel_uid,   @big_op_id, 'read'),
    (@ivy_uid,     @big_op_id, 'write'),
    (@justin_uid,  @big_op_id, 'read'),
    (@kevin_uid,   @big_op_id, 'write'),
    (@liz_uid,     @big_op_id, 'read'),
    (@martin_uid,  @big_op_id, 'read'),
    (@nick_uid,    @big_op_id, 'read'),
    (@ophelia_uid, @big_op_id, 'write'),
    (@paul_uid,    @big_op_id, 'read'),
    (@quinton_uid, @big_op_id, 'write'),
    (@rachel_uid,  @big_op_id, 'read'),
    (@tess_uid,    @big_op_id, 'write'),
    (@ursula_uid,  @big_op_id, 'write'),
    (@vera_uid,    @big_op_id, 'read'),
    (@winston_uid, @big_op_id, 'write'),
    (@xavier_uid,  @big_op_id, 'read'),
    (@yavonne_uid, @big_op_id, 'write'),
    (@zoey_uid,    @big_op_id, 'read');

INSERT INTO tags
    (name, color_name, operation_id)
VALUES
    ('Europa', 'red',    @alice_op_id),
    ('Titan',  'orange', @alice_op_id),
    ('Io',     'yellow', @alice_op_id),
    ('Ceres',  'green',  @alice_op_id),
    ('Triton', 'blue',   @alice_op_id),

    ('Doc',     'red',    @bob_op_id),
    ('Grumpy',  'orange', @bob_op_id),
    ('Happy',   'yellow', @bob_op_id),
    ('Sleepy',  'green',  @bob_op_id),
    ('Bashful', 'blue',   @bob_op_id),
    ('Sneezy',  'indigo', @bob_op_id),
    ('Dopey',   'violet', @bob_op_id),

    ('application',  'lightRed',    @co_op_id),
    ('presentation', 'lightOrange', @co_op_id),
    ('session',      'lightYellow', @co_op_id),
    ('transport',    'lightGreen',  @co_op_id),
    ('network',      'lightBlue',   @co_op_id),
    ('data link',    'lightIndigo', @co_op_id),
    ('physical',     'lightViolet', @co_op_id);

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

INSERT INTO evidence
    (uuid, operation_id, operator_id, content_type, full_image_key, thumb_image_key, occurred_at, `description`)
VALUES
    (@a_op_evi_uuid_1, @alice_op_id, @alice_uid, 'image',     'seed_movie_full',    'seed_movie_thumb',    now(), CONCAT_WS(CHAR(10 using utf8), '# Movie Reel', '', 'Cinema''s favorite feature. Action. Excitement. Comedy.', 'Reviews:', '', '* A loud, long and pointless spectacle.', '* Schumacher''s storytelling is limp, and the characters lack energy.', '', '[Click here](https://www.rottentomatoes.com/m/1077027_batman_and_robin) for more info' )),
    (@a_op_evi_uuid_2, @alice_op_id, @alice_uid, 'image',     'seed_popcorn_full',  'seed_popcorn_thumb',  now(), 'Popcorn Box'),
    (@a_op_evi_uuid_3, @alice_op_id, @alice_uid, 'codeblock', 'seed_go_aoc201614',  'seed_go_aoc201614',   now(), 'Go AOC 2016 Day 14 (https://adventofcode.com/2016/day/14)'),

    (@b_op_evi_uuid_1, @bob_op_id,   @bob_uid,   'image',     'seed_magazine_full', 'seed_magazine_thumb', now(), 'Magazine'),
    (@b_op_evi_uuid_2, @bob_op_id,   @bob_uid,   'image',     'seed_chicken_full',  'seed_chicken_thumb',  now(), 'Chickens with a Frisbee'),
    (@b_op_evi_uuid_3, @bob_op_id,   @bob_uid,   'codeblock', 'seed_py_aoc201717',  'seed_py_aoc201717',   now(), 'Python AOC 2017 Day 17 (https://adventofcode.com/2017/day/17)'),

    (@c_op_evi_uuid_1, @co_op_id,    @alice_uid, 'image',     'seed_pocky_full',    'seed_pocky_thumb',    now(), 'A miko holding a bunch of ofudas'),
    (@c_op_evi_uuid_2, @co_op_id,    @bob_uid,   'image',     'seed_rocky_full',    'seed_rocky_thumb',    now(), 'A raccoon juggling some leaves'),
    (@c_op_evi_uuid_3, @co_op_id,    @bob_uid,   'codeblock', 'seed_rs_aoc201501',  'seed_rs_aoc201501',   now(), 'Rust AOC 2015 Day 1 (https://adventofcode.com/2015/day/1)');


SET @a_op_evi_1 = 1;
SET @a_op_evi_2 = 2;
SET @a_op_evi_3 = 3;
SET @b_op_evi_1 = 4;
SET @b_op_evi_2 = 5;
SET @b_op_evi_3 = 6;
SET @c_op_evi_1 = 7;
SET @c_op_evi_2 = 8;
SET @c_op_evi_3 = 9;

SET @a_op_evt_uuid_1 = 'a10F0000-0000-4000-a000-000000000000';
SET @a_op_evt_uuid_2 = 'a20F0000-0000-4000-a000-000000000000';

SET @b_op_evt_uuid_1 = 'b10F0000-0000-4000-b000-000000000000';
SET @b_op_evt_uuid_2 = 'b20F0000-0000-4000-b000-000000000000';

SET @c_op_evt_uuid_1 = 'c10F0000-0000-4000-8000-000000000000';
SET @c_op_evt_uuid_2 = 'c20F0000-0000-4000-8000-000000000000';

INSERT INTO findings
    (`uuid`, `operation_id`, `category`, `title`, `description`, `ready_to_report`, `ticket_link`)
VALUES
    (@a_op_evt_uuid_1, @alice_op_id, 'OPSEC', 'Main Event',                'body', true,  'http://google.com'), -- 1
    (@a_op_evt_uuid_2, @alice_op_id, 'OPSEC', 'Side Show left',            'body', true,  null), -- 2
    (@b_op_evt_uuid_1, @bob_op_id,   'OPSEC', 'Bob Sees an Issue',         'body', false, null), -- 3
    (@b_op_evt_uuid_2, @bob_op_id,   'OPSEC', 'Bob Suspects Fowl Play',    'body', false, null), -- 4
    (@c_op_evt_uuid_1, @co_op_id,    'OPSEC', 'I get Pocky',               'body', false, null), -- 5
    (@c_op_evt_uuid_2, @co_op_id,    'OPSEC', 'Bob gets stuck with Rocky', 'body', false, null); -- 6

SET @a_op_evt_1 = 1;
SET @a_op_evt_2 = 2;
SET @b_op_evt_1 = 3;
SET @b_op_evt_2 = 4;
SET @c_op_evt_1 = 5;
SET @c_op_evt_2 = 6;

INSERT INTO evidence_finding_map 
    (evidence_id, finding_id)
VALUES
    (@a_op_evi_1, @a_op_evt_1),
    (@a_op_evi_2, @a_op_evt_2),
    (@a_op_evi_3, @a_op_evt_2),
    (@b_op_evi_1, @b_op_evt_1),
    (@b_op_evi_2, @b_op_evt_2),
    (@b_op_evi_3, @b_op_evt_2),
    (@c_op_evi_1, @c_op_evt_1),
    (@c_op_evi_2, @c_op_evt_2),
    (@c_op_evi_3, @c_op_evt_2);

INSERT INTO tag_evidence_map
    (tag_id, evidence_id)
VALUES
    (@tag_moon_1_id, @a_op_evi_1),
    (@tag_moon_2_id, @a_op_evi_1),
    (@tag_moon_3_id, @a_op_evi_2),
    (@tag_moon_2_id, @a_op_evi_2),
    (@tag_moon_4_id, @a_op_evi_3),
    (@tag_moon_5_id, @a_op_evi_3),

    (@tag_dwarf_4_id, @b_op_evi_1),
    (@tag_dwarf_5_id, @b_op_evi_1),
    (@tag_dwarf_1_id, @b_op_evi_2),
    (@tag_dwarf_4_id, @b_op_evi_2),
    (@tag_dwarf_7_id, @b_op_evi_3),
    (@tag_dwarf_6_id, @b_op_evi_3),

    (@tag_osi_6_id, @c_op_evi_1),
    (@tag_osi_5_id, @c_op_evi_1),
    (@tag_osi_7_id, @c_op_evi_2),
    (@tag_osi_1_id, @c_op_evi_2),
    (@tag_osi_3_id, @c_op_evi_3),
    (@tag_osi_2_id, @c_op_evi_3);
    
SET @api_key = 'DAYPFGHnm1Pqes-l0Fm76_y1';
SET @secret_key = 0x1EA9AE5B294BCE747EB6A4A8B5900E73EC38EDBB9215A28A4C9A33A5711892E3418AE449830DCD7893AE54FEA46D00509A2613A9801A88829B70ED41C5C1F9D9;
-- secret_key: HqmuWylLznR+tqSotZAOc+w47buSFaKKTJozpXEYkuNBiuRJgw3NeJOuVP6kbQBQmiYTqYAaiIKbcO1BxcH52Q==

INSERT INTO api_keys
    (user_id, access_key, secret_key)
VALUES
    (@smoker_uid, @api_key, @secret_key)
