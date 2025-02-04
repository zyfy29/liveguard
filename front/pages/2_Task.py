import streamlit as st

from front.api_client import get_tasks, retry_task


def task_page():
    st.title('Summarize Tasks')
    tasks = get_tasks()
    for task in tasks:
        task['live_time'] = task['live_time']['Time']
        task['updated'] = task['updated']['Time']
        task['select'] = False

    edited_task_dicts = st.data_editor(tasks)

    count = 0
    if st.button('Restore'):
        for task in edited_task_dicts:
            if task['select']:
                res = retry_task(task['id'])
                if res.get('code') == 0:
                    count += 1
    if count > 0:
        st.success(f'{count} tasks restored')


task_page()
