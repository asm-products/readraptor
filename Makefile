.PHONY: restore

restore:
  heroku pgbackups:capture --expire --app readraptor
  curl `heroku pgbackups:url --app readraptor` -o db/latest.dump
  pg_restore --verbose --clean --no-acl --no-owner -h localhost -d rr_development db/latest.dump
