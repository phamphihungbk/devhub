FROM postgres:15

# Optional: Add custom configuration or initialization scripts
# COPY ./init.sql /docker-entrypoint-initdb.d/

# Optional: Add custom config
# COPY ./postgresql.conf /etc/postgresql/postgresql.conf
# ENV POSTGRES_CONFIG_FILE=/etc/postgresql/postgresql.conf