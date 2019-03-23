"""create user table

Revision ID: a6c327aacb7b
Revises: 
Create Date: 2019-01-02 01:27:28.240713

"""
from alembic import op
import sqlalchemy as sa


# revision identifiers, used by Alembic.
revision = 'a6c327aacb7b'
down_revision = None
branch_labels = None
depends_on = None


def upgrade():
    op.create_table(
        'users',
        sa.Column('id', sa.BigInteger, primary_key=True),
        sa.Column('fullname', sa.String(50), nullable=False),
        sa.Column('username', sa.String(20)),
        sa.Column('email', sa.String(50)),
        sa.Column('password', sa.String(255)),
        sa.Column('status', sa.String(20)),
        sa.Column('created_at',
            sa.TIMESTAMP(timezone=True),
            server_default=sa.text('now()'),
            nullable=False
        ),
    )


def downgrade():
    op.drop_table('users')
