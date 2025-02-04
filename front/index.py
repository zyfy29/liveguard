from pathlib import Path

import streamlit as st


def get_readme_content():
    readme_path = Path(__file__).parent / 'readme.md'
    with open(readme_path, encoding='utf-8') as f:
        text = f.read()
    return text


st.set_page_config(page_title='Liveguard Panel', layout='wide')
st.sidebar.success("Select a page above.")

index_content = get_readme_content()
st.markdown(index_content)
