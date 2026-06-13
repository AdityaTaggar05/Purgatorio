mkdir -p ./backend/keys &&
openssl genrsa -out ./backend/keys/jwt_private.pem 2048 &&
openssl rsa -in ./backend/keys/jwt_private.pem -pubout -out ./backend/keys/jwt_public.pem

echo PURG_DB_PASS=purg_db_pass123 >> .env
