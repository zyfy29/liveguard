import streamlit as st

from front.api_client import delete_task, get_tasks, retry_task


def task_page():
    st.set_page_config(layout='wide')
    st.title('Summarize Tasks')
    tasks = get_tasks()['data']
    for task in tasks:
        task['live_time'] = task['live_time']['Time']
        task['updated'] = task['updated']['Time']
        task['info'] = task['error_info']
        del task['error_info']
        task['select'] = False

    edited_task_dicts = st.data_editor(tasks)

    count = 0
    col1, col2, _ = st.columns(3)
    with col1:
        if st.button('Restore'):
            for task in edited_task_dicts:
                if task['select']:
                    res = retry_task(task['id'])
                    if res.get('code') == 0:
                        count += 1
    with col2:
        if st.button('Delete'):
            for task in edited_task_dicts:
                if task['select']:
                    res = delete_task(task['id'])
                    if res.get('code') == 0:
                        count += 1
    if count > 0:
        st.success(f'{count} tasks affected.')


task_page()
