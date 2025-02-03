from pathlib import Path

from sqlalchemy import create_engine, Column, Integer, Text, TIMESTAMP
from sqlalchemy.ext.declarative import declarative_base
from sqlalchemy.orm import sessionmaker

Base = declarative_base()


class Task(Base):
    __tablename__ = 'task'
    id = Column(Integer, primary_key=True)
    owner_name = Column(Text, nullable=False)
    title = Column(Text, nullable=False)
    live_id = Column(Text, nullable=False, unique=True)
    live_time = Column(TIMESTAMP)
    status = Column(Text, nullable=False)
    error_info = Column(Text, nullable=False)
    details = Column(Text, nullable=False)
    created = Column(TIMESTAMP, nullable=False)
    updated = Column(TIMESTAMP)


db_path = Path(__file__).parent.parent.joinpath('liveguard.db')
engine = create_engine(f'sqlite:///{str(db_path.absolute())}')
Session = sessionmaker(bind=engine)


def fetch_tasks():
    with Session() as session:
        tasks = session.query(Task).all()
        return tasks


def update_task(task_dict):
    with Session() as session:
        task = session.query(Task).filter(Task.id == task_dict['id']).one()
        task.owner_name = task_dict['owner_name']
        task.title = task_dict['title']
        task.live_id = task_dict['live_id']
        task.live_time = task_dict['live_time']
        task.status = task_dict['status']
        task.error_info = task_dict['error_info']
        task.created = task_dict['created']
        task.updated = task_dict['updated']
        session.commit()
