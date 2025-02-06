import streamlit as st

from front.api_client import create_task, get_lives, get_members


def live_page():
    st.set_page_config(layout='wide')
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

    count = 0
    if st.button('Submit'):
        for live in edited_lives:
            if live['select']:
                res = create_task(live['id'])
                if res.get('code') == 0:
                    count += 1
    if count > 0:
        st.success(f'{count} tasks created')


live_page()
