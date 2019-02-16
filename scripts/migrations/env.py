from __future__ import with_statement

import enum

from logging.config import fileConfig

from sqlalchemy import engine_from_config
from sqlalchemy import pool, Column, String
from sqlalchemy.ext.declarative import declarative_base

from alembic import context

# this is the Alembic Config object, which provides
# access to the values within the .ini file in use.
config = context.config

# Interpret the config file for Python logging.
# This line sets up loggers basically.
fileConfig(config.config_file_name)

# add your model's MetaData object here
# for 'autogenerate' support
# from myapp import mymodel
# target_metadata = mymodel.Base.metadata
target_metadata = None

# other values from the config, defined by the needs of env.py,
# can be acquired:
# my_important_option = config.get_main_option("my_important_option")
# ... etc.

# class User(MetaData):
#     __tablename__ = 'user'
#     id = Column(Integer, primary_key=True)
#     created_at = Column(DateTime)
#     updated_at = Column(DateTime)

#     email = Column(String(256), nullable=True, index=True)
#     user_id = Column(Integer, nullable=False, index=True)
#     chat_id = Column(BigInteger, nullable=False)
#     email = Column(String(256), index=True)
#     name = Column(String(256))
#     username = Column(String(256), index=True)
#     has_unread_messages = Column(Boolean, nullable=False)
#     authtoken = Column(String(64))
#     is_authorized = Column(Boolean, nullable=False)
#     is_token_expired = Column(Boolean, nullable=False)
#     last_message_at = Column(DateTime, nullable=False)
#     last_message = Column(Integer, ForeignKey('message.id'), nullable=False)
#     user_photo_id = Column(String(512))
#     is_active = Column(Boolean, nullable=False)


# class Message(MetaData):
#     __tablename__ = 'message'
#     id = Column(Integer, primary_key=True)
#     created_at = Column(DateTime)
#     updated_at = Column(DateTime)

#     user_id = Column(Integer, ForeignKey('user.id'), nullable=False,
#                      index=True)
#     message_id = Column(Integer, nullable=False, index=True)
#     chat_id = Column(BigInteger, nullable=False)
#     text = Column(String(4000), nullable=True)
#     document_id = Column(String(512), nullable=True)
#     photo_id = Column(String(512), nullable=True)


def run_migrations_offline():
    """Run migrations in 'offline' mode.

    This configures the context with just a URL
    and not an Engine, though an Engine is acceptable
    here as well.  By skipping the Engine creation
    we don't even need a DBAPI to be available.

    Calls to context.execute() here emit the given string to the
    script output.

    """
    url = config.get_main_option("sqlalchemy.url")
    context.configure(
        url=url, target_metadata=target_metadata, literal_binds=True
    )

    with context.begin_transaction():
        context.run_migrations()


def run_migrations_online():
    """Run migrations in 'online' mode.

    In this scenario we need to create an Engine
    and associate a connection with the context.

    """
    connectable = engine_from_config(
        config.get_section(config.config_ini_section),
        prefix="sqlalchemy.",
        poolclass=pool.NullPool,
    )

    with connectable.connect() as connection:
        context.configure(
            connection=connection, target_metadata=target_metadata
        )

        with context.begin_transaction():
            context.run_migrations()


if context.is_offline_mode():
    run_migrations_offline()
else:
    run_migrations_online()
