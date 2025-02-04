import streamlit as st

from front.api_client import get_tasks


def task_page():
    st.title('Summarize Tasks')
    tasks = get_tasks()
    edited_task_dicts = st.data_editor(tasks)
