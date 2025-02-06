import streamlit as st
from dateutil import parser

from front.api_client import get_task_detail, get_tasks


def format_iso_time(iso_str):
    dt_obj = parser.isoparse(iso_str)
    return dt_obj.strftime("%Y年%-m月%-d日%H:%M:%S")


def story_page():
    st.set_page_config(layout='wide')
    st.title('Story')
    done_tasks = get_tasks('succeed')['data']
    display_task = st.selectbox(
        "Select a story to read",
        done_tasks,
        format_func=lambda x: f'{x.get('owner_name')}-{x.get('title')}-{x.get('live_time').get('Time')}'
    )

    task_detail = get_task_detail(display_task.get('id'))['data']

    st.title(task_detail.get('title'))
    st.text(f"{task_detail.get('owner_name')}\t{format_iso_time(task_detail.get('live_time').get('Time'))}")
    st.text(
        f"{task_detail.get('duration') // 60}分钟\t总结{len(task_detail.get('summary'))}字\t转录{len(task_detail.get('transcript'))}字")
    st.divider()

    st.markdown(task_detail.get('summary'))
    st.divider()

    st.title('Transcript')
    st.markdown(task_detail.get('transcript'))


story_page()
