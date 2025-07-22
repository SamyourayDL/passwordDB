-- INSERT INTO users (name)
-- VALUES
--     ('Alice'),
--     ('Bob');

INSERT INTO passwords (user_id, service_name, secret_enc) VALUES
    (8, 'gmail', decode('c4045bab291ce97bb12db90a80efad81a2feb0527ef695f5339a0adfe9ff8dfafd516280fd0245', 'hex')),
    (8, 'github', decode('6d6d032c3695bf8af6e257386fae37642f732457cff1e56950bda4996798b40d4c6165d31a', 'hex'));