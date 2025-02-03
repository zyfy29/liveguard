import streamlit as st

from api_client import *

st.set_page_config(page_title='Task Table', layout='wide')


def main_page():
    from db import fetch_tasks, update_task

    st.title('Summarize Tasks')

    if st.button('Retry Failed Tasks'):
        res = retry_tasks()
        if res.get('code') == 0:
            st.success('OK')
        else:
            st.error("Operation failed")

    tasks = fetch_tasks()

    # Convert tasks to a list of dictionaries, excluding the 'details' field
    task_dicts = [
        {
            'id': task.id,
            'owner_name': task.owner_name,
            'title': task.title,
            'live_id': task.live_id,
            'live_time': task.live_time,
            'status': task.status,
            'error_info': task.error_info,
            'created': task.created,
            'updated': task.updated
        }
        for task in tasks or []
    ]
    task_dicts.sort(key=lambda x: x['id'], reverse=True)

    # Display the data using st.data_editor
    edited_task_dicts = st.data_editor(task_dicts)

    # Check for changes and update the database
    for original, edited in zip(task_dicts, edited_task_dicts):
        if original != edited:
            update_task(edited)


def live_history_page():
    st.title('Live History')
    members = get_members()
    member = st.selectbox(
        "Select a member",
        members,
        format_func=lambda x: x.get('name')
    )

    lives, next_time = get_lives(owner_id=member.get('member_id'))

    # Add a new column for checkboxes
    for live in lives:
        live['select'] = False

    edited_lives = st.data_editor(lives, width=1000, height=300)

    # todo use session_state?
    # def go_front():
    #     nonlocal lives, next_time
    #     lives, next_time = get_lives(owner_id=member.get('member_id'), next_time=next_time)
    #     st.write(next_time)
    #
    # st.button("Next", key='next', on_click=lambda: go_front())

    # Add a submit button below the table
    count = 0
    if st.button('Submit'):
        for live in edited_lives:
            if live['select']:
                res = create_task(live['id'])
                if res.get('code') == 0:
                    count += 1
    if count > 0:
        st.success(f'{count} tasks created')


page_names_to_funcs = {
    'Sub': live_history_page,
    'Main': main_page,
}

page_name = st.sidebar.selectbox("Choose page", page_names_to_funcs.keys())
page_names_to_funcs[page_name]()
