from pathlib import Path

import streamlit as st

from front.api_client import get_my_name


def get_readme_content():
    readme_path = Path(__file__).parent / 'readme.md'
    with open(readme_path, encoding='utf-8') as f:
        text = f.read()
    return text


def index_page():
    st.set_page_config(page_title='Liveguard Panel', layout='wide')

    index_content = get_readme_content()
    st.markdown(index_content)

    my_name_resp = get_my_name()
    if my_name_resp.get('code') == 0:
        st.sidebar.success(f"Logged in as {my_name_resp.get('data').get('nickname')}.")
    else:
        st.sidebar.error(f"Login failed:\n{my_name_resp.get('msg')}")

index_page()
