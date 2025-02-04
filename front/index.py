import streamlit as st

from front.pages.live import live_page
from front.pages.task import task_page

st.set_page_config(page_title='Liveguard Panel', layout='wide')

page_names_to_funcs = {
    'Live': live_page,
    'Task': task_page,
}

page_name = st.sidebar.selectbox("Choose page", page_names_to_funcs.keys())
page_names_to_funcs[page_name]()
