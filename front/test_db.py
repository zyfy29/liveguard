from unittest.mock import patch

import pytest

from front.db import fetch_tasks, Task


@pytest.fixture
def mock_session():
    with patch('front.db.Session') as MockSession:
        yield MockSession


@pytest.mark.usefixtures("mock_session")
def test_fetch_tasks(mock_session):
    # Arrange
    mock_task = Task(id=1, owner_name='John Doe', title='Sample Task', live_id='123', live_time=None, status='pending',
                     error_info=None, details='Details', _created=None, _updated=None)
    mock_session_instance = mock_session.return_value
    mock_session_instance.query.return_value.all.return_value = [mock_task]

    # Act
    tasks = fetch_tasks()

    # Assert
    assert len(tasks) == 1
    assert tasks[0].id == 1
    assert tasks[0].owner_name == 'John Doe'
    assert tasks[0].title == 'Sample Task'
