"""Create base schema

Revision ID: 200e7ea0e64b
Revises: 
Create Date: 2019-02-16 12:15:59.212678

"""
from alembic import op
import sqlalchemy as sa


# revision identifiers, used by Alembic.
revision = '200e7ea0e64b'
down_revision = None
branch_labels = None
depends_on = None


def upgrade():
    op.create_table(
        'admins',
        sa.Column('id', sa.Integer, primary_key=True, index=True),
        sa.Column('created_at', sa.DateTime),
        sa.Column('updated_at', sa.DateTime),

        sa.Column('email', sa.String(256), nullable=False, index=True, unique=True),
        sa.Column('name', sa.String(256)),
        sa.Column('hashed_password', sa.String(256), nullable=False),
        sa.Column('is_superuser', sa.Boolean, nullable=False),
        sa.Column('is_active', sa.Boolean, nullable=False),
        sa.Column('email_confirmed', sa.Boolean, nullable=False),
        sa.Column('auth_token', sa.String(64)),
        sa.Column('password_reset_token', sa.String(64)),
    )
    # FIXME separate messangers should have separate tables
    op.create_table(
        'users',
        sa.Column('user_id', sa.Integer, primary_key=True, nullable=False, index=True),
        sa.Column('created_at', sa.DateTime),
        sa.Column('updated_at', sa.DateTime),

        sa.Column('chat_id', sa.BigInteger, nullable=False),
        sa.Column('email', sa.String(256), index=True),
        sa.Column('name', sa.String(256)),
        sa.Column('username', sa.String(256)),
        sa.Column('has_unread_messages', sa.Boolean, nullable=False),
    )
    op.create_table(
        'messages',
        sa.Column('id', sa.Integer, primary_key=True),
        sa.Column('user_id', sa.Integer, sa.ForeignKey('users.user_id', ondelete="CASCADE"), nullable=False, index=True),
        sa.Column('message_id', sa.Integer, index=True, unique=True),
        sa.Column('is_broadcast', sa.Boolean, nullable=False),
        sa.Column('from_admin', sa.Boolean, nullable=False),
        sa.Column('created_at', sa.DateTime),
        sa.Column('updated_at', sa.DateTime),
        sa.Column('text', sa.String(4000), nullable=False),
        sa.Column('reply_to_message', sa.Integer),
        sa.Column('document_id', sa.String, nullable=False),
        sa.Column('photo_id', sa.String, nullable=False),
    )
    op.add_column('users',
        sa.Column('last_message_id', sa.Integer, sa.ForeignKey('messages.message_id', ondelete="CASCADE"), index=True),
    )


def downgrade():
    op.drop_table('admins')
    op.drop_table('users')
    op.drop_table('messages')
