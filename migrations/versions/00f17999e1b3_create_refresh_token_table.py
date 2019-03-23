"""create refresh_token table

Revision ID: 00f17999e1b3
Revises: a6c327aacb7b
Create Date: 2019-01-03 21:34:12.823192

"""
from alembic import op
import sqlalchemy as sa


# revision identifiers, used by Alembic.
revision = '00f17999e1b3'
down_revision = 'a6c327aacb7b'
branch_labels = None
depends_on = None


def upgrade():
    op.create_table(
        'refresh_token',
        sa.Column('id', sa.BigInteger, primary_key=True),
        sa.Column('user_id', sa.BigInteger, nullable=False),
        sa.Column('refresh_token', sa.String(255), nullable=False),
        sa.Column('revoked', sa.SmallInteger(), server_default=sa.schema.DefaultClause("0"), nullable=False),
    )


def downgrade():
    op.drop_table('refresh_token')
