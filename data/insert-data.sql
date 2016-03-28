\connect go_twitter_bot;

DO $$
    DECLARE user_id UUID;
    DECLARE twitter_account_id UUID;
BEGIN
    INSERT INTO users(email, hashed_password, is_admin, date_created)
    VALUES ('your@email.com', '$2a$10$oF2TzJDQO7VuKQR3y.5bne.vGIOEWGNpE8T1VVLNLLX.QKKj8bifa', true, LOCALTIMESTAMP(0))
    RETURNING id INTO user_id;

    INSERT INTO twitter_accounts(user_id, username, date_created, consumer_key, consumer_secret, access_token, access_token_secret)
    VALUES (user_id, 'myusername', LOCALTIMESTAMP(0), 'CONSUMER_KEY', 'CONSUMER_SECRET', 'ACCESS_TOKEN', 'ACCESS_TOKEN_SECRET')
    RETURNING id INTO twitter_account_id;



    INSERT INTO tweets(twitter_account_id, tweet, post_on, is_posted, date_created)
    VALUES (
        twitter_account_id,
        'Some people, when confronted with a problem, think
"I know, I''ll use regular expressions." Now they have two problems.',
        '2016-03-22 19:30:00',
        false,
        LOCALTIMESTAMP(0)
    );

    INSERT INTO tweets(twitter_account_id, tweet, post_on, is_posted, date_created)
    VALUES (
        twitter_account_id,
        'Some people, when confronted with a problem, think
"I know, I''ll use binary." Now they have 10 problems.',
        '2016-03-23 19:30:00',
        false,
        LOCALTIMESTAMP(0)
    );

    INSERT INTO tweets(twitter_account_id, tweet, post_on, is_posted, date_created)
    VALUES (
        twitter_account_id,
        'Some people, when confronted with a problem, think
"I know, I''ll use threading." Now have two they pborlesm.',
        '2016-03-24 19:30:00',
        false,
        LOCALTIMESTAMP(0)
    );

END$$;
